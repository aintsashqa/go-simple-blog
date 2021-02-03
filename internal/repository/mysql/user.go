package mysql

import (
	"context"
	"fmt"

	"github.com/aintsashqa/go-simple-blog/internal/domain"
	"github.com/aintsashqa/go-simple-blog/pkg/database"
)

type UserRepos struct {
	database database.DatabasePrivoder
}

func NewUserRepos(database database.DatabasePrivoder) *UserRepos {
	return &UserRepos{database: database}
}

func (r *UserRepos) Create(ctx context.Context, user domain.User) error {
	query := fmt.Sprintf("insert into %s (id, email, username, encrypted_password, created_at, updated_at, deleted_at) values (?, ?, ?, ?, ?, ?, ?)", usersTable)
	return r.database.Exec(ctx, query, user.ID, user.Email, user.Username, user.Password, user.CreatedAt, user.UpdatedAt, user.DeletedAt)
}

func (r *UserRepos) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	var user domain.User
	query := fmt.Sprintf("select * from %s where (email = ? and deleted_at is null)", usersTable)
	err := r.database.Get(ctx, &user, query, email)
	return user, err
}
