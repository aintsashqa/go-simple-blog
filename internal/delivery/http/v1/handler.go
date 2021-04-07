package v1

import (
	"github.com/aintsashqa/go-simple-blog/internal/service"
	"github.com/go-chi/chi"
)

type Handler struct {
	user service.User
	post service.Post
}

func NewHandler(user service.User, post service.Post) *Handler {
	return &Handler{user: user, post: post}
}

func (h *Handler) Init(r chi.Router) {
	r.Route("/v1", func(r chi.Router) {

		r.Route("/user", func(r chi.Router) {
			r.Post("/sign-up", h.SignUp)
			r.Post("/sign-in", h.SignIn)
			r.Get("/{id}", h.GetSingleUser)

			r.Group(func(r chi.Router) {
				r.Use(h.authenticateMiddleware)
				r.Put("/{id}", h.UpdateUser)
			})
		})

		r.Route("/post", func(r chi.Router) {
			r.Get("/", h.GetAllPublishedPosts)
			r.Get("/{id}", h.GetSinglePost)

			r.Group(func(r chi.Router) {
				r.Use(h.authenticateMiddleware)
				r.Post("/", h.CreatePost)
				r.Put("/{id}", h.UpdatePost)
				r.Get("/{id}/publish", h.PublishPost)
			})
		})
	})
}
