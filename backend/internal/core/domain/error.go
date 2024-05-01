package domain

type ErrorResponse struct {
	Code  int         `json:"code"`
	Error interface{} `json:"error"`
}

type ValidationError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}
