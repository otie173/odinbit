package http

import (
	"github.com/go-chi/chi/v5"
)

type Router struct {
	chi.Router
}

func NewRouter(router chi.Router) *Router {
	return &Router{
		Router: router,
	}
}

func (r *Router) setupRoutes(handler *Handler) {
	r.Get("/ping", handler.ping)
	r.Get("/textures", handler.getTextures)
}
