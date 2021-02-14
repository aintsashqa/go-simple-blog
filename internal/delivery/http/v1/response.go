package v1

import (
	"net/http"

	"github.com/go-chi/render"
)

const (
	internalErrorMsg             string = "Internal error"
	validationFailedMsg          string = "Validation failed"
	invalidUrlQueryParamErrorMsg string = "Invalid url query parameter `%s`"
	authorizationFailedMsg       string = "Authorization failed"
)

type responseError struct {
	Message string   `json:"message"`
	Errors  []string `json:"errors"`
}

func errorFn(w http.ResponseWriter, r *http.Request, status int, message string, errors ...string) {
	respond(w, r, status, responseError{
		Message: message,
		Errors:  errors,
	})
}

func respond(w http.ResponseWriter, r *http.Request, status int, payload interface{}) {
	render.Status(r, status)
	render.JSON(w, r, payload)
}
