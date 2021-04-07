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
	PostCacheKey           = "post-cache-key-%s"
	PostCollectionCacheKey = "post-collection-cache-key-%d-%d-%s"
)

type PostCache struct {
	repo       repository.Post
	provider   cache.CachePrivoder
	serializer serializer.PostSerializer
}

func NewPostCache(repo repository.Post, provider cache.CachePrivoder, serializer serializer.PostSerializer) *PostCache {
	return &PostCache{repo: repo, provider: provider, serializer: serializer}
}

func (c *PostCache) Find(ctx context.Context, id uuid.UUID) (domain.Post, error) {
	key := fmt.Sprintf(PostCacheKey, id)

	if value, err := c.provider.Get(ctx, key); err == nil {
		return c.serializer.Deserialize(value)
	}

	post, err := c.repo.Find(ctx, id)
	if err != nil {
		return domain.Post{}, err
	}

	value, err := c.serializer.Serialize(post)
	if err != nil {
		return domain.Post{}, err
	}

	err = c.provider.Set(ctx, key, value)
	return post, err
}

func (c *PostCache) FindWithPrimaryAndUserID(ctx context.Context, postID uuid.UUID, userID uuid.UUID) (domain.Post, error) {
	key := fmt.Sprintf(PostCacheKey, postID)

	if value, err := c.provider.Get(ctx, key); err == nil {
		return c.serializer.Deserialize(value)
	}

	post, err := c.repo.FindWithPrimaryAndUserID(ctx, postID, userID)
	if err != nil {
		return domain.Post{}, err
	}

	value, err := c.serializer.Serialize(post)
	if err != nil {
		return domain.Post{}, err
	}

	err = c.provider.Set(ctx, key, value)
	return post, err
}

func (c *PostCache) GetAllPublished(ctx context.Context, offset int, count int) ([]domain.Post, error) {
	return c.repo.GetAllPublished(ctx, offset, count)
}

func (c *PostCache) GetAllPublishedWithUserID(ctx context.Context, id uuid.UUID, offset int, count int) ([]domain.Post, error) {
	return c.repo.GetAllPublishedWithUserID(ctx, id, offset, count)
}

func (c *PostCache) GetAllWithUserID(ctx context.Context, id uuid.UUID, offset int, count int) ([]domain.Post, error) {
	return c.repo.GetAllWithUserID(ctx, id, offset, count)
}

func (c *PostCache) AllPublishedCount(ctx context.Context) (int, error) {
	return c.repo.AllPublishedCount(ctx)
}

func (c *PostCache) AllPublishedCountWithUserID(ctx context.Context, id uuid.UUID) (int, error) {
	return c.repo.AllPublishedCountWithUserID(ctx, id)
}

func (c *PostCache) TotalCountWithUserID(ctx context.Context, id uuid.UUID) (int, error) {
	return c.repo.TotalCountWithUserID(ctx, id)
}

func (c *PostCache) Create(ctx context.Context, post domain.Post) error {
	err := c.repo.Create(ctx, post)
	if err != nil {
		return err
	}

	value, err := c.serializer.Serialize(post)
	if err != nil {
		return err
	}

	key := fmt.Sprintf(PostCacheKey, post.ID)
	return c.provider.Set(ctx, key, value)
}

func (c *PostCache) Update(ctx context.Context, post domain.Post) error {
	err := c.repo.Update(ctx, post)
	if err != nil {
		return err
	}

	value, err := c.serializer.Serialize(post)
	if err != nil {
		return err
	}

	key := fmt.Sprintf(PostCacheKey, post.ID)
	return c.provider.Set(ctx, key, value)
}

func (c *PostCache) Publish(ctx context.Context, post domain.Post) error {
	err := c.repo.Publish(ctx, post)
	if err != nil {
		return err
	}

	value, err := c.serializer.Serialize(post)
	if err != nil {
		return err
	}

	key := fmt.Sprintf(PostCacheKey, post.ID)
	return c.provider.Set(ctx, key, value)
}

func (c *PostCache) SoftDelete(ctx context.Context, post domain.Post) error {
	key := fmt.Sprintf(PostCacheKey, post.ID)

	if err := c.provider.Delete(ctx, key); err != nil {
		return err
	}

	return c.repo.SoftDelete(ctx, post)
}
