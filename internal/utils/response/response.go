package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status       string `json:"status"`
	ErrorMessage string `json:"error"`
}

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

func WriteJson(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)
}

func GeneralError(err error) Response {
	return Response{
		Status:       StatusError,
		ErrorMessage: err.Error(),
	}
}

func ValidationError(errors validator.ValidationErrors) Response {
	var errMessages []string

	for _, err := range errors {
		switch err.ActualTag() {
		case "required":
			errMessages = append(errMessages, fmt.Sprintf("field %s is required", err.Field()))
		default:
			errMessages = append(errMessages, fmt.Sprintf("field %s is invalid", err.Field()))
		}
	}

	return Response{
		Status:       StatusError,
		ErrorMessage: strings.Join(errMessages, ", "),
	}
}
