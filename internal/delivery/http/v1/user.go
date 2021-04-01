package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/aintsashqa/go-simple-blog/internal/domain"
	"github.com/aintsashqa/go-simple-blog/internal/service"
	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	uuid "github.com/satori/go.uuid"
)

const (
	GetSingleUserCacheKey string = "get-single-user[id=%s]"
)

type signUpRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *signUpRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Email,
			validation.Required,
			validation.Length(5, 255),
			is.Email,
		),
		validation.Field(&r.Password,
			validation.Required,
			validation.Length(6, 255),
		),
	)
}

type signUpResponse struct {
	Message string `json:"message"`
}

// @Summary Sign up
// @Description Sign up with account details
// @ID user-sign-up
// @Tags User
// @Accept json
// @Produce json
// @Param payload body signUpRequest true "Sign up with account details"
// @Success 201 {object} signUpResponse
// @Failure 400 {object} responseError
// @Failure 500 {object} responseError
// @Router /user/sign-up [post]
func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {
	var input signUpRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Print(err)
		errorFn(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	if err := input.Validate(); err != nil {
		log.Print(err)
		errorFn(w, r, http.StatusBadRequest, validationFailedMsg, strings.Split(err.Error(), "; ")...)
		return
	}

	if err := h.user.SignUp(r.Context(), service.SignUpUserInput{
		Email:    input.Email,
		Password: input.Password,
	}); err != nil {
		log.Print(err)
		errorFn(w, r, http.StatusInternalServerError, internalErrorMsg, err.Error())
		return
	}

	respond(w, r, http.StatusCreated, signUpResponse{
		Message: "OK",
	})
}

type signInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *signInRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Email,
			validation.Required,
			validation.Length(5, 255),
			is.Email,
		),
		validation.Field(&r.Password,
			validation.Required,
			validation.Length(6, 255),
		),
	)
}

type signInResponse struct {
	AccessToken string `json:"access_token"`
}

// @Summary Sign in
// @Description Sign in with account details
// @ID user-sign-in
// @Tags User
// @Accept json
// @Produce json
// @Param payload body signInRequest true "Sign in with account details"
// @Success 200 {object} signInResponse
// @Failure 400 {object} responseError
// @Failure 500 {object} responseError
// @Router /user/sign-in [post]
func (h *Handler) signIn(w http.ResponseWriter, r *http.Request) {
	var input signInRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Print(err)
		errorFn(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	if err := input.Validate(); err != nil {
		log.Print(err)
		errorFn(w, r, http.StatusBadRequest, validationFailedMsg, strings.Split(err.Error(), "; ")...)
		return
	}

	tokens, err := h.user.SignIn(r.Context(), service.SignInUserInput{
		Email:    input.Email,
		Password: input.Password,
	})

	if err != nil {
		log.Print(err)
		errorFn(w, r, http.StatusInternalServerError, internalErrorMsg, err.Error())
		return
	}

	respond(w, r, http.StatusOK, signInResponse{
		AccessToken: tokens.AccessToken,
	})
}

type getSingleUserResponse struct {
	User domain.User `json:"user"`
}

// @Summary Get single user
// @Description Get single user by id
// @ID user-get-single
// @Tags User
// @Accept json
// @Produce json
// @Param id path string string "User with id"
// @Success 200 {object} getSingleUserResponse
// @Failure 500 {object} responseError
// @Router /user/{id} [get]
func (h *Handler) getSingleUser(w http.ResponseWriter, r *http.Request) {
	id := uuid.FromStringOrNil(chi.URLParam(r, "id"))

	var response getSingleUserResponse
	currentCacheKey := fmt.Sprintf(GetSingleUserCacheKey, id)

	founded, err := h.getFromCache(r.Context(), currentCacheKey, &response)
	if err != nil {
		log.Print(err)
		errorFn(w, r, http.StatusInternalServerError, internalErrorMsg, err.Error())
		return
	}

	if !founded {

		user, err := h.user.Find(r.Context(), id)
		if err != nil {
			log.Print(err)
			errorFn(w, r, http.StatusInternalServerError, internalErrorMsg, err.Error())
			return
		}

		response = getSingleUserResponse{
			User: user,
		}

		h.saveToCache(r.Context(), currentCacheKey, response)
	}

	respond(w, r, http.StatusOK, response)
}

type updateUserRequest struct {
	Username string `json:"username"`
}

func (r *updateUserRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Username, validation.Required, validation.Length(6, 255)),
	)
}

type updateUserResponse struct {
	User domain.User `json:"user"`
}

// @Summary Update user
// @Description Update user with id
// @ID user-update
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "User with id"
// @Param payload body updateUserRequest true "User details"
// @Success 202 {object} updateUserResponse
// @Failure 400 {object} responseError
// @Failure 401 {object} responseError
// @Failure 403 {object} responseError
// @Failure 500 {object} responseError
// @Security ApiKeyAuth
// @Router /user/{id} [put]
func (h *Handler) updateUser(w http.ResponseWriter, r *http.Request) {
	var input updateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Print(err)
		errorFn(w, r, http.StatusInternalServerError, internalErrorMsg, err.Error())
		return
	}

	if err := input.Validate(); err != nil {
		log.Print(err)
		errorFn(w, r, http.StatusBadRequest, validationFailedMsg, err.Error())
		return
	}

	authorizedUserID := r.Context().Value("user_id").(uuid.UUID)
	id := uuid.FromStringOrNil(chi.URLParam(r, "id"))
	if authorizedUserID != id {
		log.Print("invalid authorized user")
		errorFn(w, r, http.StatusForbidden, authorizationFailedMsg, "invalid authorized user")
		return
	}

	opt := service.UpdateUserInput{
		ID:       id,
		Username: input.Username,
	}

	user, err := h.user.Update(r.Context(), opt)
	if err != nil {
		log.Print(err)
		errorFn(w, r, http.StatusInternalServerError, internalErrorMsg, err.Error())
		return
	}

	respond(w, r, http.StatusAccepted, updateUserResponse{
		User: user,
	})
}

func (h *Handler) authenticateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if len(header) == 0 {
			log.Print("header `Authorization` is empty")
			errorFn(w, r, http.StatusUnauthorized, authorizationFailedMsg, "header `Authorization` is empty")
			return
		}

		headerPieces := strings.Split(header, " ")
		if len(headerPieces) != 2 {
			log.Print("invalid `Authorization` header")
			errorFn(w, r, http.StatusUnauthorized, authorizationFailedMsg, "invalid `Authorization` header")
			return
		}

		id, err := h.user.Authenticate(r.Context(), service.AuthenticateUserInput{
			Token: headerPieces[1],
		})
		if err != nil {
			log.Print(err)
			errorFn(w, r, http.StatusUnauthorized, authorizationFailedMsg, err.Error())
			return
		}

		if id == uuid.Nil {
			log.Print("invalid user id")
			errorFn(w, r, http.StatusUnauthorized, authorizationFailedMsg, "invalid user id")
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
