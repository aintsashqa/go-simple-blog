//go:generate mockgen -source=service.go -destination=mocks/mock.go
package service

import (
	"context"
	"time"

	"github.com/aintsashqa/go-simple-blog/internal/domain"
	"github.com/aintsashqa/go-simple-blog/internal/repository"
	"github.com/aintsashqa/go-simple-blog/pkg/auth"
	"github.com/aintsashqa/go-simple-blog/pkg/hash"
	uuid "github.com/satori/go.uuid"
)

type (
	SignUpUserInput struct {
		Email    string
		Password string
	}

	SignInUserInput struct {
		Email    string
		Password string
	}

	Tokens struct {
		AccessToken string
	}

	AuthenticateUserInput struct {
		Token string
	}

	UpdateUserInput struct {
		ID       uuid.UUID
		Username string
	}

	User interface {
		SignUp(context.Context, SignUpUserInput) (domain.User, error)
		SignIn(context.Context, SignInUserInput) (Tokens, error)
		Authenticate(context.Context, AuthenticateUserInput) (uuid.UUID, error)
		Find(context.Context, uuid.UUID) (domain.User, error)
		Self(context.Context, uuid.UUID) (domain.User, error)
		Update(context.Context, UpdateUserInput) (domain.User, error)
	}

	CreatePostInput struct {
		Title       string
		Slug        string
		Content     string
		UserID      uuid.UUID
		IsPublished bool
	}

	UpdatePostInput struct {
		ID          uuid.UUID
		Title       string
		Slug        string
		Content     string
		IsPublished bool
	}

	PaginatePostOptions struct {
		UserID       uuid.UUID
		CurrentPage  int
		PostsPerPage int
	}

	PostPagination struct {
		Posts        []domain.Post
		PostsCount   int
		PreviousPage int
		CurrentPage  int
		NextPage     int
		PostsPerPage int
	}

	SoftDeletePostInput struct {
		UserID uuid.UUID
		PostID uuid.UUID
	}

	Post interface {
		Find(context.Context, uuid.UUID) (domain.Post, error)
		GetAllPublishedPaginate(context.Context, PaginatePostOptions) (PostPagination, error)
		GetAllSelfPaginate(context.Context, PaginatePostOptions) (PostPagination, error)
		Create(context.Context, CreatePostInput) (domain.Post, error)
		Update(context.Context, UpdatePostInput) (domain.Post, error)
		Publish(context.Context, uuid.UUID) (domain.Post, error)
		SoftDelete(context.Context, SoftDeletePostInput) error
	}

	Service struct {
		User
		Post
	}

	DataProvider interface {
		UserProvider() repository.User
		PostProvider() repository.Post
	}

	ServiceDependencies struct {
		DataProvider                  DataProvider
		Hasher                        hash.HashProvider
		Authorization                 auth.AuthorizationProvider
		AuthorizationTokenExpiresTime time.Duration
	}
)

func NewService(deps ServiceDependencies) *Service {
	return &Service{
		User: NewUserService(deps.DataProvider.UserProvider(), deps.Hasher, deps.Authorization, deps.AuthorizationTokenExpiresTime),
		Post: NewPostService(deps.DataProvider.PostProvider()),
	}
}
