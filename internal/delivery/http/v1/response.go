package v1

import (
	"net/http"

	"github.com/go-chi/render"
)

type responseError struct {
	Message string `json:"message"`
}

func errorFn(w http.ResponseWriter, r *http.Request, status int, message string) {
	respond(w, r, status, responseError{
		Message: message,
	})
}

func respond(w http.ResponseWriter, r *http.Request, status int, payload interface{}) {
	render.Status(r, status)
	render.JSON(w, r, payload)
}
