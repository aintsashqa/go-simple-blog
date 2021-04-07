package redis

import (
	"context"
	"fmt"

	"github.com/aintsashqa/go-simple-blog/internal/domain"
	"github.com/aintsashqa/go-simple-blog/internal/repository"
	"github.com/aintsashqa/go-simple-blog/internal/serializer"
	"github.com/aintsashqa/go-simple-blog/pkg/cache"
	uuid "github.com/satori/go.uuid"
)

const (
	UserCacheKey string = "user-cache-key-%s"
)

type UserCache struct {
	repo       repository.User
	provider   cache.CachePrivoder
	serializer serializer.UserSerializer
}

func NewUserCache(repo repository.User, provider cache.CachePrivoder, serializer serializer.UserSerializer) *UserCache {
	return &UserCache{repo: repo, provider: provider, serializer: serializer}
}

func (c *UserCache) Create(ctx context.Context, user domain.User) error {
	err := c.repo.Create(ctx, user)
	if err != nil {
		return err
	}

	value, err := c.serializer.Serialize(user)
	if err != nil {
		return err
	}

	key := fmt.Sprintf(UserCacheKey, user.ID)
	return c.provider.Set(ctx, key, value)
}

func (c *UserCache) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	return c.repo.GetByEmail(ctx, email)
}

func (c *UserCache) Find(ctx context.Context, id uuid.UUID) (domain.User, error) {
	key := fmt.Sprintf(UserCacheKey, id)

	if value, err := c.provider.Get(ctx, key); err == nil {
		return c.serializer.Deserialize(value)
	}

	user, err := c.repo.Find(ctx, id)
	if err != nil {
		return domain.User{}, err
	}

	value, err := c.serializer.Serialize(user)
	if err != nil {
		return domain.User{}, err
	}

	err = c.provider.Set(ctx, key, value)
	return user, err
}

func (c *UserCache) Self(ctx context.Context, id uuid.UUID) (domain.User, error) {
	key := fmt.Sprintf(UserCacheKey, id)

	if value, err := c.provider.Get(ctx, key); err == nil {
		return c.serializer.Deserialize(value)
	}

	user, err := c.repo.Self(ctx, id)
	if err != nil {
		return domain.User{}, err
	}

	value, err := c.serializer.Serialize(user)
	if err != nil {
		return domain.User{}, err
	}

	err = c.provider.Set(ctx, key, value)
	return user, err
}

func (c *UserCache) Update(ctx context.Context, user domain.User) error {
	err := c.repo.Update(ctx, user)
	if err != nil {
		return err
	}

	value, err := c.serializer.Serialize(user)
	if err != nil {
		return err
	}

	key := fmt.Sprintf(UserCacheKey, user.ID)
	return c.provider.Set(ctx, key, value)
}
