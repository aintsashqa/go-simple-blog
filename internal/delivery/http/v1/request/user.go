package request

import (
	"encoding/json"
	"net/http"

	"github.com/aintsashqa/go-simple-blog/internal/delivery/http/v1/errors"
	"github.com/aintsashqa/go-simple-blog/internal/delivery/http/v1/response"
	"github.com/aintsashqa/go-simple-blog/internal/service"
	"github.com/go-chi/chi"
	uuid "github.com/satori/go.uuid"
)

type SignUpUserRequestDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (dto *SignUpUserRequestDto) FromRequest(r *http.Request) (response.ErrorResponseDto, error) {
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response := response.NewErrorResponseDto(http.StatusInternalServerError, errors.ErrUnavailableRequestBody.Error())
		return response, errors.ErrUnavailableRequestBody
	}

	return response.ErrorResponseDto{}, nil
}

func (dto *SignUpUserRequestDto) TransformToObject() service.SignUpUserInput {
	return service.SignUpUserInput{
		Email:    dto.Email,
		Password: dto.Password,
	}
}

type SignInUserRequestDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (dto *SignInUserRequestDto) FromRequest(r *http.Request) (response.ErrorResponseDto, error) {
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response := response.NewErrorResponseDto(http.StatusInternalServerError, errors.ErrUnavailableRequestBody.Error())
		return response, errors.ErrUnavailableRequestBody
	}

	return response.ErrorResponseDto{}, nil
}

func (dto *SignInUserRequestDto) TransformToObject() service.SignInUserInput {
	return service.SignInUserInput{
		Email:    dto.Email,
		Password: dto.Password,
	}
}

type UpdateUserRequestDto struct {
	ID       uuid.UUID `json:"-"`
	Username string    `json:"username"`
}

func (dto *UpdateUserRequestDto) FromRequest(r *http.Request) (response.ErrorResponseDto, error) {
	authorizeUserID, casted := r.Context().Value("user_id").(uuid.UUID)
	if !casted {
		response := response.NewErrorResponseDto(http.StatusForbidden, errors.ErrInvalidTokenUserId.Error())
		return response, errors.ErrInvalidTokenUserId
	}

	urlUserID := uuid.FromStringOrNil(chi.URLParam(r, "id"))

	if authorizeUserID != urlUserID {
		response := response.NewErrorResponseDto(http.StatusForbidden, errors.ErrInvalidAuthorizedUserID.Error())
		return response, errors.ErrInvalidAuthorizedUserID
	}

	dto.ID = urlUserID

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response := response.NewErrorResponseDto(http.StatusInternalServerError, errors.ErrUnavailableRequestBody.Error())
		return response, errors.ErrUnavailableRequestBody
	}

	return response.ErrorResponseDto{}, nil
}

func (dto *UpdateUserRequestDto) TransformToObject() service.UpdateUserInput {
	return service.UpdateUserInput{
		ID:       dto.ID,
		Username: dto.Username,
	}
}
