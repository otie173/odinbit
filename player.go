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
	player                   rl.Texture2D
	playerWalk1, playerWalk2 rl.Texture2D
	playerTexture            rl.Texture2D
	playerPosition           rl.Vector2   = rl.NewVector2(0.0, 0.0)
	targetPosition           rl.Vector2   = rl.NewVector2(0.0, 0.0)
	playerRectangle          rl.Rectangle = rl.NewRectangle(playerPosition.X, playerPosition.Y, TILE_SIZE, TILE_SIZE)
	playerFlippedRectangle   rl.Rectangle = rl.NewRectangle(playerPosition.X, playerPosition.Y, -TILE_SIZE, TILE_SIZE)
	playerDirection          bool         = false
	playerBlockDistance      int          = 3
	cam                      rl.Camera2D
	prevCamPosition          rl.Vector2
	canBuild                 bool
	lastMoveTime             time.Time
	playerHealth             int  = 3
	isWalk                   bool = false
	stepTimer                int
	stepInterval             int = 15
	currentStep              int = 0
	lastTargetPosition       rl.Vector2
	animationComplete        bool
	directionChangeTimer     int
	directionChangeDelay     int = 5
)

type Player struct {
	X                float32 `json:"x"`
	Y                float32 `json:"y"`
	TargetX          float32 `json:"target_x"`
	TargetY          float32 `json:"target_y"`
	Health           int     `json:"health"`
	WoodCount        int     `json:"wood"`
	StoneCount       int     `json:"stone"`
	MetalCount       int     `json:"metal"`
	PickaxeOpen      bool    `json:"pickaxe_open"`
	AxeOpen          bool    `json:"axe_open"`
	ShovelOpen       bool    `json:"shovel_open"`
	WallOpen         bool    `json:"wall_open"`
	WallWindowOpen   bool    `json:"wall_window_open"`
	FloorOpen        bool    `json:"floor_open"`
	DoorOpen         bool    `json:"door_open"`
	DoorOpenOpen     bool    `json:"door_open_open"`
	ChestOpen        bool    `json:"chest_open"`
	WallCount        int     `json:"wall_count"`
	WallWindowCount  int     `json:"wall_window_count"`
	FloorCount       int     `json:"floor_count"`
	DoorCount        int     `json:"door_count"`
	ChestCount       int     `json:"chest_count"`
	DoorOpenCount    int     `json:"door_open_count"`
	BigBarrelOpen    bool    `json:"big_barrel_open"`
	BookshelfOpen    bool    `json:"bookshelf_open"`
	ChairOpen        bool    `json:"chair_open"`
	ClosetOpen       bool    `json:"closet_open"`
	Fence1Open       bool    `json:"fence1_open"`
	Fence2Open       bool    `json:"fence2_open"`
	Floor2Open       bool    `json:"floor2_open"`
	Floor4Open       bool    `json:"floor4_open"`
	LampOpen         bool    `json:"lamp_open"`
	ShelfOpen        bool    `json:"shelf_open"`
	SignOpen         bool    `json:"sign_open"`
	SmallBarrelOpen  bool    `json:"small_barrel_open"`
	TableOpen        bool    `json:"table_open"`
	TrashOpen        bool    `json:"trash_open"`
	LootboxOpen      bool    `json:"lootbox_open"`
	TombstoneOpen    bool    `json:"tombstone_open"`
	SaplingOpen      bool    `json:"sapling_open"`
	SeedOpen         bool    `json:"seed_open"`
	CabbageOpen      bool    `json:"cabbage_open"`
	BigBarrelCount   int     `json:"big_barrel_count"`
	BookshelfCount   int     `json:"bookshelf_count"`
	ChairCount       int     `json:"chair_count"`
	ClosetCount      int     `json:"closet_count"`
	Fence1Count      int     `json:"fence1_count"`
	Fence2Count      int     `json:"fence2_count"`
	Floor2Count      int     `json:"floor2_count"`
	Floor4Count      int     `json:"floor4_count"`
	LampCount        int     `json:"lamp_count"`
	ShelfCount       int     `json:"shelf_count"`
	SignCount        int     `json:"sign_count"`
	SmallBarrelCount int     `json:"small_barrel_count"`
	TableCount       int     `json:"table_count"`
	TrashCount       int     `json:"trash_count"`
	LootboxCount     int     `json:"lootbox_count"`
	TombstoneCount   int     `json:"tombstone_count"`
	SaplingCount     int     `json:"sapling_count"`
	SeedCount        int     `json:"seed_count"`
	CabaggeCount     int     `json:"cabbage_count"`
}

