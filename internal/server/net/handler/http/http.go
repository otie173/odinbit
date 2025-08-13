package http

import (
	"net/http"
)

type HTTPHandler struct {
}

func New() *HTTPHandler {
	return &HTTPHandler{}
}

func (m *HTTPHandler) Load() {

}

func (m *HTTPHandler) Run() error {
	return http.ListenAndServe("0.0.0.0:8080", nil)
}
