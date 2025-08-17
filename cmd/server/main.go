package main

import (
	"github.com/otie173/odinbit/internal/pkg/server"
)

func main() {
	server := server.New("0.0.0.0:8080")
	server.Load()
	server.Run()
}
