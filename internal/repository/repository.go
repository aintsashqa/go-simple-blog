//go:generate mockgen -source=repository.go -destination=mocks/mock.go
package repository

import (
	"context"

	"github.com/aintsashqa/go-simple-blog/internal/domain"
	"github.com/aintsashqa/go-simple-blog/internal/repository/mysql"
	"github.com/aintsashqa/go-simple-blog/pkg/database"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/guregu/null.v4"
)

type (
	User interface {
		Create(context.Context, domain.User) error
		GetByEmail(context.Context, string) (domain.User, error)
	}

	Post interface {
		Find(context.Context, uuid.UUID) (domain.Post, error)
		GetAllPublished(context.Context, int, int) ([]domain.Post, error)
		GetAllPublishedWithUserID(context.Context, uuid.UUID, int, int) ([]domain.Post, error)
		AllPublishedCount(context.Context) (int, error)
		Create(context.Context, domain.Post) error
		Update(context.Context, domain.Post) error
		Publish(context.Context, uuid.UUID, null.Time) error
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
