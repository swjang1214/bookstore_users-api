package errors

import "net/http"

type RestError struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	Error   string `json:"error"`
}

func NewBadRequestError(msg string) *RestError {
	return &RestError{
		Message: msg,
		Status:  http.StatusBadRequest,
		Error:   "bad_request",
	}
}
func NewNotFoundError(msg string) *RestError {
	return &RestError{
		Message: msg,
		Status:  http.StatusNotFound,
		Error:   "not_found",
	}
}
