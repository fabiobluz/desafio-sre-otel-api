package model

// CEPRequest representa a estrutura de requisição para validação de CEP
type CEPRequest struct {
	CEP string `json:"cep"`
}

// ViaCEPResponse representa a resposta da API ViaCEP
type ViaCEPResponse struct {
	CEP        string `json:"cep"`
	Localidade string `json:"localidade"`
	Erro       string `json:"erro"`
}

// WeatherAPIResponse representa a resposta da API WeatherAPI
type WeatherAPIResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

// WeatherResponse representa a resposta final do sistema
type WeatherResponse struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

// ServiceBError representa erros específicos do Service B
type ServiceBError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Type    string `json:"type"` // "zipcode_not_found", "weather_api_error", "validation_error"
}

// HTTPStatus representa códigos de status HTTP comuns
type HTTPStatus struct {
	OK                  int `json:"ok"`                    // 200
	BadRequest          int `json:"bad_request"`           // 400
	UnprocessableEntity int `json:"unprocessable_entity"`  // 422
	NotFound            int `json:"not_found"`             // 404
	InternalServerError int `json:"internal_server_error"` // 500
}

// NewHTTPStatus retorna uma instância com os códigos de status padrão
func NewHTTPStatus() HTTPStatus {
	return HTTPStatus{
		OK:                  200,
		BadRequest:          400,
		UnprocessableEntity: 422,
		NotFound:            404,
		InternalServerError: 500,
	}
}

// ErrorResponse representa uma resposta de erro padrão
type ErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// TemperatureConversion representa funções de conversão de temperatura
type TemperatureConversion struct {
	Celsius    float64 `json:"celsius"`
	Fahrenheit float64 `json:"fahrenheit"`
	Kelvin     float64 `json:"kelvin"`
}

// ConvertCelsiusToFahrenheit converte Celsius para Fahrenheit
func (tc *TemperatureConversion) ConvertCelsiusToFahrenheit(celsius float64) float64 {
	return celsius*1.8 + 32
}

// ConvertCelsiusToKelvin converte Celsius para Kelvin
func (tc *TemperatureConversion) ConvertCelsiusToKelvin(celsius float64) float64 {
	return celsius + 273
}

// SetTemperatures define todas as temperaturas baseadas no Celsius
func (tc *TemperatureConversion) SetTemperatures(celsius float64) {
	tc.Celsius = celsius
	tc.Fahrenheit = tc.ConvertCelsiusToFahrenheit(celsius)
	tc.Kelvin = tc.ConvertCelsiusToKelvin(celsius)
}
