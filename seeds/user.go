package seeds

import (
	"context"

	"github.com/aintsashqa/go-simple-blog/internal/domain"
	"github.com/aintsashqa/go-simple-blog/pkg/database"
	"github.com/aintsashqa/go-simple-blog/pkg/hash/bcrypt"
	"github.com/jaswdr/faker"
)

func UserSeed(ctx context.Context, faker faker.Faker, tx database.DatabaseInterface) error {
	hasher := bcrypt.NewBcryptProvider()
	trancate := "truncate table users"
	query := "insert into users (id, email, username, encrypted_password, created_at, updated_at, deleted_at) values (?, ?, ?, ?, ?, ?, ?)"

	if err := tx.Exec(ctx, trancate); err != nil {
		return err
	}

	for i := 0; i < 10; i++ {
		temp := domain.User{
			Email:    faker.Internet().Email(),
			Username: faker.Internet().User(),
			Password: hasher.Make("secret"),
		}
		temp.Init()

		if err := tx.Exec(ctx, query, temp.ID, temp.Email, temp.Username, temp.Password, temp.CreatedAt, temp.UpdatedAt, temp.DeletedAt); err != nil {
			return err
		}
	}

	return nil
}
