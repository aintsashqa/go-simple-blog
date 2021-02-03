package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aintsashqa/go-simple-blog/internal/domain"
	mock_repository "github.com/aintsashqa/go-simple-blog/internal/repository/mocks"
	"github.com/aintsashqa/go-simple-blog/internal/service"
	"github.com/aintsashqa/go-simple-blog/pkg/auth"
	mock_auth "github.com/aintsashqa/go-simple-blog/pkg/auth/mocks"
	mock_hash "github.com/aintsashqa/go-simple-blog/pkg/hash/mocks"
	"github.com/golang/mock/gomock"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type UserServiceSuite struct {
	suite.Suite
	*require.Assertions

	Controller *gomock.Controller

	MockUserRepository *mock_repository.MockUser
	MockHashProvider   *mock_hash.MockHashProvider
	MockAuthProvider   *mock_auth.MockAuthorizationProvider

	TokenExpiresTime time.Duration

	CurrentService service.User
}

func TestUserServiceSuite(t *testing.T) {
	suite.Run(t, new(UserServiceSuite))
}

func (s *UserServiceSuite) SetupTest() {
	s.Assertions = require.New(s.T())
	s.Controller = gomock.NewController(s.T())
	s.MockUserRepository = mock_repository.NewMockUser(s.Controller)
	s.MockHashProvider = mock_hash.NewMockHashProvider(s.Controller)
	s.MockAuthProvider = mock_auth.NewMockAuthorizationProvider(s.Controller)
	s.TokenExpiresTime = time.Duration(time.Hour * 10)
	s.CurrentService = service.NewUserService(s.MockUserRepository, s.MockHashProvider, s.MockAuthProvider, s.TokenExpiresTime)
}

func (s *UserServiceSuite) TearDownTest() {
	s.Controller.Finish()
}

func (s *UserServiceSuite) TestSignUpMethod() {
	type MockUserRepositoryBehavior func(m *mock_repository.MockUser, returns error)
	type MockHashProviderBehavior func(m *mock_hash.MockHashProvider, input string, returns string)

	mockUserRepositoryBehavior := func(m *mock_repository.MockUser, returns error) {
		m.EXPECT().
			Create(context.Background(), gomock.AssignableToTypeOf(domain.User{})).
			Return(returns).
			Times(1)
	}

	mockHashProviderBehavior := func(m *mock_hash.MockHashProvider, input string, returns string) {
		m.EXPECT().
			Make(input).
			Return(returns).
			Times(1)
	}

	methodCases := []struct {
		Name                       string
		PasswordHash               string
		ServiceInput               service.SignUpUserInput
		RepositoryResultError      error
		MockUserRepositoryBehavior MockUserRepositoryBehavior
		MockHashProviderBehavior   MockHashProviderBehavior
	}{
		{
			Name:         "Success",
			PasswordHash: "secret-hash",
			ServiceInput: service.SignUpUserInput{
				Email:    "test@example.com",
				Password: "secret",
			},
			RepositoryResultError:      nil,
			MockUserRepositoryBehavior: mockUserRepositoryBehavior,
			MockHashProviderBehavior:   mockHashProviderBehavior,
		},
		{
			Name:         "RepositoryFailure",
			PasswordHash: "secret-hash",
			ServiceInput: service.SignUpUserInput{
				Email:    "test@example.com",
				Password: "secret",
			},
			RepositoryResultError:      errors.New("RepositoryResultError"),
			MockUserRepositoryBehavior: mockUserRepositoryBehavior,
			MockHashProviderBehavior:   mockHashProviderBehavior,
		},
	}

	for _, currentCase := range methodCases {
		s.Suite.Run(currentCase.Name, func() {
			currentCase.MockUserRepositoryBehavior(s.MockUserRepository, currentCase.RepositoryResultError)
			currentCase.MockHashProviderBehavior(s.MockHashProvider, currentCase.ServiceInput.Password, currentCase.PasswordHash)
			err := s.CurrentService.SignUp(context.Background(), currentCase.ServiceInput)
			s.Assertions.Equal(currentCase.RepositoryResultError, err)
		})
	}
}

