package http

import (
	"github.com/go-chi/chi/v5"
)

type Router struct {
	mux *chi.Mux
}

func NewRouter(mux *chi.Mux) *Router {
	return &Router{
		mux: mux,
	}
}

func (r *Router) setupRoutes(handler *Handler) {
	r.mux.Get("/ping", handler.ping)
	r.mux.Get("/textures", handler.getTextures)
}
