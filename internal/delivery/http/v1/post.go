package v1

import (
	"net/http"

	"github.com/aintsashqa/go-simple-blog/internal/delivery/http/v1/errors"
	requsetdto "github.com/aintsashqa/go-simple-blog/internal/delivery/http/v1/request"
	responsedto "github.com/aintsashqa/go-simple-blog/internal/delivery/http/v1/response"
	repoerrors "github.com/aintsashqa/go-simple-blog/internal/repository/errors"
	"github.com/go-chi/chi"
	uuid "github.com/satori/go.uuid"
)

// @Summary Get all published
// @Description Get all published with pagination
// @ID post-get-all-published
// @Tags Post
// @Accept json
// @Produce json
// @Param current_page query int false "Number of current page"
// @Param count_per_page query int false "Number of posts count"
// @Param user_id query string false "Posts with user id"
// @Success 200 {object} response.PostPaginationResponseDto
// @Failure 500 {object} response.ErrorResponseDto
// @Router /post [get]
func (h *Handler) GetAllPublishedPosts(w http.ResponseWriter, r *http.Request) {
	request := requsetdto.PostPaginationRequestDto{}
	response := responsedto.PostPaginationResponseDto{}

	request.FromRequest(r)

	opt := request.TransformToObject()
	pagination, err := h.Service.Post.GetAllPublishedPaginate(r.Context(), opt)
	if err != nil {

		h.Service.Logger.Errorf("v1.GetAllPublishedPosts error: %s", err)

		errorResp := responsedto.NewErrorResponseDto(http.StatusInternalServerError, errors.ErrInternal.Error())
		errorRespond(w, r, errorResp)
		return
	}

	response.TransformFromObject(pagination)
	respond(w, r, http.StatusOK, response)
}

// @Summary Get all self posts
// @Description Get all self post with pagination
// @ID post-get-all-self
// @Tags Post
// @Accept json
// @Produce json
// @Param current_page query int false "Number of current page"
// @Param count_per_page query int false "Number of posts count"
// @Success 200 {object} response.PostPaginationResponseDto
// @Failure 403 {object} response.ErrorResponseDto
// @Failure 500 {object} response.ErrorResponseDto
// @Security ApiKeyAuth
// @Router /post/self [get]
func (h *Handler) GetAllSelfPosts(w http.ResponseWriter, r *http.Request) {
	request := requsetdto.SelfPostPaginationRequestDto{}
	response := responsedto.PostPaginationResponseDto{}

	if response, err := request.FromRequest(r); err != nil {

		h.Service.Logger.Errorf("v1.GetAllSelfPosts error: %s", err)

		errorRespond(w, r, response)
		return
	}

	opt := request.TransformToObject()
	pagination, err := h.Service.Post.GetAllSelfPaginate(r.Context(), opt)
	if err != nil {

		h.Service.Logger.Errorf("v1.GetAllSelfPosts error: %s", err)

		errorResp := responsedto.NewErrorResponseDto(http.StatusInternalServerError, errors.ErrInternal.Error())
		errorRespond(w, r, errorResp)
		return
	}

	response.TransformFromObject(pagination)
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
func (h *Handler) GetSinglePost(w http.ResponseWriter, r *http.Request) {
	response := responsedto.PostResponseDto{}

	id := uuid.FromStringOrNil(chi.URLParam(r, "id"))

	post, err := h.Service.Post.Find(r.Context(), id)
	if err != nil {

		h.Service.Logger.Errorf("v1.GetSinglePost error: %s", err)

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
func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	request := requsetdto.CreatePostRequestDto{}
	response := responsedto.PostResponseDto{}

	if response, err := request.FromRequest(r); err != nil {

		h.Service.Logger.Errorf("v1.CreatePost error: %s", err)

		errorRespond(w, r, response)
		return
	}

	opt := request.TransformToObject()
	post, err := h.Service.Post.Create(r.Context(), opt)
	if err != nil {

		h.Service.Logger.Errorf("v1.CreatePost error: %s", err)

		if errorResp, isValidation := ValidationErrorsHandler(err); isValidation {
			errorRespond(w, r, errorResp)
			return
		}

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
func (h *Handler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	request := requsetdto.UpdatePostRequestDto{}
	response := responsedto.PostResponseDto{}

	if response, err := request.FromRequest(r); err != nil {

		h.Service.Logger.Errorf("v1.UpdatePost error: %s", err)

		errorRespond(w, r, response)
		return
	}

	opt := request.TransformToObject()
	post, err := h.Service.Post.Update(r.Context(), opt)
	if err != nil {

		h.Service.Logger.Errorf("v1.UpdatePost error: %s", err)

		if errorResp, isValidation := ValidationErrorsHandler(err); isValidation {
			errorRespond(w, r, errorResp)
			return
		}

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
func (h *Handler) PublishPost(w http.ResponseWriter, r *http.Request) {
	response := responsedto.PostResponseDto{}

	id := uuid.FromStringOrNil(chi.URLParam(r, "id"))
	post, err := h.Service.Post.Publish(r.Context(), id)
	if err != nil {

		h.Service.Logger.Errorf("v1.PublishPost error: %s", err)

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

// @Summary Delete post
// @Description Delete post with id
// @ID post-delete
// @Tags Post
// @Accept json
// @Produce json
// @Param id path string true "Post with id"
// @Success 204
// @Failure 401 {object} response.ErrorResponseDto
// @Failure 403 {object} response.ErrorResponseDto
// @Failure 404 {object} response.ErrorResponseDto
// @Failure 500 {object} response.ErrorResponseDto
// @Security ApiKeyAuth
// @Router /post/{id} [delete]
func (h *Handler) DeletePost(w http.ResponseWriter, r *http.Request) {
	request := requsetdto.DeletePostRequestDto{}

	if response, err := request.FromRequest(r); err != nil {

		h.Service.Logger.Errorf("v1.DeletePost error: %s", err)

		errorRespond(w, r, response)
		return
	}

	opt := request.TransformToObject()
	if err := h.Service.Post.SoftDelete(r.Context(), opt); err != nil {

		h.Service.Logger.Errorf("v1.DeletePost error: %s", err)

		var errorResp responsedto.ErrorResponseDto
		if err == repoerrors.ErrPostNotFound {
			errorResp = responsedto.NewErrorResponseDto(http.StatusNotFound, err.Error())
		} else {
			errorResp = responsedto.NewErrorResponseDto(http.StatusInternalServerError, errors.ErrInternal.Error())
		}

		errorRespond(w, r, errorResp)
		return
	}

	respond(w, r, http.StatusNoContent, nil)
}
