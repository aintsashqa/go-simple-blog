package seeds

import (
	"context"

	"github.com/aintsashqa/go-simple-blog/pkg/database"
	"github.com/jaswdr/faker"
)

type SeedFunc func(context.Context, faker.Faker, database.DatabasePrivoder) error

type Seeder struct {
	database database.DatabasePrivoder
}

func NewSeeder(database database.DatabasePrivoder) *Seeder {
	return &Seeder{database: database}
}

func (s *Seeder) Seed(ctx context.Context, fslice ...SeedFunc) error {
	faker := faker.New()

	// TODO: start database transaction

	for _, f := range fslice {
		if err := f(ctx, faker, s.database); err != nil {
			// TODO: rollback database transaction
			return err
		}
	}

	// TODO: commit database transaction
	return nil
}
