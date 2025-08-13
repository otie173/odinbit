package tcp

import "log"

type TCPHandler struct {
}

func New() *TCPHandler {
	return &TCPHandler{}
}

func (m *TCPHandler) Load() {

}

func (m *TCPHandler) Run() {
	log.Println("Менеджер для TCP загружен")
}
