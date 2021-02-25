package mysql_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aintsashqa/go-simple-blog/internal/domain"
	"github.com/aintsashqa/go-simple-blog/internal/repository"
	"github.com/aintsashqa/go-simple-blog/internal/repository/mysql"
	mock_database "github.com/aintsashqa/go-simple-blog/pkg/database/mocks"
	"github.com/golang/mock/gomock"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gopkg.in/guregu/null.v4"
)

type PostRepositorySuite struct {
	suite.Suite
	*require.Assertions

	Controller *gomock.Controller

	MockDatabasePrivoder *mock_database.MockDatabasePrivoder

	CurrentRepository repository.Post
}

func TestPostRepositorySuite(t *testing.T) {
	suite.Run(t, new(PostRepositorySuite))
}

func (s *PostRepositorySuite) SetupTest() {
	s.Assertions = require.New(s.T())
	s.Controller = gomock.NewController(s.T())
	s.MockDatabasePrivoder = mock_database.NewMockDatabasePrivoder(s.Controller)
	s.CurrentRepository = mysql.NewPostRepos(s.MockDatabasePrivoder)
}

func (s *PostRepositorySuite) TearDownTest() {
	s.Controller.Finish()
}

func (s *PostRepositorySuite) TestFindMethod() {
	type MockDatabasePrivoderBehavior func(*mock_database.MockDatabasePrivoder, context.Context, uuid.UUID, error)

	mockDatabasePrivoderBehavior := func(m *mock_database.MockDatabasePrivoder, inputContext context.Context, inputID uuid.UUID, returns error) {
		m.EXPECT().
			Get(inputContext, gomock.AssignableToTypeOf(&domain.Post{}), gomock.Any(), gomock.Any()).
			Return(returns).
			Times(1).
			Do(func(_ context.Context, post *domain.Post, _ string, _ uuid.UUID) error {
				post.ID = inputID
				return nil
			})
	}

	databaseResultError := errors.New("DatabaseResultError")

	id := uuid.NewV4()

	methodCases := []struct {
		Name                         string
		InputID                      uuid.UUID
		DatabaseResultError          error
		MethodResultValue            domain.Post
		MethodResultError            error
		MockDatabasePrivoderBehavior MockDatabasePrivoderBehavior
	}{
		{
			Name:                "Success",
			InputID:             id,
			DatabaseResultError: nil,
			MethodResultValue: domain.Post{
				Model: domain.Model{
					ID: id,
				},
			},
			MethodResultError:            nil,
			MockDatabasePrivoderBehavior: mockDatabasePrivoderBehavior,
		},
		{
			Name:                         "DatabaseFailure",
			InputID:                      uuid.Nil,
			DatabaseResultError:          databaseResultError,
			MethodResultValue:            domain.Post{},
			MethodResultError:            databaseResultError,
			MockDatabasePrivoderBehavior: mockDatabasePrivoderBehavior,
		},
	}

	for _, currentCase := range methodCases {
		s.Suite.Run(currentCase.Name, func() {
			ctx := context.Background()
			currentCase.MockDatabasePrivoderBehavior(s.MockDatabasePrivoder, ctx, currentCase.InputID, currentCase.DatabaseResultError)
			result, err := s.CurrentRepository.Find(ctx, currentCase.InputID)
			s.Assertions.Equal(currentCase.MethodResultValue.ID, result.ID)
			s.Assertions.Equal(currentCase.MethodResultError, err)
		})
	}
}

func (s *PostRepositorySuite) TestGetAllPublishedMethod() {
	type MockDatabasePrivoderBehavior func(*mock_database.MockDatabasePrivoder, context.Context, int, error)

	mockDatabasePrivoderBehavior := func(m *mock_database.MockDatabasePrivoder, inputContext context.Context, inputCount int, returns error) {
		m.EXPECT().
			Select(inputContext, gomock.AssignableToTypeOf(&[]domain.Post{}), gomock.Any(), gomock.Any(), inputCount).
			Return(returns).
			Times(1).
			Do(func(_ context.Context, posts *[]domain.Post, _ string, _ int, count int) error {
				if count > 0 {
					for i := 0; i < count; i++ {
						*posts = append(*posts, domain.Post{})
					}
				}
				return nil
			})
	}

	databaseResultError := errors.New("DatabaseResultError")

	methodCases := []struct {
		Name                         string
		InputOffset                  int
		InputCount                   int
		DatabaseResultError          error
		MethodResultError            error
		MockDatabasePrivoderBehavior MockDatabasePrivoderBehavior
	}{
		{
			Name:                         "Success",
			InputOffset:                  0,
			InputCount:                   10,
			DatabaseResultError:          nil,
			MethodResultError:            nil,
			MockDatabasePrivoderBehavior: mockDatabasePrivoderBehavior,
		},
		{
			Name:                         "SuccessEmpty",
			InputOffset:                  0,
			InputCount:                   0,
			DatabaseResultError:          nil,
			MethodResultError:            nil,
			MockDatabasePrivoderBehavior: mockDatabasePrivoderBehavior,
		},
		{
			Name:                         "DatabaseFailure",
			InputOffset:                  0,
			InputCount:                   10,
			DatabaseResultError:          databaseResultError,
			MethodResultError:            databaseResultError,
			MockDatabasePrivoderBehavior: mockDatabasePrivoderBehavior,
		},
	}

	for _, currentCase := range methodCases {
		s.Suite.Run(currentCase.Name, func() {
			ctx := context.Background()
			currentCase.MockDatabasePrivoderBehavior(s.MockDatabasePrivoder, ctx, currentCase.InputCount, currentCase.DatabaseResultError)
			result, err := s.CurrentRepository.GetAllPublished(ctx, currentCase.InputOffset, currentCase.InputCount)
			s.Assertions.NotNil(result)
			s.Assertions.Len(result, currentCase.InputCount)
			s.Assertions.Equal(currentCase.MethodResultError, err)
		})
	}
}

