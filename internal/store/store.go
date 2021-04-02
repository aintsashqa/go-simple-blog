package store

import (
	"github.com/aintsashqa/go-simple-blog/internal/repository"
	"github.com/aintsashqa/go-simple-blog/internal/serializer"
	"github.com/aintsashqa/go-simple-blog/internal/store/redis"
	"github.com/aintsashqa/go-simple-blog/pkg/cache"
)

type CacheStore struct {
	User repository.User
	Post repository.Post
}

func NewCacheStore(repos *repository.Repository, cache cache.CachePrivoder, serializer *serializer.Serializer) *CacheStore {
	return &CacheStore{
		User: redis.NewUserCache(repos.User, cache, serializer.User),
		Post: redis.NewPostCache(repos.Post, cache, serializer.Post),
	}
}

func (s *CacheStore) UserProvider() repository.User {
	return s.User
}

func (s *CacheStore) PostProvider() repository.Post {
	return s.Post
}
