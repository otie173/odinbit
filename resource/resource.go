package resource

import (
	"embed"
	"log"
	"odinbit/utils/build"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// /go:embed block/*.png
// /go:embed sound/*.ogg
// /go:embed soundtrack/*.ogg
// /go:embed ui/*.png
//
//go:embed entity/*.png
var resource embed.FS

func LoadTexture(path string) rl.Texture2D {
	fileData, err := resource.ReadFile(path)
	if err != nil {
		if build.GetBuildType() == build.Debug {
			log.Fatalf("Error! Failed to read embed file: %v\n", err)
		}
	}

	image := rl.LoadImageFromMemory(".png", fileData, int32(len(fileData)))
	texture := rl.LoadTextureFromImage(image)
	rl.UnloadImage(image)

	return texture
}
