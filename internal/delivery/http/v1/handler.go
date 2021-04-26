package v1

import (
	"github.com/aintsashqa/go-simple-blog/internal/service"
	"github.com/go-chi/chi"
)

type Handler struct {
	Service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{Service: service}
}

func (h *Handler) Init(r chi.Router) {
	r.Route("/v1", func(r chi.Router) {

		r.Route("/user", func(r chi.Router) {
			r.Post("/sign-up", h.SignUp)
			r.Post("/sign-in", h.SignIn)
			r.Get("/{id}", h.GetSingleUser)

			r.Group(func(r chi.Router) {
				r.Use(h.authenticateMiddleware)
				r.Get("/self", h.GetSelfUser)
				r.Put("/{id}", h.UpdateUser)
			})
		})

		r.Route("/post", func(r chi.Router) {
			r.Get("/", h.GetAllPublishedPosts)
			r.Get("/{id}", h.GetSinglePost)

			r.Group(func(r chi.Router) {
				r.Use(h.authenticateMiddleware)
				r.Post("/", h.CreatePost)
				r.Get("/self", h.GetAllSelfPosts)
				r.Put("/{id}", h.UpdatePost)
				r.Get("/{id}/publish", h.PublishPost)
				r.Delete("/{id}", h.DeletePost)
			})
		})
	})
}