// same as GetAllPublished?
// func (s *PostRepositorySuite) TestGetAllPublishedWithUserIDMethod() {}

func (s *PostRepositorySuite) TestCreateMethod() {
	type MockDatabasePrivoderBehavior func(*mock_database.MockDatabasePrivoder, context.Context, error)

	mockDatabasePrivoderBehavior := func(m *mock_database.MockDatabasePrivoder, input context.Context, returns error) {
		m.EXPECT().
			Exec(input, gomock.Any(), gomock.Any()).
			Return(returns).
			Times(1)
	}

	databaseResultError := errors.New("DatabaseResultError")

	methodCases := []struct {
		Name                         string
		InputPost                    domain.Post
		DatabaseResultError          error
		MethodResultError            error
		MockDatabasePrivoderBehavior MockDatabasePrivoderBehavior
	}{
		{
			Name:                         "Success",
			InputPost:                    domain.Post{},
			DatabaseResultError:          nil,
			MethodResultError:            nil,
			MockDatabasePrivoderBehavior: mockDatabasePrivoderBehavior,
		},
		{
			Name:                         "DatabaseFailure",
			InputPost:                    domain.Post{},
			DatabaseResultError:          databaseResultError,
			MethodResultError:            databaseResultError,
			MockDatabasePrivoderBehavior: mockDatabasePrivoderBehavior,
		},
	}

	for _, currentCase := range methodCases {
		s.Suite.Run(currentCase.Name, func() {
			ctx := context.Background()
			currentCase.MockDatabasePrivoderBehavior(s.MockDatabasePrivoder, ctx, currentCase.DatabaseResultError)
			err := s.CurrentRepository.Create(ctx, currentCase.InputPost)
			s.Assertions.Equal(currentCase.MethodResultError, err)
		})
	}
}

func (s *PostRepositorySuite) TestUpdateMethod() {
	type MockDatabasePrivoderBehavior func(*mock_database.MockDatabasePrivoder, context.Context, error)

	mockDatabasePrivoderBehavior := func(m *mock_database.MockDatabasePrivoder, input context.Context, returns error) {
		m.EXPECT().
			Exec(input, gomock.Any(), gomock.Any()).
			Return(returns).
			Times(1)
	}

	databaseResultError := errors.New("DatabaseResultError")

	methodCases := []struct {
		Name                         string
		InputPost                    domain.Post
		DatabaseResultError          error
		MethodResultError            error
		MockDatabasePrivoderBehavior MockDatabasePrivoderBehavior
	}{
		{
			Name:                         "Success",
			InputPost:                    domain.Post{},
			DatabaseResultError:          nil,
			MethodResultError:            nil,
			MockDatabasePrivoderBehavior: mockDatabasePrivoderBehavior,
		},
		{
			Name:                         "DatabaseFailure",
			InputPost:                    domain.Post{},
			DatabaseResultError:          databaseResultError,
			MethodResultError:            databaseResultError,
			MockDatabasePrivoderBehavior: mockDatabasePrivoderBehavior,
		},
	}

	for _, currentCase := range methodCases {
		s.Suite.Run(currentCase.Name, func() {
			ctx := context.Background()
			currentCase.MockDatabasePrivoderBehavior(s.MockDatabasePrivoder, ctx, currentCase.DatabaseResultError)
			err := s.CurrentRepository.Update(ctx, currentCase.InputPost)
			s.Assertions.Equal(currentCase.MethodResultError, err)
		})
	}
}

func (s *PostRepositorySuite) TestPublishMethod() {
	type MockDatabasePrivoderBehavior func(*mock_database.MockDatabasePrivoder, context.Context, error)

	mockDatabasePrivoderBehavior := func(m *mock_database.MockDatabasePrivoder, input context.Context, returns error) {
		m.EXPECT().
			Exec(input, gomock.Any(), gomock.Any()).
			Return(returns).
			Times(1)
	}

	databaseResultError := errors.New("DatabaseResultError")

	methodCases := []struct {
		Name                         string
		InputID                      uuid.UUID
		InputValue                   null.Time
		DatabaseResultError          error
		MethodResultError            error
		MockDatabasePrivoderBehavior MockDatabasePrivoderBehavior
	}{
		{
			Name:                         "Success",
			InputID:                      uuid.NewV4(),
			InputValue:                   null.NewTime(time.Now(), true),
			DatabaseResultError:          nil,
			MethodResultError:            nil,
			MockDatabasePrivoderBehavior: mockDatabasePrivoderBehavior,
		},
		{
			Name:                         "DatabaseFailure",
			InputID:                      uuid.NewV4(),
			InputValue:                   null.NewTime(time.Now(), true),
			DatabaseResultError:          databaseResultError,
			MethodResultError:            databaseResultError,
			MockDatabasePrivoderBehavior: mockDatabasePrivoderBehavior,
		},
	}

	for _, currentCase := range methodCases {
		s.Suite.Run(currentCase.Name, func() {
			ctx := context.Background()
			currentCase.MockDatabasePrivoderBehavior(s.MockDatabasePrivoder, ctx, currentCase.DatabaseResultError)
			err := s.CurrentRepository.Publish(ctx, currentCase.InputID, currentCase.InputValue)
			s.Assertions.Equal(currentCase.MethodResultError, err)
		})
	}
}
