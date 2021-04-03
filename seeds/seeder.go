package seeds

import (
	"context"

	"github.com/aintsashqa/go-simple-blog/pkg/database"
	"github.com/jaswdr/faker"
)

type SeedFunc func(context.Context, faker.Faker, database.DatabaseInterface) error

type Seeder struct {
	database database.DatabasePrivoder
}

func NewSeeder(database database.DatabasePrivoder) *Seeder {
	return &Seeder{database: database}
}

func (s *Seeder) Seed(ctx context.Context, fslice ...SeedFunc) error {
	faker := faker.New()

	tx, err := s.database.BeginTx(ctx)
	if err != nil {
		return err
	}

	for _, f := range fslice {
		if err := f(ctx, faker, tx); err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}

			return err
		}
	}

	return tx.Commit()
}
