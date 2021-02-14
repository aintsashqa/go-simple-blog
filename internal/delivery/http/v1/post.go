package v1

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/aintsashqa/go-simple-blog/internal/domain"
	"github.com/aintsashqa/go-simple-blog/internal/service"
	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	uuid "github.com/satori/go.uuid"
)

type getAllPublishedPostsResponseMetadata struct {
	PreviousPage int `json:"previous_page"`
	CurrentPage  int `json:"current_page"`
	NextPage     int `json:"next_page"`
	PostsPerPage int `json:"posts_per_page"`
}

type getAllPublishedPostsResponse struct {
	Posts    []domain.Post                        `json:"posts"`
	Metadata getAllPublishedPostsResponseMetadata `json:"metadata"`
}

// @Summary Get all published
// @Description Get all published with pagination
// @ID post-get-all-published
// @Tags Post
// @Accept json
// @Produce json
// @Param current_page query int true "Number of current page"
// @Param posts_per_page query int true "Number of posts count"
// @Param user_id query string false "Posts with user id"
// @Success 200 {object} getAllPublishedPostsResponse
// @Failure 400 {object} responseError
// @Failure 500 {object} responseError
// @Router /post [get]
func (h *Handler) getAllPublishedPosts(w http.ResponseWriter, r *http.Request) {
	currentPage, err := strconv.Atoi(r.URL.Query().Get("current_page"))
	if err != nil {
		log.Print(err)
		errorFn(w, r, http.StatusBadRequest, fmt.Sprintf(invalidUrlQueryParamErrorMsg, "current_page"), err.Error())
		return
	}

	postsPerPage, err := strconv.Atoi(r.URL.Query().Get("posts_per_page"))
	if err != nil {
		log.Print(err)
		errorFn(w, r, http.StatusBadRequest, fmt.Sprintf(invalidUrlQueryParamErrorMsg, "posts_per_page"), err.Error())
		return
	}

	opt := service.PublishedPostsOptions{
		UserID:       uuid.FromStringOrNil(r.URL.Query().Get("user_id")),
		CurrentPage:  currentPage,
		PostsPerPage: postsPerPage,
	}

	pagination, err := h.post.GetAllPublishedPaginate(r.Context(), opt)
	if err != nil {
		log.Print(err)
		errorFn(w, r, http.StatusInternalServerError, internalErrorMsg, err.Error())
		return
	}

	respond(w, r, http.StatusOK, getAllPublishedPostsResponse{
		Posts: pagination.Posts,
		Metadata: getAllPublishedPostsResponseMetadata{
			PreviousPage: pagination.PreviousPage,
			CurrentPage:  pagination.CurrentPage,
			NextPage:     pagination.NextPage,
			PostsPerPage: pagination.PostsPerPage,
		},
	})
}

type getSinglePostResponse struct {
	Post domain.Post `json:"post"`
}

// @Summary Get single post
// @Description Get single post by id
// @ID post-get-single
// @Tags Post
// @Accept json
// @Produce json
// @Param id path string string "Post with id"
// @Success 200 {object} getSinglePostResponse
// @Failure 500 {object} responseError
// @Router /post/{id} [get]
func (h *Handler) getSinglePost(w http.ResponseWriter, r *http.Request) {
	id := uuid.FromStringOrNil(chi.URLParam(r, "id"))
	post, err := h.post.Find(r.Context(), id)
	if err != nil {
		log.Print(err)
		errorFn(w, r, http.StatusInternalServerError, internalErrorMsg, err.Error())
		return
	}

	respond(w, r, http.StatusOK, getSinglePostResponse{
		Post: post,
	})
}

type createPostRequest struct {
	Title       string `json:"title"`
	Slug        string `json:"slug"`
	Content     string `json:"content"`
	IsPublished bool   `json:"is_published"`
}

func (r *createPostRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Title, validation.Required, validation.Length(8, 255)),
		validation.Field(&r.Slug, validation.Required, validation.Length(8, 255)),
		validation.Field(&r.Content, validation.Required, validation.Length(500, 0)),
		validation.Field(&r.IsPublished),
	)
}

type createPostResponse struct {
	Post domain.Post `json:"post"`
}

