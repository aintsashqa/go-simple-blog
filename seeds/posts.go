package seeds

import (
	"context"
	"math/rand"
	"time"

	"github.com/aintsashqa/go-simple-blog/internal/domain"
	"github.com/aintsashqa/go-simple-blog/pkg/database"
	"github.com/gosimple/slug"
	"github.com/jaswdr/faker"
	"gopkg.in/guregu/null.v4"
)

func PostSeed(ctx context.Context, faker faker.Faker, database database.DatabasePrivoder) error {
	rand.Seed(time.Now().Unix())

	userQuery := "select * from users"
	trancate := "truncate table posts"
	query := "insert into posts (id, title, slug, content, user_id, created_at, updated_at, published_at, deleted_at) values (?, ?, ?, ?, ?, ?, ?, ?, ?)"

	var users []domain.User
	if err := database.Select(ctx, &users, userQuery); err != nil {
		return err
	}

	if err := database.Exec(ctx, trancate); err != nil {
		return err
	}

	for _, user := range users {
		for i := 0; i < 15; i++ {
			title := faker.Lorem().Sentence(3)
			slug := slug.Make(title)
			isPublished := rand.Intn(3)%2 == 0

			temp := domain.Post{
				Title:       title,
				Slug:        slug,
				Content:     faker.Lorem().Text(1000),
				UserID:      user.ID,
				PublishedAt: null.NewTime(time.Now(), isPublished),
			}
			temp.Init()

			if err := database.Exec(ctx, query, temp.ID, temp.Title, temp.Slug, temp.Content, temp.UserID, temp.CreatedAt, temp.UpdatedAt, temp.PublishedAt, temp.DeletedAt); err != nil {
				return err
			}
		}
	}

	return nil
}
