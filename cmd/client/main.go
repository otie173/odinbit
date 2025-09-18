package main

import (
	"github.com/otie173/odinbit/internal/client/common"
	"github.com/otie173/odinbit/internal/pkg/client"
	"github.com/otie173/odinbit/resources"
)

func main() {

	c := client.New("Odinbit", common.ScreenWidth, common.ScreenHeight)
	c.Load()
	resources.Load()
	c.Run()
}
