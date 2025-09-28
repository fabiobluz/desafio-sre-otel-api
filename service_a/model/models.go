package model

// CEPRequest representa a estrutura de requisição para validação de CEP
type CEPRequest struct {
	CEP string `json:"cep"`
}

// CEPValidationError representa erros de validação do CEP
type CEPValidationError struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// ServiceAResponse representa a resposta padrão do Service A
type ServiceAResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
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