func savePlayerFile() {
	playerData := Player{
		X: playerPosition.X, Y: playerPosition.Y, TargetX: targetPosition.X, TargetY: targetPosition.Y,
		Health: playerHealth, WoodCount: woodCount, StoneCount: stoneCount, MetalCount: metalCount,
		PickaxeOpen: pickaxeIsOpen, AxeOpen: axeIsOpen, ShovelOpen: shovelIsOpen,
		WallOpen: wallIsOpen, WallWindowOpen: wallWindowIsOpen, FloorOpen: floorIsOpen,
		DoorOpen: doorIsOpen, DoorOpenOpen: doorOpenIsOpen, ChestOpen: chestIsOpen,
		WallCount: wallCount, WallWindowCount: wallWindowCount, FloorCount: floorCount,
		DoorCount: doorCount, ChestCount: chestCount, DoorOpenCount: doorOpenCount,
		BigBarrelOpen: bigBarrelIsOpen, BookshelfOpen: bookshelfIsOpen, ChairOpen: chairIsOpen,
		ClosetOpen: closetIsOpen, Fence1Open: fence1IsOpen, Fence2Open: fence2IsOpen,
		Floor2Open: floor2IsOpen, Floor4Open: floor4IsOpen, LampOpen: lampIsOpen,
		ShelfOpen: shelfIsOpen, SignOpen: signIsOpen, SmallBarrelOpen: smallBarrelIsOpen,
		TableOpen: tableIsOpen, TrashOpen: trashIsOpen, LootboxOpen: lootboxIsOpen, TombstoneOpen: tombstoneIsOpen, SaplingOpen: saplingIsOpen, SeedOpen: seedIsOpen, CabbageOpen: cabbageIsOpen,
		BigBarrelCount: bigBarrelCount, BookshelfCount: bookshelfCount, ChairCount: chairCount,
		ClosetCount: closetCount, Fence1Count: fence1Count, Fence2Count: fence2Count,
		Floor2Count: floor2Count, Floor4Count: floor4Count, LampCount: lampCount,
		ShelfCount: shelfCount, SignCount: signCount, SmallBarrelCount: smallBarrelCount,
		TableCount: tableCount, TrashCount: trashCount, LootboxCount: lootboxCount, TombstoneCount: tombstoneCount, SaplingCount: saplingCount, SeedCount: seedCount, CabaggeCount: cabbageCount,
	}
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

	bigBarrelIsOpen = playerData.BigBarrelOpen
	bookshelfIsOpen = playerData.BookshelfOpen
	chairIsOpen = playerData.ChairOpen
	closetIsOpen = playerData.ClosetOpen
	fence1IsOpen = playerData.Fence1Open
	fence2IsOpen = playerData.Fence2Open
	floor2IsOpen = playerData.Floor2Open
	floor4IsOpen = playerData.Floor4Open
	lampIsOpen = playerData.LampOpen
	shelfIsOpen = playerData.ShelfOpen
	signIsOpen = playerData.SignOpen
	smallBarrelIsOpen = playerData.SmallBarrelOpen
	tableIsOpen = playerData.TableOpen
	trashIsOpen = playerData.TrashOpen
	lootboxIsOpen = playerData.LootboxOpen
	tombstoneIsOpen = playerData.TombstoneOpen
	saplingIsOpen = playerData.SaplingOpen
	seedIsOpen = playerData.SeedOpen
	cabbageIsOpen = playerData.CabbageOpen

	bigBarrelCount = playerData.BigBarrelCount
	bookshelfCount = playerData.BookshelfCount
	chairCount = playerData.ChairCount
	closetCount = playerData.ClosetCount
	fence1Count = playerData.Fence1Count
	fence2Count = playerData.Fence2Count
	floor2Count = playerData.Floor2Count
	floor4Count = playerData.Floor4Count
	lampCount = playerData.LampCount
	shelfCount = playerData.ShelfCount
	signCount = playerData.SignCount
	smallBarrelCount = playerData.SmallBarrelCount
	tableCount = playerData.TableCount
	trashCount = playerData.TrashCount
	lootboxCount = playerData.LootboxCount
	tombstoneCount = playerData.TombstoneCount
	saplingCount = playerData.SaplingCount
	seedCount = playerData.SeedCount
	cabbageCount = playerData.CabaggeCount
}

