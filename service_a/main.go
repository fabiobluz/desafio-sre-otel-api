package main

// Inclua aqui todos os imports necessários
import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"regexp"
	"service_a/model"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// Handler principal do Serviço A (sem métricas)
func CEPValidationHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tracer := otel.Tracer("service-a-tracer")

	// Cria o Span principal
	ctx, span := tracer.Start(ctx, "CEPValidationHandler")
	defer span.End()
	span.SetAttributes(attribute.String("http.route", "/cep-weather"))

	// 1. Decodificação do Input
	var reqBody model.CEPRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		span.SetAttributes(attribute.String("error.message", "invalid request body"))
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": "invalid request body"})
		return
	}

	// 2. Validação do Input (8 dígitos numéricos e string)
	cep := strings.TrimSpace(reqBody.CEP)
	if len(cep) != 8 || !regexp.MustCompile(`^\d{8}$`).MatchString(cep) {
		span.SetAttributes(attribute.Bool("cep.valid", false))
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]string{"message": "invalid zipcode"})
		return
	}

	span.SetAttributes(attribute.Bool("cep.valid", true), attribute.String("cep.value", cep))

	// 3. Chamada HTTP para Serviço B (propagação do contexto)
	ctx, clientSpan := tracer.Start(ctx, "call-service-b", trace.WithSpanKind(trace.SpanKindClient))
	defer clientSpan.End()

	// Criar novo body com o CEP validado
	requestBody := model.CEPRequest{CEP: cep}
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		clientSpan.RecordError(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "internal server error"})
		return
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "http://service-b:8081/weather", bytes.NewBuffer(jsonBody))
	if err != nil {
		clientSpan.RecordError(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "internal server error"})
		return
	}

	// Propaga o contexto de tracing no Header da requisição HTTP
	req.Header.Set("Content-Type", "application/json")
	propagator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	propagator.Inject(ctx, propagation.HeaderCarrier(req.Header))

	// Executa a requisição
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		clientSpan.RecordError(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "could not connect to service B"})
		return
	}
	defer resp.Body.Close()

	// 4. Retorna a resposta (incluindo o status code do Serviço B)
	w.WriteHeader(resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		// Loga o erro do Serviço B no Span principal
		span.SetAttributes(attribute.Int("http.status_code", resp.StatusCode))
		io.Copy(w, resp.Body)
		return
	}

	io.Copy(w, resp.Body)
}

func main() {
	// Inicialização da telemetria (apenas traces)
	shutdown, err := InitTelemetry("service-a")
	if err != nil {
		log.Fatalf("failed to initialize telemetry: %v", err)
	}
	defer func() {
		if err := shutdown(context.Background()); err != nil {
			log.Fatalf("failed to shut down telemetry: %v", err)
		}
	}()

	http.HandleFunc("/cep-weather", CEPValidationHandler)
	log.Println("Service A running on :8080")
	http.ListenAndServe(":8080", nil)
}
