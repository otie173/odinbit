package player

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/otie173/odinbit/internal/client/texture"
	"github.com/otie173/odinbit/internal/server/common"
)

var (
	GamePlayer = Player{
		Name:     "otie173",
		CurrentX: common.WorldSize / 2,
		TargetX:  common.WorldSize / 2,
		CurrentY: common.WorldSize / 2,
		TargetY:  common.WorldSize / 2,
	}
)

const (
	LocalPlayer PlayerType = iota
	NetworkPlayer
)

type PlayerType byte

type Player struct {
	Name              string
	CurrentX, TargetX float32
	CurrentY, TargetY float32
}

func ChangePos(x, y float32) {
	GamePlayer.TargetX = x
	GamePlayer.TargetY = y
}

func UpdatePos() {
	dx := GamePlayer.TargetX - GamePlayer.CurrentX
	dy := GamePlayer.TargetY - GamePlayer.CurrentY

	duration := 0.4
	step := rl.GetFrameTime() / float32(duration)

	GamePlayer.CurrentX += dx * step
	GamePlayer.CurrentY += dy * step
}

func DrawPlayer() {
	playerVec := rl.NewVector2(GamePlayer.CurrentX*12, GamePlayer.CurrentY*12)
	playerRec := rl.NewRectangle(0, 0, 12, 12)
	rl.DrawTextureRec(texture.PlayerTexture, playerRec, playerVec, rl.White)
	//rl.DrawTexture(texture.PlayerTexture, int32(GamePlayer.CurrentX*12), int32(GamePlayer.CurrentY*12), rl.White)
}
