package v1

import (
	"net/http"

	"github.com/aintsashqa/go-simple-blog/internal/delivery/http/v1/response"
	"github.com/go-chi/render"
)

const (
	internalErrorMsg             string = "Internal error"
	validationFailedMsg          string = "Validation failed"
	invalidUrlQueryParamErrorMsg string = "Invalid url query parameter `%s`"
	authorizationFailedMsg       string = "Authorization failed"
)

func respond(w http.ResponseWriter, r *http.Request, status int, payload interface{}) {
	render.Status(r, status)
	render.JSON(w, r, payload)
}

func errorRespond(w http.ResponseWriter, r *http.Request, err response.ErrorResponseDto) {
	respond(w, r, err.Code, err)
}
