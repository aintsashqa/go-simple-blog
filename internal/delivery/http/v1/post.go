package v1

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aintsashqa/go-simple-blog/internal/delivery/http/v1/errors"
	requsetdto "github.com/aintsashqa/go-simple-blog/internal/delivery/http/v1/request"
	responsedto "github.com/aintsashqa/go-simple-blog/internal/delivery/http/v1/response"
	repoerrors "github.com/aintsashqa/go-simple-blog/internal/repository/errors"
	"github.com/go-chi/chi"
	uuid "github.com/satori/go.uuid"
)

const (
	GetAllPublishedPostsCacheKey string = "get-all-published-posts[current-page=%d;posts-per-page=%d;user-id=%s]"
	GetSinglePostCacheKey        string = "get-single-post[id=%s]"
)

// @Summary Get all published
// @Description Get all published with pagination
// @ID post-get-all-published
// @Tags Post
// @Accept json
// @Produce json
// @Param current_page query int false "Number of current page"
// @Param posts_per_page query int false "Number of posts count"
// @Param user_id query string false "Posts with user id"
// @Success 200 {object} response.PostPaginationResponseDto
// @Failure 500 {object} response.ErrorResponseDto
// @Router /post [get]
func (h *Handler) getAllPublishedPosts(w http.ResponseWriter, r *http.Request) {
	request := requsetdto.PostPaginationRequestDto{}
	response := responsedto.PostPaginationResponseDto{}

	request.FromRequest(r)

	currentCacheKey := fmt.Sprintf(GetAllPublishedPostsCacheKey, request.CurrentPage, request.CountPerPage, request.UserID)

	founded, err := h.getFromCache(r.Context(), currentCacheKey, &response)
	if err != nil {

		log.Print(err)

		errorResp := responsedto.NewErrorResponseDto(http.StatusInternalServerError, errors.ErrInternal.Error())
		errorRespond(w, r, errorResp)

		return
	}

	if !founded {

		opt := request.TransformToObject()
		pagination, err := h.post.GetAllPublishedPaginate(r.Context(), opt)
		if err != nil {

			log.Print(err)

			errorResp := responsedto.NewErrorResponseDto(http.StatusInternalServerError, errors.ErrInternal.Error())
			errorRespond(w, r, errorResp)

			return
		}

		response.TransformFromObject(pagination)
		h.saveToCache(r.Context(), currentCacheKey, response)
	}

	respond(w, r, http.StatusOK, response)
}

// @Summary Get single post
// @Description Get single post by id
// @ID post-get-single
// @Tags Post
// @Accept json
// @Produce json
// @Param id path string string "Post with id"
// @Success 200 {object} response.PostResponseDto
// @Failure 404 {object} response.ErrorResponseDto
// @Failure 500 {object} response.ErrorResponseDto
// @Router /post/{id} [get]
func (h *Handler) getSinglePost(w http.ResponseWriter, r *http.Request) {
	response := responsedto.PostResponseDto{}

	id := uuid.FromStringOrNil(chi.URLParam(r, "id"))
	currentCacheKey := fmt.Sprintf(GetSinglePostCacheKey, id)

	founded, err := h.getFromCache(r.Context(), currentCacheKey, &response)
	if err != nil {

		log.Print(err)

		errorResp := responsedto.NewErrorResponseDto(http.StatusInternalServerError, errors.ErrInternal.Error())
		errorRespond(w, r, errorResp)

		return
	}

	if !founded {

		post, err := h.post.Find(r.Context(), id)
		if err != nil {

			log.Print(err)

			var errorResp responsedto.ErrorResponseDto
			if err == repoerrors.ErrPostNotFound {
				errorResp = responsedto.NewErrorResponseDto(http.StatusNotFound, err.Error())
			} else {
				errorResp = responsedto.NewErrorResponseDto(http.StatusInternalServerError, errors.ErrInternal.Error())
			}

			errorRespond(w, r, errorResp)

			return
		}

		response.TransformFromObject(post)
		h.saveToCache(r.Context(), currentCacheKey, response)
	}

	respond(w, r, http.StatusOK, response)
}

