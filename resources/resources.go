package resources

import (
	"log"

	"github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	Pypx, PypxB rl.Font
)

func Load() {
	raygui.LoadStyle("resources/ui/styles/odinbit_style.rgs")
	barrier := rl.LoadTexture("resources/textures/barrier.png")
	log.Println(barrier)
	Pypx = rl.LoadFont("resources/ui/fonts/pypx.ttf")
	PypxB = rl.LoadFont("resources/ui/fonts/pypx-B.ttf")
}