func loadPlayer() {
	player = loadTexture("assets/images/characters/player.png")
	playerWalk1 = loadTexture("assets/images/characters/playerWalk1.png")
	playerWalk2 = loadTexture("assets/images/characters/playerWalk2.png")
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
	// Вычисляем разницу между текущей и целевой позицией
	dx := targetPosition.X - playerPosition.X
	dy := targetPosition.Y - playerPosition.Y

	// Определяем, нужно ли двигаться
	const minMoveDistance = 0.05
	distanceSquared := dx*dx + dy*dy

	if distanceSquared >= minMoveDistance*minMoveDistance {
		// Проверяем, изменилась ли целевая позиция
		if targetPosition != lastTargetPosition {
			directionChangeTimer++
			if directionChangeTimer >= directionChangeDelay {
				if currentStep == 0 || animationComplete {
					currentStep = 1
				} else {
					currentStep = 3 - currentStep // Переключаем между 1 и 2
				}
				stepTimer = 0
				animationComplete = false
				lastTargetPosition = targetPosition
				directionChangeTimer = 0
			}
		} else {
			directionChangeTimer = 0
		}

		// Выполняем интерполяцию
		lerpFactor := 0.075 * rl.GetFrameTime() * float32(monitorFPS)
		playerPosition.X += dx * lerpFactor
		playerPosition.Y += dy * lerpFactor
		isWalk = true
	} else {
		// Если расстояние меньше minMoveDistance, устанавливаем точную позицию
		playerPosition = targetPosition
		isWalk = false
		currentStep = 0
		animationComplete = true
		directionChangeTimer = 0
	}
}

func updatePlayerTexture() {
	if isWalk && !animationComplete {
		stepTimer++
		if stepTimer >= stepInterval {
			if currentStep == 0 {
				currentStep = 1
			} else if currentStep == 1 {
				currentStep = 2
			} else {
				currentStep = 0
				animationComplete = true
			}
			stepTimer = 0
		}

		switch currentStep {
		case 1:
			playerTexture = playerWalk1
		case 2:
			playerTexture = playerWalk2
		default:
			playerTexture = player
		}
	} else {
		playerTexture = player
	}
}

func canMoveAgain() bool {
	const moveDelay time.Duration = 150 * time.Millisecond
	return time.Since(lastMoveTime) >= moveDelay
}

func drawPlayer() {
	if playerDirection {
		rl.DrawTextureRec(playerTexture, playerRectangle, rl.NewVector2(playerPosition.X, playerPosition.Y), rl.White)
	} else if !playerDirection {
		rl.DrawTextureRec(playerTexture, playerFlippedRectangle, rl.NewVector2(playerPosition.X, playerPosition.Y), rl.White)
	}
}
