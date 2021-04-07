package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/aintsashqa/go-simple-blog/internal/domain"
	"github.com/aintsashqa/go-simple-blog/internal/repository/errors"
	"github.com/aintsashqa/go-simple-blog/pkg/database"
	uuid "github.com/satori/go.uuid"
)

type PostRepos struct {
	database database.DatabasePrivoder
}

func NewPostRepos(database database.DatabasePrivoder) *PostRepos {
	return &PostRepos{database: database}
}

func (r *PostRepos) Find(ctx context.Context, id uuid.UUID) (domain.Post, error) {
	var post domain.Post
	query := fmt.Sprintf("select * from %s where (id = ? and deleted_at is null)", postsTable)
	err := r.database.Get(ctx, &post, query, id)
	if err == sql.ErrNoRows {
		return post, errors.ErrPostNotFound
	}
	return post, err
}

func (r *PostRepos) FindWithPrimaryAndUserID(ctx context.Context, postID uuid.UUID, userID uuid.UUID) (domain.Post, error) {
	var post domain.Post
	query := fmt.Sprintf("select * from %s where (id = ? and user_id = ? and deleted_at is null)", postsTable)
	err := r.database.Get(ctx, &post, query, postID, userID)
	if err == sql.ErrNoRows {
		return post, errors.ErrPostNotFound
	}
	return post, err
}

func (r *PostRepos) GetAllPublished(ctx context.Context, offset, count int) ([]domain.Post, error) {
	var posts []domain.Post
	query := fmt.Sprintf("select * from %s where (published_at is not null and deleted_at is null) limit ?, ?", postsTable)
	err := r.database.Select(ctx, &posts, query, offset, count)
	if posts == nil {
		posts = []domain.Post{}
	}
	return posts, err
}

func (r *PostRepos) GetAllPublishedWithUserID(ctx context.Context, id uuid.UUID, offset, count int) ([]domain.Post, error) {
	var posts []domain.Post
	query := fmt.Sprintf("select * from %s where (user_id = ? and published_at is not null and deleted_at is null) limit ?, ?", postsTable)
	err := r.database.Select(ctx, &posts, query, id, offset, count)
	if posts == nil {
		posts = []domain.Post{}
	}
	return posts, err
}

func (r *PostRepos) AllPublishedCount(ctx context.Context) (int, error) {
	var count int
	query := fmt.Sprintf("select count(*) from %s where (published_at is not null and deleted_at is null)", postsTable)
	err := r.database.QueryRow(ctx, &count, query)
	return count, err
}

func (r *PostRepos) AllPublishedCountWithUserID(ctx context.Context, id uuid.UUID) (int, error) {
	var count int
	query := fmt.Sprintf("select count(*) from %s where (user_id = ? and published_at is not null and deleted_at is null)", postsTable)
	err := r.database.QueryRow(ctx, &count, query, id)
	return count, err
}

func (r *PostRepos) Create(ctx context.Context, post domain.Post) error {
	query := fmt.Sprintf("insert into %s (id, title, slug, content, user_id, created_at, updated_at, published_at, deleted_at) values (?, ?, ?, ?, ?, ?, ?, ?, ?)", postsTable)
	return r.database.Exec(ctx, query, post.ID, post.Title, post.Slug, post.Content, post.UserID, post.CreatedAt, post.UpdatedAt, post.PublishedAt, post.DeletedAt)
}

func (r *PostRepos) Update(ctx context.Context, post domain.Post) error {
	query := fmt.Sprintf("update %s set title = ?, slug = ?, content = ?, updated_at = ?, published_at = ? where (id = ? and deleted_at is null)", postsTable)
	err := r.database.Exec(ctx, query, post.Title, post.Slug, post.Content, post.UpdatedAt, post.PublishedAt, post.ID)
	if err == sql.ErrNoRows {
		return errors.ErrPostNotFound
	}
	return err
}

func (r *PostRepos) Publish(ctx context.Context, post domain.Post) error {
	query := fmt.Sprintf("update %s set published_at = ?, updated_at = ? where (id = ? and deleted_at is null)", postsTable)
	err := r.database.Exec(ctx, query, post.PublishedAt, post.UpdatedAt, post.ID)
	if err == sql.ErrNoRows {
		return errors.ErrPostNotFound
	}
	return err
}

func (r *PostRepos) SoftDelete(ctx context.Context, post domain.Post) error {
	query := fmt.Sprintf("update %s set updated_at = ?, deleted_at = ? where (id = ? and deleted_at is null)", postsTable)
	err := r.database.Exec(ctx, query, post.UpdatedAt, post.DeletedAt, post.ID)
	if err == sql.ErrNoRows {
		return errors.ErrPostNotFound
	}
	return err
}
