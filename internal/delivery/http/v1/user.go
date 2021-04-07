package v1

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/aintsashqa/go-simple-blog/internal/delivery/http/v1/errors"
	requestdto "github.com/aintsashqa/go-simple-blog/internal/delivery/http/v1/request"
	responsedto "github.com/aintsashqa/go-simple-blog/internal/delivery/http/v1/response"
	repoerrors "github.com/aintsashqa/go-simple-blog/internal/repository/errors"
	"github.com/aintsashqa/go-simple-blog/internal/service"
	"github.com/go-chi/chi"
	uuid "github.com/satori/go.uuid"
)

// @Summary Sign up
// @Description Sign up with account details
// @ID user-sign-up
// @Tags User
// @Accept json
// @Produce json
// @Param payload body request.SignUpUserRequestDto true "Sign up with account details"
// @Success 201 {object} response.UserResponseDto
// @Failure 400 {object} response.ErrorResponseDto
// @Failure 500 {object} response.ErrorResponseDto
// @Router /user/sign-up [post]
func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	request := requestdto.SignUpUserRequestDto{}
	response := responsedto.UserResponseDto{}

	if response, err := request.FromRequest(r); err != nil {

		log.Print(err)

		errorRespond(w, r, response)
		return
	}

	opt := request.TransformToObject()
	user, err := h.user.SignUp(r.Context(), opt)
	if err != nil {

		log.Print(err)

		if errorResp, isValidation := ValidationErrorsHandler(err); isValidation {
			errorRespond(w, r, errorResp)
			return
		}

		errorResp := responsedto.NewErrorResponseDto(http.StatusInternalServerError, errors.ErrInternal.Error())
		errorRespond(w, r, errorResp)
		return
	}

	response.TransformFromObject(user)
	respond(w, r, http.StatusCreated, response)
}

// @Summary Sign in
// @Description Sign in with account details
// @ID user-sign-in
// @Tags User
// @Accept json
// @Produce json
// @Param payload body request.SignInUserRequestDto true "Sign in with account details"
// @Success 200 {object} response.TokenResponseDto
// @Failure 404 {object} response.ErrorResponseDto
// @Failure 500 {object} response.ErrorResponseDto
// @Router /user/sign-in [post]
func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	request := requestdto.SignInUserRequestDto{}
	response := responsedto.TokenResponseDto{}

	if response, err := request.FromRequest(r); err != nil {

		log.Print(err)

		errorRespond(w, r, response)
		return
	}

	opt := request.TransformToObject()
	tokens, err := h.user.SignIn(r.Context(), opt)
	if err != nil {

		log.Print(err)

		// TODO: check passwords not equals to send StatusBadRequest
		var errorResp responsedto.ErrorResponseDto
		if err == repoerrors.ErrUserNotFound {
			errorResp = responsedto.NewErrorResponseDto(http.StatusNotFound, err.Error())
		} else {
			errorResp = responsedto.NewErrorResponseDto(http.StatusInternalServerError, errors.ErrInternal.Error())

		}

		errorRespond(w, r, errorResp)
		return
	}

	response.TransformFromObject(tokens)
	respond(w, r, http.StatusOK, response)
}

// @Summary Get self user
// @Description Get self user by authorized information
// @ID user-get-self
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} response.UserResponseDto
// @Failure 403 {object} response.ErrorResponseDto
// @Failure 404 {object} response.ErrorResponseDto
// @Failure 500 {object} response.ErrorResponseDto
// @Security ApiKeyAuth
// @Router /user/self [get]
func (h *Handler) GetSelfUser(w http.ResponseWriter, r *http.Request) {
	request := requestdto.SelfUserRequestDto{}
	response := responsedto.UserResponseDto{}

	if response, err := request.FromRequest(r); err != nil {

		log.Print(err)

		errorRespond(w, r, response)
		return
	}

	user, err := h.user.Self(r.Context(), request.ID)
	if err != nil {

		log.Print(err)

		var errorResp responsedto.ErrorResponseDto
		if err == repoerrors.ErrUserNotFound {
			errorResp = responsedto.NewErrorResponseDto(http.StatusNotFound, err.Error())
		} else {
			errorResp = responsedto.NewErrorResponseDto(http.StatusInternalServerError, errors.ErrInternal.Error())
		}

		errorRespond(w, r, errorResp)
		return
	}

	response.TransformFromObject(user)
	respond(w, r, http.StatusOK, response)
}

