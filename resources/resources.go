package resources

import (
	"github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	Pypx, PypxB rl.Font
)

func Load() {
	raygui.LoadStyle("resources/ui/styles/odinbit_style.rgs")
	Pypx = rl.LoadFont("resources/ui/fonts/pypx.ttf")
	PypxB = rl.LoadFont("resources/ui/fonts/pypx-B.ttf")
}