// @Summary Create post
// @Description Create new post
// @ID post-create
// @Tags Post
// @Accept json
// @Produce json
// @Param payload body request.CreatePostRequestDto true "Post details"
// @Success 201 {object} response.PostResponseDto
// @Failure 400 {object} response.ErrorResponseDto
// @Failure 401 {object} response.ErrorResponseDto
// @Failure 403 {object} response.ErrorResponseDto
// @Failure 500 {object} response.ErrorResponseDto
// @Security ApiKeyAuth
// @Router /post [post]
func (h *Handler) createPost(w http.ResponseWriter, r *http.Request) {
	request := requsetdto.CreatePostRequestDto{}
	response := responsedto.PostResponseDto{}

	if response, err := request.FromRequest(r); err != nil {

		log.Print(err)

		errorRespond(w, r, response)

		return
	}

	opt := request.TransformToObject()
	post, err := h.post.Create(r.Context(), opt)
	if err != nil {

		log.Print(err)

		errorResp := responsedto.NewErrorResponseDto(http.StatusInternalServerError, errors.ErrInternal.Error())
		errorRespond(w, r, errorResp)

		return
	}

	response.TransformFromObject(post)
	respond(w, r, http.StatusCreated, response)
}

// @Summary Update post
// @Description Update post with id
// @ID post-update
// @Tags Post
// @Accept json
// @Produce json
// @Param id path string true "Post with id"
// @Param payload body request.UpdatePostRequestDto true "Post details"
// @Success 202 {object} response.PostResponseDto
// @Failure 400 {object} response.ErrorResponseDto
// @Failure 401 {object} response.ErrorResponseDto
// @Failure 404 {object} response.ErrorResponseDto
// @Failure 500 {object} response.ErrorResponseDto
// @Security ApiKeyAuth
// @Router /post/{id} [put]
func (h *Handler) updatePost(w http.ResponseWriter, r *http.Request) {
	request := requsetdto.UpdatePostRequestDto{}
	response := responsedto.PostResponseDto{}

	if response, err := request.FromRequest(r); err != nil {

		log.Print(err)

		errorRespond(w, r, response)

		return
	}

	opt := request.TransformToObject()
	post, err := h.post.Update(r.Context(), opt)
	if err != nil {

		log.Print(err)

		var errorResp responsedto.ErrorResponseDto
		if err == repoerrors.ErrPostNotFound {
			errorResp = responsedto.NewErrorResponseDto(http.StatusNotFound, err.Error())
		} else {
			errorResp = responsedto.NewErrorResponseDto(http.StatusInternalServerError, errors.ErrInternal.Error())
		}

		errorRespond(w, r, errorResp)

		return
	}

	response.TransformFromObject(post)
	respond(w, r, http.StatusAccepted, response)
}

// @Summary Publish post
// @Description Publish post with id
// @ID post-publish
// @Tags Post
// @Accept json
// @Produce json
// @Param id path string true "Post with id"
// @Success 202 {object} response.PostResponseDto
// @Failure 401 {object} response.ErrorResponseDto
// @Failure 404 {object} response.ErrorResponseDto
// @Failure 500 {object} response.ErrorResponseDto
// @Security ApiKeyAuth
// @Router /post/{id}/publish [get]
func (h *Handler) publishPost(w http.ResponseWriter, r *http.Request) {
	response := responsedto.PostResponseDto{}

	id := uuid.FromStringOrNil(chi.URLParam(r, "id"))
	post, err := h.post.Publish(r.Context(), id)
	if err != nil {

		log.Print(err)

		var errorResp responsedto.ErrorResponseDto
		if err == repoerrors.ErrPostNotFound {
			errorResp = responsedto.NewErrorResponseDto(http.StatusNotFound, err.Error())
		} else {
			errorResp = responsedto.NewErrorResponseDto(http.StatusInternalServerError, errors.ErrInternal.Error())
		}

		errorRespond(w, r, errorResp)

		return
	}

	response.TransformFromObject(post)
	respond(w, r, http.StatusAccepted, response)
}
