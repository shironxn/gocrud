package domain

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Errors  error  `json:"errors"`
}

type ValidationError struct {
	Field  string `json:"field"`
	Errors string `json:"errors,omitempty"`
}
