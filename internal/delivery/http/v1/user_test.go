package v1_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	v1 "github.com/aintsashqa/go-simple-blog/internal/delivery/http/v1"
	rerr "github.com/aintsashqa/go-simple-blog/internal/delivery/http/v1/errors"
	"github.com/aintsashqa/go-simple-blog/internal/domain"
	derr "github.com/aintsashqa/go-simple-blog/internal/domain"
	repoerr "github.com/aintsashqa/go-simple-blog/internal/repository/errors"
	"github.com/aintsashqa/go-simple-blog/internal/service"
	mock_service "github.com/aintsashqa/go-simple-blog/internal/service/mocks"
	"github.com/golang/mock/gomock"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	ErrorResponseBodyInformationNull string = `{"code":%d,"message":"%s","information":null}`
	ErrorResponseBody                string = `{"code":%d,"message":"%s","information":["%s"]}`
	SignUpResponseBody               string = `{"id":"%s","email":"%s","username":"new username","created_at":"%v","updated_at":"%v"}`
	SignInResponseBody               string = `{"access_token":"%s"}`
)

type UserHTTPHandlerSuite struct {
	suite.Suite
	*require.Assertions

	Controller *gomock.Controller

	MockUserService *mock_service.MockUser

	CurrentHTTPHandler *v1.Handler
}

func TestUserHTTPHandlerSuite(t *testing.T) {
	suite.Run(t, new(UserHTTPHandlerSuite))
}

func (s *UserHTTPHandlerSuite) SetupTest() {
	s.Assertions = require.New(s.T())
	s.Controller = gomock.NewController(s.T())
	s.MockUserService = mock_service.NewMockUser(s.Controller)
	service := service.Service{User: s.MockUserService}
	s.CurrentHTTPHandler = v1.NewHandler(&service)
}

func (s *UserHTTPHandlerSuite) TearDownTest() {
	s.Controller.Finish()
}

func (s *UserHTTPHandlerSuite) TestSignUpMethod() {
	type MockUserServiceSignUpMethodBehavior func(*mock_service.MockUser, context.Context, service.SignUpUserInput, domain.User, error)

	mockUserServiceSignUpMethodBehavior := func(m *mock_service.MockUser, ctx context.Context, input service.SignUpUserInput, returnsUser domain.User, returnsError error) {
		m.EXPECT().
			SignUp(context.Background(), input).
			Return(returnsUser, returnsError).
			Times(1)
	}

	someInternalError := errors.New("SomeInternalError")
	id := uuid.NewV4()
	now := time.Now()

	methodCases := []struct {
		Name                                string
		RequestBody                         string
		ServiceInput                        service.SignUpUserInput
		ServiceResult                       domain.User
		ServiceResultError                  error
		MockUserServiceSignUpMethodBehavior MockUserServiceSignUpMethodBehavior
		ResponseBody                        string
		ResponseStatusCode                  int
	}{
		{
			Name:        "Success",
			RequestBody: `{"email":"root@example.com","password":"secret"}`,
			ServiceInput: service.SignUpUserInput{
				Email:    "root@example.com",
				Password: "secret",
			},
			ServiceResult: domain.User{
				Model: domain.Model{
					ID:        id,
					CreatedAt: now,
					UpdatedAt: now,
				},
				Email:    "root@example.com",
				Username: "new username",
				Password: "hashed",
			},
			ServiceResultError:                  nil,
			MockUserServiceSignUpMethodBehavior: mockUserServiceSignUpMethodBehavior,
			ResponseBody:                        fmt.Sprintf(SignUpResponseBody, id, "root@example.com", now.Format(time.RFC3339Nano), now.Format(time.RFC3339Nano)),
			ResponseStatusCode:                  http.StatusCreated,
		},
		{
			Name:                                "RequestFailure - Unavailable request body",
			RequestBody:                         `{"email":"root@example.com","password":"secret"`,
			ServiceInput:                        service.SignUpUserInput{},
			ServiceResult:                       domain.User{},
			ServiceResultError:                  nil,
			MockUserServiceSignUpMethodBehavior: nil,
			ResponseBody:                        fmt.Sprintf(ErrorResponseBodyInformationNull, http.StatusInternalServerError, rerr.ErrUnavailableRequestBody.Error()),
			ResponseStatusCode:                  http.StatusInternalServerError,
		},
		{
			Name:        "RequestFailure - Invalid request body",
			RequestBody: `{"email":"root@example.com"}`,
			ServiceInput: service.SignUpUserInput{
				Email:    "root@example.com",
				Password: "",
			},
			ServiceResult:                       domain.User{},
			ServiceResultError:                  derr.ErrUserPasswordEmptyValue,
			MockUserServiceSignUpMethodBehavior: mockUserServiceSignUpMethodBehavior,
			ResponseBody:                        fmt.Sprintf(ErrorResponseBody, http.StatusBadRequest, rerr.ErrInvalidRequestBody.Error(), derr.ErrUserPasswordEmptyValue.Error()),
			ResponseStatusCode:                  http.StatusBadRequest,
		},
		{
			Name:        "InternalFailure",
			RequestBody: `{"email":"root@example.com","password":"secret"}`,
			ServiceInput: service.SignUpUserInput{
				Email:    "root@example.com",
				Password: "secret",
			},
			ServiceResult:                       domain.User{},
			ServiceResultError:                  someInternalError,
			MockUserServiceSignUpMethodBehavior: mockUserServiceSignUpMethodBehavior,
			ResponseBody:                        fmt.Sprintf(ErrorResponseBodyInformationNull, http.StatusInternalServerError, rerr.ErrInternal.Error()),
			ResponseStatusCode:                  http.StatusInternalServerError,
		},
	}

	ctx := context.Background()
	for _, currentCase := range methodCases {
		s.Suite.Run(currentCase.Name, func() {
			handler := http.HandlerFunc(s.CurrentHTTPHandler.SignUp)
			if currentCase.MockUserServiceSignUpMethodBehavior != nil {
				currentCase.MockUserServiceSignUpMethodBehavior(s.MockUserService, ctx, currentCase.ServiceInput, currentCase.ServiceResult, currentCase.ServiceResultError)
			}
			responseRecorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodPost, "/v1/user/sign-up", bytes.NewBufferString(currentCase.RequestBody))
			handler.ServeHTTP(responseRecorder, request)
			s.Assertions.Equal(currentCase.ResponseStatusCode, responseRecorder.Code)
			s.Assertions.Equal(currentCase.ResponseBody+"\n", responseRecorder.Body.String())
		})
	}
}

