package domain

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

type ErrorValidationResponse struct {
	Message string            `json:"message"`
	Errors  []ValidationError `json:"errors"`
}

type ValidationError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}