// @Summary Create post
// @Description Create new post
// @ID post-create
// @Tags Post
// @Accept json
// @Produce json
// @Param payload body createPostRequest true "Post details"
// @Success 201 {object} createPostResponse
// @Failure 400 {object} responseError
// @Failure 401 {object} responseError
// @Failure 500 {object} responseError
// @Security ApiKeyAuth
// @Router /post [post]
func (h *Handler) createPost(w http.ResponseWriter, r *http.Request) {
	var input createPostRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Print(err)
		errorFn(w, r, http.StatusInternalServerError, internalErrorMsg, err.Error())
		return
	}

	if err := input.Validate(); err != nil {
		log.Print(err)
		errorFn(w, r, http.StatusBadRequest, validationFailedMsg, err.Error())
		return
	}

	id := r.Context().Value("user_id").(uuid.UUID)
	opt := service.CreatePostInput{
		Title:       input.Title,
		Slug:        input.Slug,
		Content:     input.Content,
		UserID:      id,
		IsPublished: input.IsPublished,
	}

	post, err := h.post.Create(r.Context(), opt)
	if err != nil {
		log.Print(err)
		errorFn(w, r, http.StatusInternalServerError, internalErrorMsg, err.Error())
		return
	}

	respond(w, r, http.StatusCreated, createPostResponse{
		Post: post,
	})
}

type updatePostRequest struct {
	Title       string `json:"title"`
	Slug        string `json:"slug"`
	Content     string `json:"content"`
	IsPublished bool   `json:"is_published"`
}

func (r *updatePostRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Title, validation.NotNil, validation.Length(8, 255)),
		validation.Field(&r.Slug, validation.NotNil, validation.Length(8, 255)),
		validation.Field(&r.Content, validation.NotNil, validation.Length(500, 0)),
		validation.Field(&r.IsPublished),
	)
}

type updatePostResponse struct {
	Post domain.Post `json:"post"`
}

// @Summary Update post
// @Description Update post with id
// @ID post-update
// @Tags Post
// @Accept json
// @Produce json
// @Param id path string true "Post with id"
// @Param payload body updatePostRequest true "Post details"
// @Success 202 {object} updatePostResponse
// @Failure 400 {object} responseError
// @Failure 401 {object} responseError
// @Failure 500 {object} responseError
// @Security ApiKeyAuth
// @Router /post/{id} [put]
func (h *Handler) updatePost(w http.ResponseWriter, r *http.Request) {
	var input updatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Print(err)
		errorFn(w, r, http.StatusInternalServerError, internalErrorMsg, err.Error())
		return
	}

	if err := input.Validate(); err != nil {
		log.Print(err)
		errorFn(w, r, http.StatusBadRequest, validationFailedMsg, err.Error())
		return
	}

	id := uuid.FromStringOrNil(chi.URLParam(r, "id"))
	opt := service.UpdatePostInput{
		ID:          id,
		Title:       input.Title,
		Slug:        input.Slug,
		Content:     input.Content,
		IsPublished: input.IsPublished,
	}

	post, err := h.post.Update(r.Context(), opt)
	if err != nil {
		log.Print(err)
		errorFn(w, r, http.StatusInternalServerError, internalErrorMsg, err.Error())
		return
	}

	respond(w, r, http.StatusAccepted, updatePostResponse{
		Post: post,
	})
}

type publishPostResponse struct {
	Message string `json:"message"`
}

// @Summary Publish post
// @Description Publish post with id
// @ID post-publish
// @Tags Post
// @Accept json
// @Produce json
// @Param id path string true "Post with id"
// @Success 202 {object} publishPostResponse
// @Failure 401 {object} responseError
// @Failure 500 {object} responseError
// @Security ApiKeyAuth
// @Router /post/{id}/publish [get]
func (h *Handler) publishPost(w http.ResponseWriter, r *http.Request) {
	id := uuid.FromStringOrNil(chi.URLParam(r, "id"))
	if err := h.post.Publish(r.Context(), id); err != nil {
		log.Print(err)
		errorFn(w, r, http.StatusInternalServerError, internalErrorMsg, err.Error())
		return
	}

	respond(w, r, http.StatusAccepted, publishPostResponse{
		Message: "Accepted",
	})
}