// @Summary Get single user
// @Description Get single user by id
// @ID user-get-single
// @Tags User
// @Accept json
// @Produce json
// @Param id path string string "User with id"
// @Success 200 {object} response.UserResponseDto
// @Failure 404 {object} response.ErrorResponseDto
// @Failure 500 {object} response.ErrorResponseDto
// @Router /user/{id} [get]
func (h *Handler) GetSingleUser(w http.ResponseWriter, r *http.Request) {
	response := responsedto.UserResponseDto{}

	id := uuid.FromStringOrNil(chi.URLParam(r, "id"))
	user, err := h.user.Find(r.Context(), id)
	if err != nil {

		log.Print(err)

		var errorResp responsedto.ErrorResponseDto
		if err == repoerrors.ErrUserNotFound {
			errorResp = responsedto.NewErrorResponseDto(http.StatusNotFound, err.Error())
		} else {
			errorResp = responsedto.NewErrorResponseDto(http.StatusInternalServerError, errors.ErrInternal.Error())
		}

		errorRespond(w, r, errorResp)
		return
	}

	response.TransformFromObject(user)
	respond(w, r, http.StatusOK, response)
}

// @Summary Update user
// @Description Update user with id
// @ID user-update
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "User with id"
// @Param payload body request.UpdateUserRequestDto true "User details"
// @Success 202 {object} response.UserResponseDto
// @Failure 400 {object} response.ErrorResponseDto
// @Failure 401 {object} response.ErrorResponseDto
// @Failure 403 {object} response.ErrorResponseDto
// @Failure 404 {object} response.ErrorResponseDto
// @Failure 500 {object} response.ErrorResponseDto
// @Security ApiKeyAuth
// @Router /user/{id} [put]
func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	request := requestdto.UpdateUserRequestDto{}
	response := responsedto.UserResponseDto{}

	if response, err := request.FromRequest(r); err != nil {

		log.Print(err)

		errorRespond(w, r, response)
		return
	}

	opt := request.TransformToObject()
	user, err := h.user.Update(r.Context(), opt)
	if err != nil {

		log.Print(err)

		if errorResp, isValidation := ValidationErrorsHandler(err); isValidation {
			errorRespond(w, r, errorResp)
			return
		}

		var errorResp responsedto.ErrorResponseDto
		if err == repoerrors.ErrUserNotFound {
			errorResp = responsedto.NewErrorResponseDto(http.StatusNotFound, err.Error())
		} else {
			errorResp = responsedto.NewErrorResponseDto(http.StatusInternalServerError, errors.ErrInternal.Error())
		}

		errorRespond(w, r, errorResp)
		return
	}

	response.TransformFromObject(user)
	respond(w, r, http.StatusAccepted, response)
}

func (h *Handler) authenticateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if len(header) == 0 {

			log.Print(errors.ErrEmptyAuthorizationHeader.Error())

			errorResp := responsedto.NewErrorResponseDto(http.StatusUnauthorized, errors.ErrEmptyAuthorizationHeader.Error())
			errorRespond(w, r, errorResp)
			return
		}

		headerPieces := strings.Split(header, " ")
		if len(headerPieces) != 2 {

			log.Print(errors.ErrInvalidAuthorizationHeader.Error())

			errorResp := responsedto.NewErrorResponseDto(http.StatusUnauthorized, errors.ErrInvalidAuthorizationHeader.Error())
			errorRespond(w, r, errorResp)
			return
		}

		id, err := h.user.Authenticate(r.Context(), service.AuthenticateUserInput{Token: headerPieces[1]})
		if err != nil {

			log.Print(err)

			errorResp := responsedto.NewErrorResponseDto(http.StatusUnauthorized, errors.ErrAuthenticationFailed.Error())
			errorRespond(w, r, errorResp)
			return
		}

		if id == uuid.Nil {

			log.Print(errors.ErrInvalidTokenUserId)

			errorResp := responsedto.NewErrorResponseDto(http.StatusUnauthorized, errors.ErrInvalidTokenUserId.Error())
			errorRespond(w, r, errorResp)
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
