package http

import (
	"net/http"

	"github.com/aintsashqa/go-simple-blog/api/swagger"
	v1 "github.com/aintsashqa/go-simple-blog/internal/delivery/http/v1"
	"github.com/aintsashqa/go-simple-blog/internal/service"
	"github.com/go-chi/chi"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Handler struct {
	Service *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		Service: services,
	}
}

func (h *Handler) Init(host string, port int) http.Handler {
	h.Service.Logger.Info("Initialize routes")

	r := chi.NewRouter()

	// Swagger docs
	h.swagger(r, host, port)

	// API
	h.api(r)

	return r
}

func (h *Handler) swagger(r chi.Router, host string, port int) {
	h.Service.Logger.Info("Initialize swagger route")

	// swagger.SwaggerInfo.Host = fmt.Sprintf("%s:%d", host, port)
	swagger.SwaggerInfo.Host = "localhost:80"

	r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL("doc.json")))
}

func (h *Handler) api(r chi.Router) {
	h.Service.Logger.Info("Initialize api routes")

	version1 := v1.NewHandler(h.Service)

	r.Route("/api", func(r chi.Router) {
		version1.Init(r)
	})
}
