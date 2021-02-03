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

	User interface {
		SignUp(context.Context, SignUpUserInput) error
		SignIn(context.Context, SignInUserInput) (Tokens, error)
		Authenticate(context.Context, AuthenticateUserInput) (uuid.UUID, error)
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

	PublishedPostsOptions struct {
		UserID uuid.UUID
		// IsPaginate   bool
		CurrentPage  int
		PostsPerPage int
	}

	PublishedPostsPagination struct {
		Posts        []domain.Post
		PreviousPage int
		CurrentPage  int
		NextPage     int
		PostsPerPage int
	}

	Post interface {
		Find(context.Context, uuid.UUID) (domain.Post, error)
		GetAllPublishedPaginate(context.Context, PublishedPostsOptions) (PublishedPostsPagination, error)
		Create(context.Context, CreatePostInput) (domain.Post, error)
		Update(context.Context, UpdatePostInput) (domain.Post, error)
		Publish(context.Context, uuid.UUID) error
	}

	Service struct {
		User
		Post
	}

	ServiceDependencies struct {
		Repository                    *repository.Repository
		Hasher                        hash.HashProvider
		Authorization                 auth.AuthorizationProvider
		AuthorizationTokenExpiresTime time.Duration
	}
)

func NewService(deps ServiceDependencies) *Service {
	return &Service{
		User: NewUserService(deps.Repository.User, deps.Hasher, deps.Authorization, deps.AuthorizationTokenExpiresTime),
		Post: NewPostService(deps.Repository.Post),
	}
}
