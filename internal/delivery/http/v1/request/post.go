package request

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/aintsashqa/go-simple-blog/internal/delivery/http/v1/errors"
	"github.com/aintsashqa/go-simple-blog/internal/delivery/http/v1/response"
	"github.com/aintsashqa/go-simple-blog/internal/service"
	"github.com/go-chi/chi"
	uuid "github.com/satori/go.uuid"
)

const (
	DefaultCurrentPage  int = 1
	DefaultCountPerPage int = 15
)

type PostPaginationRequestDto struct {
	CurrentPage  int       `json:"-"`
	CountPerPage int       `json:"-"`
	UserID       uuid.UUID `json:"-"`
}

func (dto *PostPaginationRequestDto) FromRequest(r *http.Request) {
	currentPage, err := strconv.Atoi(r.URL.Query().Get("current_page"))
	if err != nil {
		currentPage = DefaultCurrentPage
	}

	countPerPage, err := strconv.Atoi(r.URL.Query().Get("count_per_page"))
	if err != nil {
		countPerPage = DefaultCountPerPage
	}

	userID := uuid.FromStringOrNil(r.URL.Query().Get("user_id"))

	dto.CurrentPage = currentPage
	dto.CountPerPage = countPerPage
	dto.UserID = userID
}

func (dto *PostPaginationRequestDto) TransformToObject() service.PaginatePostOptions {
	return service.PaginatePostOptions{
		CurrentPage:  dto.CurrentPage,
		PostsPerPage: dto.CountPerPage,
		UserID:       dto.UserID,
	}
}

type SelfPostPaginationRequestDto struct {
	CurrentPage  int       `json:"-"`
	CountPerPage int       `json:"-"`
	UserID       uuid.UUID `json:"-"`
}

func (dto *SelfPostPaginationRequestDto) FromRequest(r *http.Request) (response.ErrorResponseDto, error) {
	userID, casted := r.Context().Value("user_id").(uuid.UUID)
	if !casted {
		response := response.NewErrorResponseDto(http.StatusForbidden, errors.ErrInvalidTokenUserId.Error())
		return response, errors.ErrInvalidTokenUserId
	}

	currentPage, err := strconv.Atoi(r.URL.Query().Get("current_page"))
	if err != nil {
		currentPage = DefaultCurrentPage
	}

	countPerPage, err := strconv.Atoi(r.URL.Query().Get("count_per_page"))
	if err != nil {
		countPerPage = DefaultCountPerPage
	}

	dto.CurrentPage = currentPage
	dto.CountPerPage = countPerPage
	dto.UserID = userID

	return response.ErrorResponseDto{}, nil
}

func (dto *SelfPostPaginationRequestDto) TransformToObject() service.PaginatePostOptions {
	return service.PaginatePostOptions{
		CurrentPage:  dto.CurrentPage,
		PostsPerPage: dto.CountPerPage,
		UserID:       dto.UserID,
	}
}

type CreatePostRequestDto struct {
	Title       string    `json:"title"`
	Slug        string    `json:"slug"`
	Content     string    `json:"content"`
	UserID      uuid.UUID `json:"-"`
	IsPublished bool      `json:"is_published"`
}

func (dto *CreatePostRequestDto) FromRequest(r *http.Request) (response.ErrorResponseDto, error) {
	id, casted := r.Context().Value("user_id").(uuid.UUID)
	if !casted {
		response := response.NewErrorResponseDto(http.StatusForbidden, errors.ErrInvalidTokenUserId.Error())
		return response, errors.ErrInvalidTokenUserId
	}

	dto.UserID = id

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response := response.NewErrorResponseDto(http.StatusInternalServerError, errors.ErrUnavailableRequestBody.Error())
		return response, errors.ErrUnavailableRequestBody
	}

	return response.ErrorResponseDto{}, nil
}

func (dto *CreatePostRequestDto) TransformToObject() service.CreatePostInput {
	return service.CreatePostInput{
		Title:       dto.Title,
		Slug:        dto.Slug,
		Content:     dto.Content,
		UserID:      dto.UserID,
		IsPublished: dto.IsPublished,
	}
}

type UpdatePostRequestDto struct {
	ID          uuid.UUID `json:"-"`
	Title       string    `json:"title"`
	Slug        string    `json:"slug"`
	Content     string    `json:"content"`
	IsPublished bool      `json:"is_published"`
}

func (dto *UpdatePostRequestDto) FromRequest(r *http.Request) (response.ErrorResponseDto, error) {
	dto.ID = uuid.FromStringOrNil(chi.URLParam(r, "id"))

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		response := response.NewErrorResponseDto(http.StatusInternalServerError, errors.ErrUnavailableRequestBody.Error())
		return response, errors.ErrUnavailableRequestBody
	}

	return response.ErrorResponseDto{}, nil
}

func (dto *UpdatePostRequestDto) TransformToObject() service.UpdatePostInput {
	return service.UpdatePostInput{
		ID:          dto.ID,
		Title:       dto.Title,
		Slug:        dto.Slug,
		Content:     dto.Content,
		IsPublished: dto.IsPublished,
	}
}

type DeletePostRequestDto struct {
	UserID uuid.UUID `json:"-"`
	PostID uuid.UUID `json:"-"`
}

func (dto *DeletePostRequestDto) FromRequest(r *http.Request) (response.ErrorResponseDto, error) {
	userID, casted := r.Context().Value("user_id").(uuid.UUID)
	if !casted {
		response := response.NewErrorResponseDto(http.StatusForbidden, errors.ErrInvalidTokenUserId.Error())
		return response, errors.ErrInvalidTokenUserId
	}

	postID := uuid.FromStringOrNil(chi.URLParam(r, "id"))

	dto.UserID = userID
	dto.PostID = postID

	return response.ErrorResponseDto{}, nil
}

func (dto *DeletePostRequestDto) TransformToObject() service.SoftDeletePostInput {
	return service.SoftDeletePostInput{
		UserID: dto.UserID,
		PostID: dto.PostID,
	}
}
