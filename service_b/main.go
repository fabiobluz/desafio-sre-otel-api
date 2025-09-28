package main

// Inclua aqui todos os imports necessários
import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"service_b/model"
)

// Função auxiliar: buscar cidade pelo CEP
func fetchCity(ctx context.Context, cep string) (string, error) {
	tracer := otel.Tracer("service-b-tracer")
	ctx, span := tracer.Start(ctx, "fetch-city-via-viacep", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	url := "http://viacep.com.br/ws/" + cep + "/json/"
	span.SetAttributes(attribute.String("http.url", url))

	resp, err := http.Get(url)
	if err != nil {
		span.RecordError(err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		span.SetAttributes(attribute.Int("http.status_code", resp.StatusCode))
		return "", err
	}

	var data model.ViaCEPResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		span.RecordError(err)
		return "", err
	}

	if data.Erro != "" {
		span.SetAttributes(attribute.Bool("city.found", false))
		return "", http.ErrNoLocation
	}

	span.SetAttributes(attribute.String("city.name", data.Localidade))
	return data.Localidade, nil
}

// Função auxiliar: buscar clima pela cidade
func fetchWeather(ctx context.Context, city string) (float64, error) {
	tracer := otel.Tracer("service-b-tracer")
	ctx, span := tracer.Start(ctx, "fetch-weather-api", trace.WithSpanKind(trace.SpanKindClient))
	defer span.End()

	// Obter API key da variável de ambiente
	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		apiKey = "demo" // Para testes, usar chave demo
	}
	encodedCity := url.QueryEscape(city)
	url := "http://api.weatherapi.com/v1/current.json?key=" + apiKey + "&q=" + encodedCity

	resp, err := http.Get(url)
	if err != nil {
		span.RecordError(err)
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		span.SetAttributes(attribute.Int("http.status_code", resp.StatusCode))
		return 0, err
	}

	var data model.WeatherAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		span.RecordError(err)
		return 0, err
	}

	span.SetAttributes(attribute.Float64("temp.celsius", data.Current.TempC))
	return data.Current.TempC, nil
}

// Handler principal do Serviço B (sem métricas)
func WeatherHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tracer := otel.Tracer("service-b-tracer")

	// O Span principal do serviço B é criado implicitamente pelo Middleware
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("http.route", "/weather"))

	// 1. Decodificação e Validação
	var reqBody model.CEPRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		span.SetAttributes(attribute.String("error.type", "decode_error"))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "invalid request body"})
		return
	}

	cep := strings.TrimSpace(reqBody.CEP)
	span.SetAttributes(attribute.String("cep.value", cep))

	// Validação adicional do CEP no Service B
	if len(cep) != 8 || !regexp.MustCompile(`^\d{8}$`).MatchString(cep) {
		span.SetAttributes(attribute.Bool("cep.valid", false))
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]string{"message": "invalid zipcode"})
		return
	}

	// 2. Busca da Cidade
	city, err := fetchCity(ctx, cep)
	if err != nil {
		if err == http.ErrNoLocation {
			span.SetAttributes(attribute.String("error.type", "not_found"))
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"message": "can not find zipcode"})
		} else {
			span.SetAttributes(attribute.String("error.type", "viacep_error"))
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": "internal server error during city lookup"})
		}
		return
	}

	// 3. Busca do Clima
	tempC, err := fetchWeather(ctx, city)
	if err != nil {
		span.SetAttributes(attribute.String("error.type", "weather_api_error"))
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "internal server error during weather lookup"})
		return
	}

	// 4. Conversões e Resposta
	_, conversionSpan := tracer.Start(ctx, "temperature-conversion")
	tempF := tempC*1.8 + 32
	tempK := tempC + 273
	conversionSpan.End()

	respData := model.WeatherResponse{
		City:  city,
		TempC: tempC,
		TempF: tempF,
		TempK: tempK,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(respData)
}

func main() {
	// Inicialização da telemetria (apenas traces)
	shutdown, err := InitTelemetry("service-b")
	if err != nil {
		log.Fatalf("failed to initialize telemetry: %v", err)
	}
	defer func() {
		if err := shutdown(context.Background()); err != nil {
			log.Fatalf("failed to shut down telemetry: %v", err)
		}
	}()

	http.HandleFunc("/weather", WeatherHandler)
	log.Println("Service B running on :8081")
	http.ListenAndServe(":8081", nil)
}