func (s *UserServiceSuite) TestSignInMethod() {
	type MockUserRepositoryBehavior func(m *mock_repository.MockUser, input service.SignInUserInput, returnsUser domain.User, returnsError error)
	type MockHashProviderBehavior func(m *mock_hash.MockHashProvider, inputHashPassword string, inputPassword string, returns error)
	type MockAuthProviderBehavior func(m *mock_auth.MockAuthorizationProvider, inputUserID uuid.UUID, inputExpiresAt time.Duration, returnsAccessToken string, returnsError error)

	mockUserRepositoryBehavior := func(m *mock_repository.MockUser, input service.SignInUserInput, returnsUser domain.User, returnsError error) {
		m.EXPECT().
			GetByEmail(context.Background(), input.Email).
			Return(returnsUser, returnsError).
			Times(1)
	}

	mockHashProviderBehavior := func(m *mock_hash.MockHashProvider, inputHashPassword string, inputPassword string, returns error) {
		m.EXPECT().
			Compare(inputHashPassword, inputPassword).
			Return(returns).
			Times(1)
	}

	mockAuthProviderBehavior := func(m *mock_auth.MockAuthorizationProvider, inputUserID uuid.UUID, inputExpiresAt time.Duration, returnsAccessToken string, returnsError error) {
		m.EXPECT().
			NewToken(auth.TokenParams{UserID: inputUserID, ExpiresAt: inputExpiresAt}).
			Return(returnsAccessToken, returnsError).
			Times(1)
	}

	repositoryResultError := errors.New("RepositoryResultError")
	hashResultError := errors.New("HashResultError")
	authResultError := errors.New("AuthResultError")

	methodCases := []struct {
		Name                       string
		ServiceInput               service.SignInUserInput
		CurrentUser                domain.User
		AccessToken                string
		RepositoryResultError      error
		HashResultError            error
		AuthResultError            error
		MethodResultValue          service.Tokens
		MethodResultError          error
		MockUserRepositoryBehavior MockUserRepositoryBehavior
		MockHashProviderBehavior   MockHashProviderBehavior
		MockAuthProviderBehavior   MockAuthProviderBehavior
	}{
		{
			Name: "Success",
			ServiceInput: service.SignInUserInput{
				Email:    "test@example.com",
				Password: "secret",
			},
			CurrentUser: domain.User{
				ID:       uuid.NewV4(),
				Password: "secret-hash",
			},
			AccessToken:           "access-token-valid",
			RepositoryResultError: nil,
			HashResultError:       nil,
			AuthResultError:       nil,
			MethodResultValue: service.Tokens{
				AccessToken: "access-token-valid",
			},
			MethodResultError:          nil,
			MockUserRepositoryBehavior: mockUserRepositoryBehavior,
			MockHashProviderBehavior:   mockHashProviderBehavior,
			MockAuthProviderBehavior:   mockAuthProviderBehavior,
		},
		{
			Name: "RepositoryFailure",
			ServiceInput: service.SignInUserInput{
				Email:    "test@example.com",
				Password: "secret",
			},
			CurrentUser: domain.User{
				ID:       uuid.NewV4(),
				Password: "secret-hash",
			},
			RepositoryResultError:      repositoryResultError,
			HashResultError:            nil,
			AuthResultError:            nil,
			MethodResultValue:          service.Tokens{},
			MethodResultError:          repositoryResultError,
			MockUserRepositoryBehavior: mockUserRepositoryBehavior,
			MockHashProviderBehavior:   nil,
			MockAuthProviderBehavior:   nil,
		},
		{
			Name: "HasherFailure",
			ServiceInput: service.SignInUserInput{
				Email:    "test@example.com",
				Password: "secret1",
			},
			CurrentUser: domain.User{
				ID:       uuid.NewV4(),
				Password: "secret-hash",
			},
			RepositoryResultError:      nil,
			HashResultError:            hashResultError,
			AuthResultError:            nil,
			MethodResultValue:          service.Tokens{},
			MethodResultError:          hashResultError,
			MockUserRepositoryBehavior: mockUserRepositoryBehavior,
			MockHashProviderBehavior:   mockHashProviderBehavior,
			MockAuthProviderBehavior:   nil,
		},
		{
			Name: "AuthFailure",
			ServiceInput: service.SignInUserInput{
				Email:    "test@example.com",
				Password: "secret",
			},
			CurrentUser: domain.User{
				ID:       uuid.NewV4(),
				Password: "secret-hash",
			},
			RepositoryResultError:      nil,
			HashResultError:            nil,
			AuthResultError:            authResultError,
			MethodResultValue:          service.Tokens{},
			MethodResultError:          authResultError,
			MockUserRepositoryBehavior: mockUserRepositoryBehavior,
			MockHashProviderBehavior:   mockHashProviderBehavior,
			MockAuthProviderBehavior:   mockAuthProviderBehavior,
		},
	}

	for _, currentCase := range methodCases {
		s.Suite.Run(currentCase.Name, func() {
			if currentCase.MockUserRepositoryBehavior != nil {
				currentCase.MockUserRepositoryBehavior(s.MockUserRepository, currentCase.ServiceInput, currentCase.CurrentUser, currentCase.RepositoryResultError)
			}
			if currentCase.MockHashProviderBehavior != nil {
				currentCase.MockHashProviderBehavior(s.MockHashProvider, currentCase.CurrentUser.Password, currentCase.ServiceInput.Password, currentCase.HashResultError)
			}
			if currentCase.MockAuthProviderBehavior != nil {
				currentCase.MockAuthProviderBehavior(s.MockAuthProvider, currentCase.CurrentUser.ID, s.TokenExpiresTime, currentCase.AccessToken, currentCase.AuthResultError)
			}
			result, err := s.CurrentService.SignIn(context.Background(), currentCase.ServiceInput)
			s.Assertions.Equal(result, currentCase.MethodResultValue)
			if currentCase.RepositoryResultError != nil {
				s.Assertions.Equal(currentCase.RepositoryResultError, err)
			} else if currentCase.HashResultError != nil {
				s.Assertions.Equal(currentCase.HashResultError, err)
			} else if currentCase.AuthResultError != nil {
				s.Assertions.Equal(currentCase.AuthResultError, err)
			} else {
				s.Assertions.NoError(err)
			}
		})
	}
}

