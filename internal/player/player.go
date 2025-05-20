package player

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/otie173/odinbit/internal/common"
	"github.com/otie173/odinbit/internal/entity"
	"github.com/otie173/odinbit/resource"
)

var (
	Player  entity.Entity
	players []*entity.Entity
)

func Load() {
	// TODO: Сделать сохранение и загрузку позиции и инвентаря из файла
	Player.Texture = resource.PlayerTexture
	SetTilePos(255, 255)

	// В будущем можно будет сделать при коннекте ,
	// чтобы не рисовать каждого отдельно -
	// добавить в слайс игрока нового и рисовать всех
	players = append(players, &Player)
}

func Draw() {
	for _, player := range players {
		// TODO: Нормальное рисование текстуры
		rl.DrawTexture(resource.PlayerTexture, int32(player.Pos.X), int32(player.Pos.Y), rl.White)
	}
}

func SetTilePos(x, y int) {
	Player.Pos = rl.NewVector2(float32(x*common.TileSize), float32(y*common.TileSize))
}
