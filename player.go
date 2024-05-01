package main

import (
	"encoding/json"
	"log"
	"math"
	"os"
	"path/filepath"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	player                 rl.Texture2D
	playerPosition         rl.Vector2   = rl.NewVector2(0.0, 0.0)
	targetPosition         rl.Vector2   = rl.NewVector2(0.0, 0.0)
	playerRectangle        rl.Rectangle = rl.NewRectangle(playerPosition.X, playerPosition.Y, TILE_SIZE, TILE_SIZE)
	playerFlippedRectangle rl.Rectangle = rl.NewRectangle(playerPosition.X, playerPosition.Y, -TILE_SIZE, TILE_SIZE)
	playerDirection        bool         = false
	playerBlockDistance    int          = 3
	cam                    rl.Camera2D
	prevCamPosition        rl.Vector2
	canBuild               bool
	lastMoveTime           time.Time
	playerHealth           int = 3
)

type Player struct {
	X               float32 `json:"x"`
	Y               float32 `json:"y"`
	TargetX         float32 `json:"target_x"`
	TargetY         float32 `json:"target_y"`
	Health          int     `json:"health"`
	WoodCount       int     `json:"wood"`
	StoneCount      int     `json:"stone"`
	MetalCount      int     `json:"metal"`
	PickaxeOpen     bool    `json:"pickaxe_open"`
	AxeOpen         bool    `json:"axe_open"`
	ShovelOpen      bool    `json:"shovel_open"`
	WallOpen        bool    `json:"wall_open"`
	WallWindowOpen  bool    `json:"wall_window_open"`
	FloorOpen       bool    `json:"floor_open"`
	DoorOpen        bool    `json:"door_open"`
	DoorOpenOpen    bool    `json:"door_open_open"`
	ChestOpen       bool    `json:"chest_open"`
	WallCount       int     `json:"wall_count"`
	WallWindowCount int     `json:"wall_window_count"`
	FloorCount      int     `json:"floor_count"`
	DoorCount       int     `json:"door_count"`
	ChestCount      int     `json:"chest_count"`
	DoorOpenCount   int     `json:"door_open_count"`
}

func savePlayerFile() {
	playerData := Player{playerPosition.X, playerPosition.Y, targetPosition.X, targetPosition.Y, playerHealth, woodCount, stoneCount, metalCount, pickaxeIsOpen, axeIsOpen, shovelIsOpen, wallIsOpen, wallWindowIsOpen, floorIsOpen, doorIsOpen, doorOpenIsOpen, chestIsOpen, wallCount, wallWindowCount, floorCount, doorCount, chestCount, doorOpenCount}
	jsonData, err := json.Marshal(playerData)
	if err != nil {
		log.Fatalf("Не удалось преобразовать информацию игрока: %v", err)
	}

	odinbitPath := getOdinbitPath()
	playerDataPath := filepath.Join(odinbitPath, "player_data.json")

	err = os.WriteFile(playerDataPath, jsonData, 0644)
	if err != nil {
		log.Fatalf("Не удалось сохранить информацию о игроке: %v", err)
	}
}

func loadPlayerFile() {
	odinbitPath := getOdinbitPath()
	playerDataPath := filepath.Join(odinbitPath, "player_data.json")

	jsonData, err := os.ReadFile(playerDataPath)
	if err != nil {
		log.Fatalf("Ошибка при чтении файла: %v", err)
	}

	var playerData Player
	err = json.Unmarshal(jsonData, &playerData)
	if err != nil {
		log.Fatalf("Ошибка при десериализации данных: %v", err)
	}

	playerPosition = rl.NewVector2(playerData.X, playerData.Y)
	targetPosition = rl.NewVector2(playerData.TargetX, playerData.TargetY)
	cam.Target = playerPosition
	playerHealth = playerData.Health

	woodCount = playerData.WoodCount
	stoneCount = playerData.StoneCount
	metalCount = playerData.MetalCount

	pickaxeIsOpen = playerData.PickaxeOpen
	axeIsOpen = playerData.AxeOpen
	shovelIsOpen = playerData.ShovelOpen

	wallIsOpen = playerData.WallOpen
	wallWindowIsOpen = playerData.WallWindowOpen
	floorIsOpen = playerData.FloorOpen
	doorIsOpen = playerData.DoorOpen
	chestIsOpen = playerData.ChestOpen
	doorOpenIsOpen = playerData.DoorOpen

	wallCount = playerData.WallCount
	wallWindowCount = playerData.WallWindowCount
	floorCount = playerData.FloorCount
	doorCount = playerData.DoorCount
	doorOpenCount = playerData.DoorOpenCount
	chestCount = playerData.ChestCount
}

func loadPlayer() {
	player = loadTexture("assets/images/characters/player.png")
	cam = rl.NewCamera2D(rl.NewVector2(float32(rl.GetScreenWidth()/2), float32(rl.GetScreenHeight()/2)), rl.NewVector2(float32(playerPosition.X), float32(playerPosition.Y)), 0.0, 6.0)
	lastMoveTime = time.Now()
}

func unloadPlayer() {
	rl.UnloadTexture(player)
}

func roundToFixed(x float32, places int) float32 {
	shift := math.Pow(10, float64(places))
	return float32(math.Round(float64(x)*shift) / shift)
}

func updateCameraTarget(cam *rl.Camera2D, playerPosition rl.Vector2, playerRectangle rl.Rectangle) {
	targetX := playerPosition.X + playerRectangle.Width/2
	targetY := playerPosition.Y + playerRectangle.Height/2

	newX := rl.Vector2Lerp(cam.Target, rl.NewVector2(targetX, cam.Target.Y), 0.05*rl.GetFrameTime()*float32(monitorFPS)).X
	newY := rl.Vector2Lerp(cam.Target, rl.NewVector2(cam.Target.X, targetY), 0.05*rl.GetFrameTime()*float32(monitorFPS)).Y

	cam.Target.X = roundToFixed(newX, 1)
	cam.Target.Y = roundToFixed(newY, 1)
}

func updatePlayerPosition() {
	playerPosition = rl.Vector2Lerp(playerPosition, targetPosition, 0.075*rl.GetFrameTime()*float32(monitorFPS))
}

func canMoveAgain() bool {
	const moveDelay time.Duration = 150 * time.Millisecond
	return time.Since(lastMoveTime) >= moveDelay
}

func drawPlayer() {
	if playerDirection {
		rl.DrawTextureRec(player, playerRectangle, rl.NewVector2(playerPosition.X, playerPosition.Y), rl.White)
	} else if !playerDirection {
		rl.DrawTextureRec(player, playerFlippedRectangle, rl.NewVector2(playerPosition.X, playerPosition.Y), rl.White)
	}
}