func (s *UserServiceSuite) TestAuthenticateMethod() {
	type MockAuthProviderBehavior func(m *mock_auth.MockAuthorizationProvider, input string, returnsID uuid.UUID, returnsError error)

	mockAuthProviderBehavior := func(m *mock_auth.MockAuthorizationProvider, input string, returnsID uuid.UUID, returnsError error) {
		m.EXPECT().
			Parse(input).
			Return(returnsID, returnsError).
			Times(1)
	}

	authResultError := errors.New("AuthResultError")
	id := uuid.NewV4()

	methodCases := []struct {
		Name                     string
		ServiceInput             service.AuthenticateUserInput
		InputToken               string
		AuthResultError          error
		MethodResultValue        uuid.UUID
		MethodResultError        error
		MockAuthProviderBehavior MockAuthProviderBehavior
	}{
		{
			Name: "Success",
			ServiceInput: service.AuthenticateUserInput{
				Token: "access-token-valid",
			},
			InputToken:               "access-token-valid",
			AuthResultError:          nil,
			MethodResultValue:        id,
			MethodResultError:        nil,
			MockAuthProviderBehavior: mockAuthProviderBehavior,
		},
		{
			Name:                     "AuthFailure",
			ServiceInput:             service.AuthenticateUserInput{},
			AuthResultError:          authResultError,
			MethodResultValue:        uuid.Nil,
			MethodResultError:        authResultError,
			MockAuthProviderBehavior: mockAuthProviderBehavior,
		},
	}

	for _, currentCase := range methodCases {
		s.Suite.Run(currentCase.Name, func() {
			currentCase.MockAuthProviderBehavior(s.MockAuthProvider, currentCase.InputToken, currentCase.MethodResultValue, currentCase.MethodResultError)
			result, err := s.CurrentService.Authenticate(context.Background(), currentCase.ServiceInput)
			s.Assertions.Equal(currentCase.MethodResultValue, result)
			s.Assertions.Equal(currentCase.MethodResultError, err)
		})
	}
}
