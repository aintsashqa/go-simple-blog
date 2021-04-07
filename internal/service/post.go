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

func (s *PostService) offset(page, perPage int) int {
	return (page - 1) * perPage
}

func (s *PostService) pagination(page, perPage, total int) (int, int) {
	previousPage := page - 1
	if previousPage < 1 {
		previousPage = 1
	}

	nextPage := page + 1
	value := total - (perPage * page)
	if value <= 0 {
		nextPage = page
	}

	return previousPage, nextPage
}

func (s *PostService) GetAllPublishedPaginate(ctx context.Context, opt PaginatePostOptions) (PostPagination, error) {
	var count int
	var posts []domain.Post
	var err error

	offset := s.offset(opt.CurrentPage, opt.PostsPerPage)

	if opt.UserID != uuid.Nil {

		posts, err = s.repo.GetAllPublishedWithUserID(ctx, opt.UserID, offset, opt.PostsPerPage)
		if err != nil {
			return PostPagination{}, err
		}

		count, err = s.repo.AllPublishedCountWithUserID(ctx, opt.UserID)
		if err != nil {
			return PostPagination{}, err
		}

	} else {

		posts, err = s.repo.GetAllPublished(ctx, offset, opt.PostsPerPage)
		if err != nil {
			return PostPagination{}, err
		}

		count, err = s.repo.AllPublishedCount(ctx)
		if err != nil {
			return PostPagination{}, err
		}
	}

	previousPage, nextPage := s.pagination(opt.CurrentPage, opt.PostsPerPage, count)

	return PostPagination{
		Posts:        posts,
		PostsCount:   count,
		PreviousPage: previousPage,
		CurrentPage:  opt.CurrentPage,
		NextPage:     nextPage,
		PostsPerPage: opt.PostsPerPage,
	}, nil
}

func (s *PostService) GetAllSelfPaginate(ctx context.Context, opt PaginatePostOptions) (PostPagination, error) {
	offset := s.offset(opt.CurrentPage, opt.PostsPerPage)

	posts, err := s.repo.GetAllWithUserID(ctx, opt.UserID, offset, opt.PostsPerPage)
	if err != nil {
		return PostPagination{}, err
	}

	count, err := s.repo.TotalCountWithUserID(ctx, opt.UserID)
	if err != nil {
		return PostPagination{}, err
	}

	previousPage, nextPage := s.pagination(opt.CurrentPage, opt.PostsPerPage, count)

	return PostPagination{
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

	if err := post.Validate(domain.CreatePostValidationAction); err != nil {
		return domain.Post{}, err
	}

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

	if err := post.Validate(domain.UpdatePostValidationAction); err != nil {
		return domain.Post{}, err
	}

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

func (s *PostService) SoftDelete(ctx context.Context, input SoftDeletePostInput) error {
	post, err := s.repo.FindWithPrimaryAndUserID(ctx, input.PostID, input.UserID)
	if err != nil {
		return err
	}

	post.Delete()
	return s.repo.SoftDelete(ctx, post)
}
