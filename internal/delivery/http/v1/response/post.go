package response

import (
	"time"

	"github.com/aintsashqa/go-simple-blog/internal/domain"
	"github.com/aintsashqa/go-simple-blog/internal/service"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/guregu/null.v4"
)

type PaginationResponseDto struct {
	Count        int `json:"count"`
	PreviousPage int `json:"previous_page"`
	CurrentPage  int `json:"current_page"`
	NextPage     int `json:"next_page"`
	CountPerPage int `json:"count_per_page"`
}

type PostPaginationResponseDto struct {
	Posts      []PostResponseDto     `json:"posts"`
	Pagination PaginationResponseDto `json:"pagination"`
}

func (dto *PostPaginationResponseDto) TransformFromObject(posts service.PublishedPostsPagination) {
	dto.Posts = []PostResponseDto{}

	for _, post := range posts.Posts {
		temp := PostResponseDto{}
		temp.TransformFromObject(post)
		dto.Posts = append(dto.Posts, temp)
	}

	dto.Pagination = PaginationResponseDto{
		Count:        posts.PostsCount,
		PreviousPage: posts.PreviousPage,
		CurrentPage:  posts.CurrentPage,
		NextPage:     posts.NextPage,
		CountPerPage: posts.PostsPerPage,
	}
}

type PostResponseDto struct {
	ID      uuid.UUID `json:"id"`
	Title   string    `json:"title"`
	Slug    string    `json:"slug"`
	Content string    `json:"content"`
	UserID  uuid.UUID `json:"user_id"`
	// User        *UserResponseDto `json:"user,omitempty"`
	IsPublished bool      `json:"is_published"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	PublishedAt null.Time `json:"published_at"`
}

func (dto *PostResponseDto) TransformFromObject(post domain.Post) {
	dto.ID = post.ID
	dto.Title = post.Title
	dto.Slug = post.Slug
	dto.Content = post.Content
	dto.UserID = post.UserID
	dto.CreatedAt = post.CreatedAt
	dto.UpdatedAt = post.UpdatedAt

	if post.PublishedAt.Valid {
		dto.IsPublished = post.PublishedAt.Valid
		dto.PublishedAt = post.PublishedAt
	}
}
