package main

import (
	"github.com/otie173/odinbit/internal/client/pkg/server"
)

func main() {
	server := server.New()
	server.Load()
	server.Run()
}