func (s *UserHTTPHandlerSuite) TestSignInMethod() {
	type MockUserServiceSignInMethodBehavior func(*mock_service.MockUser, context.Context, service.SignInUserInput, service.Tokens, error)

	mockUserServiceSignInMethodBehavior := func(m *mock_service.MockUser, ctx context.Context, input service.SignInUserInput, returnsTokens service.Tokens, returnsError error) {
		m.EXPECT().
			SignIn(context.Background(), input).
			Return(returnsTokens, returnsError).
			Times(1)
	}

	someInternalError := errors.New("SomeInternalError")

	methodCases := []struct {
		Name                                string
		RequestBody                         string
		ServiceInput                        service.SignInUserInput
		ServiceResult                       service.Tokens
		ServiceResultError                  error
		MockUserServiceSignInMethodBehavior MockUserServiceSignInMethodBehavior
		ResponseBody                        string
		ResponseStatusCode                  int
	}{
		{
			Name:        "Success",
			RequestBody: fmt.Sprintf(`{"email":"root@example.com","password":"secret"}`),
			ServiceInput: service.SignInUserInput{
				Email:    "root@example.com",
				Password: "secret",
			},
			ServiceResult: service.Tokens{
				AccessToken: "VALID_ACCESS_TOKEN",
			},
			ServiceResultError:                  nil,
			MockUserServiceSignInMethodBehavior: mockUserServiceSignInMethodBehavior,
			ResponseBody:                        fmt.Sprintf(SignInResponseBody, "VALID_ACCESS_TOKEN"),
			ResponseStatusCode:                  http.StatusOK,
		},
		{
			Name:                                "RequestFailure - Unavailable request body",
			RequestBody:                         fmt.Sprintf(`{"email":"root@example.com","password":"secret"`),
			ServiceInput:                        service.SignInUserInput{},
			ServiceResult:                       service.Tokens{},
			ServiceResultError:                  nil,
			MockUserServiceSignInMethodBehavior: nil,
			ResponseBody:                        fmt.Sprintf(ErrorResponseBodyInformationNull, http.StatusInternalServerError, rerr.ErrUnavailableRequestBody.Error()),
			ResponseStatusCode:                  http.StatusInternalServerError,
		},
		{
			Name:        "RepositoryFailure",
			RequestBody: fmt.Sprintf(`{"email":"notfound@example.com","password":"secret"}`),
			ServiceInput: service.SignInUserInput{
				Email:    "notfound@example.com",
				Password: "secret",
			},
			ServiceResult:                       service.Tokens{},
			ServiceResultError:                  repoerr.ErrUserNotFound,
			MockUserServiceSignInMethodBehavior: mockUserServiceSignInMethodBehavior,
			ResponseBody:                        fmt.Sprintf(ErrorResponseBodyInformationNull, http.StatusNotFound, repoerr.ErrUserNotFound.Error()),
			ResponseStatusCode:                  http.StatusNotFound,
		},
		{
			Name:        "InternalFailure",
			RequestBody: fmt.Sprintf(`{"email":"notfound@example.com","password":"secret"}`),
			ServiceInput: service.SignInUserInput{
				Email:    "notfound@example.com",
				Password: "secret",
			},
			ServiceResult:                       service.Tokens{},
			ServiceResultError:                  someInternalError,
			MockUserServiceSignInMethodBehavior: mockUserServiceSignInMethodBehavior,
			ResponseBody:                        fmt.Sprintf(ErrorResponseBodyInformationNull, http.StatusInternalServerError, rerr.ErrInternal.Error()),
			ResponseStatusCode:                  http.StatusInternalServerError,
		},
	}

	ctx := context.Background()
	for _, currentCase := range methodCases {
		s.Suite.Run(currentCase.Name, func() {
			handler := http.HandlerFunc(s.CurrentHTTPHandler.SignIn)
			if currentCase.MockUserServiceSignInMethodBehavior != nil {
				currentCase.MockUserServiceSignInMethodBehavior(s.MockUserService, ctx, currentCase.ServiceInput, currentCase.ServiceResult, currentCase.ServiceResultError)
			}
			responseRecorder := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodPost, "/v1/user/sign-ip", bytes.NewBufferString(currentCase.RequestBody))
			handler.ServeHTTP(responseRecorder, request)
			s.Assertions.Equal(currentCase.ResponseStatusCode, responseRecorder.Code)
			s.Assertions.Equal(currentCase.ResponseBody+"\n", responseRecorder.Body.String())
		})
	}
}
