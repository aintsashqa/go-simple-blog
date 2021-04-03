//go:generate mockgen -source=repository.go -destination=mocks/mock.go
package repository

import (
	"context"

	"github.com/aintsashqa/go-simple-blog/internal/domain"
	"github.com/aintsashqa/go-simple-blog/internal/repository/mysql"
	"github.com/aintsashqa/go-simple-blog/pkg/database"
	uuid "github.com/satori/go.uuid"
)

type (
	User interface {
		Create(context.Context, domain.User) error
		GetByEmail(context.Context, string) (domain.User, error)
		Find(context.Context, uuid.UUID) (domain.User, error)
		Update(context.Context, domain.User) error
	}

	Post interface {
		Find(context.Context, uuid.UUID) (domain.Post, error)
		GetAllPublished(context.Context, int, int) ([]domain.Post, error)
		GetAllPublishedWithUserID(context.Context, uuid.UUID, int, int) ([]domain.Post, error)
		AllPublishedCount(context.Context) (int, error)
		AllPublishedCountWithUserID(context.Context, uuid.UUID) (int, error)
		Create(context.Context, domain.Post) error
		Update(context.Context, domain.Post) error
		Publish(context.Context, domain.Post) error
	}

	Repository struct {
		User
		Post
	}
)

func NewRepository(database database.DatabasePrivoder) *Repository {
	return &Repository{
		User: mysql.NewUserRepos(database),
		Post: mysql.NewPostRepos(database),
	}
}

func (r *Repository) UserProvider() User {
	return r.User
}

func (r *Repository) PostProvider() Post {
	return r.Post
}
