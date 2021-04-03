package service

import (
	"context"
	"time"

	"github.com/aintsashqa/go-simple-blog/internal/domain"
	"github.com/aintsashqa/go-simple-blog/internal/repository"
	"github.com/gosimple/slug"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/guregu/null.v4"
)

type PostService struct {
	repo repository.Post
}

func NewPostService(repo repository.Post) *PostService {
	return &PostService{repo: repo}
}

func (s *PostService) Find(ctx context.Context, id uuid.UUID) (domain.Post, error) {
	return s.repo.Find(ctx, id)
}

func (s *PostService) GetAllPublishedPaginate(ctx context.Context, opt PublishedPostsOptions) (PublishedPostsPagination, error) {
	var count int
	var posts []domain.Post
	var err error

	offset := (opt.CurrentPage - 1) * opt.PostsPerPage

	if opt.UserID != uuid.Nil {

		posts, err = s.repo.GetAllPublishedWithUserID(ctx, opt.UserID, offset, opt.PostsPerPage)
		if err != nil {
			return PublishedPostsPagination{}, err
		}

		count, err = s.repo.AllPublishedCountWithUserID(ctx, opt.UserID)
		if err != nil {
			return PublishedPostsPagination{}, err
		}

	} else {

		posts, err = s.repo.GetAllPublished(ctx, offset, opt.PostsPerPage)
		if err != nil {
			return PublishedPostsPagination{}, err
		}

		count, err = s.repo.AllPublishedCount(ctx)
		if err != nil {
			return PublishedPostsPagination{}, err
		}
	}

	previousPage := opt.CurrentPage - 1
	if previousPage < 1 {
		previousPage = 1
	}

	nextPage := opt.CurrentPage + 1
	value := count - (opt.PostsPerPage * opt.CurrentPage)
	if value <= 0 {
		nextPage = opt.CurrentPage
	}

	return PublishedPostsPagination{
		Posts:        posts,
		PostsCount:   count,
		PreviousPage: previousPage,
		CurrentPage:  opt.CurrentPage,
		NextPage:     nextPage,
		PostsPerPage: opt.PostsPerPage,
	}, nil
}

func (s *PostService) Create(ctx context.Context, input CreatePostInput) (domain.Post, error) {
	slugStr := input.Slug
	if len(slugStr) == 0 {
		slugStr = slug.Make(input.Title)
	}

	post := domain.Post{
		Title:       input.Title,
		Slug:        slugStr,
		Content:     input.Content,
		UserID:      input.UserID,
		PublishedAt: null.NewTime(time.Now(), input.IsPublished),
	}
	post.Init()

	err := s.repo.Create(ctx, post)
	return post, err
}

func (s *PostService) Update(ctx context.Context, input UpdatePostInput) (domain.Post, error) {
	post, err := s.repo.Find(ctx, input.ID)
	if err != nil {
		return domain.Post{}, err
	}

	slugStr := input.Slug
	if len(slugStr) == 0 {
		slugStr = slug.Make(input.Title)
	}

	post.Title = input.Title
	post.Slug = slugStr
	post.Content = input.Content
	post.PublishedAt = null.NewTime(time.Now(), input.IsPublished)
	post.Update()

	if err := s.repo.Update(ctx, post); err != nil {
		return domain.Post{}, err
	}

	return post, nil
}

func (s *PostService) Publish(ctx context.Context, id uuid.UUID) (domain.Post, error) {
	post, err := s.repo.Find(ctx, id)
	if err != nil {
		return domain.Post{}, err
	}

	post.PublishedAt = null.NewTime(time.Now(), true)
	post.Update()

	err = s.repo.Publish(ctx, post)
	return post, err
}
