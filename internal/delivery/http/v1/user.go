package v1

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/aintsashqa/go-simple-blog/internal/service"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	uuid "github.com/satori/go.uuid"
)

type signUpRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *signUpRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Email, validation.Required, validation.Length(5, 255), is.Email),
		validation.Field(&r.Password, validation.Required, validation.Length(6, 255)),
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
		errorFn(w, r, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.user.SignUp(r.Context(), service.SignUpUserInput{
		Email:    input.Email,
		Password: input.Password,
	}); err != nil {
		log.Print(err)
		errorFn(w, r, http.StatusInternalServerError, err.Error())
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
		validation.Field(&r.Email, validation.Required, validation.Length(5, 255), is.Email),
		validation.Field(&r.Password, validation.Required, validation.Length(6, 255)),
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
		errorFn(w, r, http.StatusBadRequest, err.Error())
		return
	}

	tokens, err := h.user.SignIn(r.Context(), service.SignInUserInput{
		Email:    input.Email,
		Password: input.Password,
	})

	if err != nil {
		log.Print(err)
		errorFn(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	respond(w, r, http.StatusOK, signInResponse{
		AccessToken: tokens.AccessToken,
	})
}

func (h *Handler) authenticateMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if len(header) == 0 {
			log.Print("header `Authorization` is empty")
			errorFn(w, r, http.StatusUnauthorized, "header `Authorization` is empty")
			return
		}

		headerPieces := strings.Split(header, " ")
		if len(headerPieces) != 2 {
			log.Print("invalid `Authorization` header")
			errorFn(w, r, http.StatusUnauthorized, "invalid `Authorization` header")
			return
		}

		id, err := h.user.Authenticate(r.Context(), service.AuthenticateUserInput{
			Token: headerPieces[1],
		})
		if err != nil {
			log.Print(err)
			errorFn(w, r, http.StatusUnauthorized, err.Error())
			return
		}

		if id == uuid.Nil {
			log.Print("invalid user id")
			errorFn(w, r, http.StatusUnauthorized, "invalid user id")
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
