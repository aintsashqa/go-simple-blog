package mysql_test

import (
	"context"
	"errors"
	"testing"

	"github.com/aintsashqa/go-simple-blog/internal/domain"
	"github.com/aintsashqa/go-simple-blog/internal/repository"
	"github.com/aintsashqa/go-simple-blog/internal/repository/mysql"
	mock_database "github.com/aintsashqa/go-simple-blog/pkg/database/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type UserRepositorySuite struct {
	suite.Suite
	*require.Assertions

	Controller *gomock.Controller

	MockDatabasePrivoder *mock_database.MockDatabasePrivoder

	CurrentRepository repository.User
}

func TestUserRepositorySuite(t *testing.T) {
	suite.Run(t, new(UserRepositorySuite))
}

func (s *UserRepositorySuite) SetupTest() {
	s.Assertions = require.New(s.T())
	s.Controller = gomock.NewController(s.T())
	s.MockDatabasePrivoder = mock_database.NewMockDatabasePrivoder(s.Controller)
	s.CurrentRepository = mysql.NewUserRepos(s.MockDatabasePrivoder)
}

func (s *UserRepositorySuite) TearDownTest() {
	s.Controller.Finish()
}

func (s *UserRepositorySuite) TestCreateMethod() {
	type MockDatabasePrivoderBehavior func(*mock_database.MockDatabasePrivoder, context.Context, error)

	mockDatabasePrivoderBehavior := func(m *mock_database.MockDatabasePrivoder, inputContext context.Context, returns error) {
		m.EXPECT().
			Exec(inputContext, gomock.Any(), gomock.Any()).
			Return(returns).
			Times(1)
	}

	databaseResultError := errors.New("DatabaseResultError")

	methodCases := []struct {
		Name                         string
		InputUser                    domain.User
		QueryString                  string
		DatabaseResultError          error
		MethodResultError            error
		MockDatabasePrivoderBehavior MockDatabasePrivoderBehavior
	}{
		{
			Name:                         "Success",
			InputUser:                    domain.User{},
			DatabaseResultError:          nil,
			MethodResultError:            nil,
			MockDatabasePrivoderBehavior: mockDatabasePrivoderBehavior,
		},
		{
			Name:                         "DatabaseFailure",
			InputUser:                    domain.User{},
			DatabaseResultError:          databaseResultError,
			MethodResultError:            databaseResultError,
			MockDatabasePrivoderBehavior: mockDatabasePrivoderBehavior,
		},
	}

	for _, currentCase := range methodCases {
		s.Suite.Run(currentCase.Name, func() {
			ctx := context.Background()
			currentCase.MockDatabasePrivoderBehavior(s.MockDatabasePrivoder, ctx, currentCase.DatabaseResultError)
			err := s.CurrentRepository.Create(ctx, currentCase.InputUser)
			if currentCase.DatabaseResultError != nil {
				s.Assertions.Equal(currentCase.MethodResultError, err)
			} else {
				s.Assertions.NoError(err)
			}
		})
	}
}

func (s *UserRepositorySuite) TestGetByEmailMethod() {
	type MockDatabasePrivoderBehavior func(*mock_database.MockDatabasePrivoder, context.Context, string, error)

	mockDatabasePrivoderBehavior := func(m *mock_database.MockDatabasePrivoder, inputContext context.Context, inputEmail string, returns error) {
		m.EXPECT().
			Get(inputContext, gomock.AssignableToTypeOf(&domain.User{}), gomock.Any(), gomock.Any()).
			Return(returns).
			Times(1).
			Do(func(_ context.Context, user *domain.User, _, _ string) error {
				user.Email = inputEmail
				return nil
			})
	}

	databaseResultError := errors.New("DatabaseResultError")

	methodCases := []struct {
		Name                         string
		InputEmail                   string
		DatabaseResultError          error
		MethodResultValue            domain.User
		MethodResultError            error
		MockDatabasePrivoderBehavior MockDatabasePrivoderBehavior
	}{
		{
			Name:                "Success",
			InputEmail:          "test@example.com",
			DatabaseResultError: nil,
			MethodResultValue: domain.User{
				Email: "test@example.com",
			},
			MethodResultError:            nil,
			MockDatabasePrivoderBehavior: mockDatabasePrivoderBehavior,
		},
		{
			Name:                "DatabaseFailure",
			InputEmail:          "test@example.com",
			DatabaseResultError: databaseResultError,
			MethodResultValue: domain.User{
				Email: "test@example.com",
			},
			MethodResultError:            databaseResultError,
			MockDatabasePrivoderBehavior: mockDatabasePrivoderBehavior,
		},
	}

	for _, currentCase := range methodCases {
		s.Suite.Run(currentCase.Name, func() {
			ctx := context.Background()
			currentCase.MockDatabasePrivoderBehavior(s.MockDatabasePrivoder, ctx, currentCase.InputEmail, currentCase.DatabaseResultError)
			result, err := s.CurrentRepository.GetByEmail(ctx, currentCase.InputEmail)
			s.Assertions.Equal(currentCase.MethodResultValue.Email, result.Email)
			s.Assertions.Equal(currentCase.MethodResultError, err)
		})
	}
}
