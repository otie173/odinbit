package main

import (
	"github.com/otie173/odinbit/internal/client/pkg/client"
	"github.com/otie173/odinbit/resources"
)

func main() {

	c := client.New("Odinbit", 1920, 1080)
	c.Load()
	resources.Load()
	c.Run()
}
