package apperr

import "fmt"

type AppError struct {
	ErrorCode string `json:"errorCode"`
	HTTPCode  int    `json:"httpCode"`
	Info      string `json:"info"`
	Message   string `json:"message"`
	Raw       error  `json:"raw"`
}

var _ error = (*AppError)(nil)

func (ae AppError) Error() string {
	return fmt.Sprintf("%#v", ae)
}

func New(errorCode string, httpCode int, info string, message string, raw error) *AppError {
	return &AppError{errorCode, httpCode, info, message, raw}
}
