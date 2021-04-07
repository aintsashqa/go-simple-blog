package v1

import (
	"net/http"

	"github.com/aintsashqa/go-simple-blog/internal/delivery/http/v1/errors"
	"github.com/aintsashqa/go-simple-blog/internal/delivery/http/v1/response"
	"github.com/aintsashqa/go-simple-blog/internal/domain"
)

func ValidationErrorsHandler(err error) (response.ErrorResponseDto, bool) {
	switch err {

	case
		// User errors
		domain.ErrUserEmailEmptyValue,
		domain.ErrUserEmailInvalidLength,
		domain.ErrUserEmailInvalidValue,
		domain.ErrUserUsernameEmptyValue,
		domain.ErrUserUsernameInvalidLength,
		domain.ErrUserPasswordEmptyValue,
		domain.ErrUserPasswordInvalidLength,

		// Post errors
		domain.ErrPostTitleEmptyValue,
		domain.ErrPostTitleInvalidLength,
		domain.ErrPostSlugInvalidLength,
		domain.ErrPostContentEmptyValue,
		domain.ErrPostContentInvalidLength:

		return response.NewErrorResponseDto(http.StatusBadRequest, errors.ErrInvalidRequestBody.Error(), err.Error()), true
	}

	return response.ErrorResponseDto{}, false
}
