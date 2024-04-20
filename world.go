package main

import (
	"embed"
	"encoding/json"
	"log"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"

	rl "github.com/gen2brain/raylib-go/raylib"
)

//go:embed assets/images/blocks/*.png
//go:embed assets/images/world/*.png
//go:embed assets/images/characters/*.png
//go:embed assets/images/gui/*.png
//go:embed assets/images/items/*.png
var assets embed.FS

var (
	world          map[rl.Rectangle]Block
	id             map[int]rl.Texture2D
	item           int
	wall           rl.Texture2D
	floor          rl.Texture2D
	door           rl.Texture2D
	chest          rl.Texture2D
	stone1         rl.Texture2D
	stone2         rl.Texture2D
	stone3         rl.Texture2D
	stone4         rl.Texture2D
	smallTree      rl.Texture2D
	normalTree     rl.Texture2D
	bigTree        rl.Texture2D
	grass1         rl.Texture2D
	grass2         rl.Texture2D
	grass3         rl.Texture2D
	grass4         rl.Texture2D
	grass5         rl.Texture2D
	grass6         rl.Texture2D
	barrier        rl.Texture2D
	worldGenerated bool
)

const (
	// Описание мира
	TILE_SIZE               float32 = 10.0
	WORLD_SIZE              int     = 384
	OBJECT_SPAWN_MULTIPLIER int     = 6

	// Перечисление для строительных блоков
	WALL = iota
	FLOOR
	DOOR
	CHEST
	SMALL_TREE
	NORMAL_TREE
	BIG_TREE
	STONE1
	STONE2
	STONE3
	STONE4
	GRASS1
	GRASS2
	GRASS3
	GRASS4
	GRASS5
	GRASS6
	BARRIER
)

type BlockData struct {
	X         float32 `json:"x"`
	Y         float32 `json:"y"`
	Passable  bool    `json:"passable"`
	TextureID int     `json:"id"`
}

type Block struct {
	img      rl.Texture2D
	rec      rl.Rectangle
	passable bool
}

type WorldInfo struct {
	Version string `json:"version"`
}

func loadID() {
	id[WALL] = wall
	id[FLOOR] = floor
	id[DOOR] = door
	id[CHEST] = chest
	id[SMALL_TREE] = smallTree
	id[NORMAL_TREE] = normalTree
	id[BIG_TREE] = bigTree
	id[STONE1] = stone1
	id[STONE2] = stone2
	id[STONE3] = stone3
	id[STONE4] = stone4
	id[GRASS1] = grass1
	id[GRASS2] = grass2
	id[GRASS3] = grass3
	id[GRASS4] = grass4
	id[GRASS5] = grass5
	id[GRASS6] = grass6
	id[BARRIER] = barrier
}

func getOdinbitPath() string {
	var odinbitPath string

	switch runtime.GOOS {
	case "windows":
		// Windows: используем AppData\Roaming
		appData := os.Getenv("APPDATA")
		if appData == "" {
			log.Fatal("Переменная окружения APPDATA не установлена")
		}
		odinbitPath = filepath.Join(appData, "odinbit")
	case "darwin":
		// macOS: используем Library/Application Support
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("Не удалось получить домашнюю директорию: %v", err)
		}
		odinbitPath = filepath.Join(home, "Library", "Application Support", "odinbit")
	case "linux":
		// Linux: используем /home/local/.share (обычно это ~/.local/share)
		xdgDataHome := os.Getenv("XDG_DATA_HOME")
		if xdgDataHome == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				log.Fatalf("Не удалось получить домашнюю директорию: %v", err)
			}
			xdgDataHome = filepath.Join(home, ".local", "share")
		}
		odinbitPath = filepath.Join(xdgDataHome, "odinbit")
	default:
		log.Fatalf("Неподдерживаемая операционная система: %s", runtime.GOOS)
	}

	// Проверка существования папки odinbit и её создание при отсутствии
	if _, err := os.Stat(odinbitPath); os.IsNotExist(err) {
		err := os.MkdirAll(odinbitPath, os.ModePerm)
		if err != nil {
			log.Fatalf("Не удалось создать папку odinbit: %v", err)
		}
	}

	return odinbitPath
}

func saveWorldInfo() {
	worldInfo := WorldInfo{Version: "indev06042024"}
	jsonData, err := json.Marshal(worldInfo)
	if err != nil {
		log.Fatalf("Не удалось преобразовать информацию мира: %v", err)
	}

	odinbitPath := getOdinbitPath()
	worldInfoPath := filepath.Join(odinbitPath, "world_info.json")

	err = os.WriteFile(worldInfoPath, jsonData, 0644)
	if err != nil {
		log.Fatalf("Не удалось сохранить информацию о мире: %v", err)
	}
}

func saveWorldFile() {
	blocksData := []BlockData{}
	var targetID int

	for rect, block := range world {
		for key, texture := range id {
			if block.img == texture {
				targetID = key
			}
		}

		blocksData = append(blocksData, BlockData{
			X:         rect.X,
			Y:         rect.Y,
			TextureID: targetID,
			Passable:  block.passable,
		})
	}

	worldData, err := json.Marshal(blocksData)
	if err != nil {
		log.Fatalf("Не удалось преобразовать данные мира: %v", err)
	}

	odinbitPath := getOdinbitPath()
	worldDataPath := filepath.Join(odinbitPath, "world_data.json")

	err = os.WriteFile(worldDataPath, worldData, 0644)
	if err != nil {
		log.Fatalf("Не удалось сохранить мир: %v", err)
	}
}

func checkWorldFile() bool {
	odinbitPath := getOdinbitPath()
	worldDataPath := filepath.Join(odinbitPath, "world_data.json")

	_, err := os.Stat(worldDataPath)
	return !os.IsNotExist(err)
}

func loadWorldFile() map[rl.Rectangle]Block {
	odinbitPath := getOdinbitPath()
	worldDataPath := filepath.Join(odinbitPath, "world_data.json")

	jsonData, err := os.ReadFile(worldDataPath)
	if err != nil {
		log.Fatalf("Ошибка при чтении файла: %v", err)
	}

	var blocksData []BlockData
	err = json.Unmarshal(jsonData, &blocksData)
	if err != nil {
		log.Fatalf("Ошибка при десериализации данных: %v", err)
	}

	world := make(map[rl.Rectangle]Block)
	for _, data := range blocksData {
		rect := rl.Rectangle{
			X:      data.X,
			Y:      data.Y,
			Width:  10.0,
			Height: 10.0,
		}
		world[rect] = Block{img: id[data.TextureID], rec: rect, passable: data.Passable}
	}

	return world
}

func loadTexture(fileName string) rl.Texture2D {
	fileBytes, err := assets.ReadFile(fileName)
	if err != nil {
		log.Fatalf("Не удалость прочитать embed файл: %v", err)
	}

	image := rl.LoadImageFromMemory(".png", fileBytes, int32(len(fileBytes)))
	texture := rl.LoadTextureFromImage(image)
	rl.UnloadImage(image)

	return texture
}

func loadWorld() {
	world = make(map[rl.Rectangle]Block, 65_536)
	id = make(map[int]rl.Texture2D, 256)
	wall = loadTexture("assets/images/blocks/wall.png")
	floor = loadTexture("assets/images/blocks/floor.png")
	door = loadTexture("assets/images/blocks/door.png")
	chest = loadTexture("assets/images/blocks/chest.png")
	smallTree = loadTexture("assets/images/world/small_tree.png")
	stone1 = loadTexture("assets/images/world/stone1.png")
	stone2 = loadTexture("assets/images/world/stone2.png")
	stone3 = loadTexture("assets/images/world/stone3.png")
	stone4 = loadTexture("assets/images/world/stone4.png")
	normalTree = loadTexture("assets/images/world/normal_tree.png")
	bigTree = loadTexture("assets/images/world/big_tree.png")
	grass1 = loadTexture("assets/images/world/grass1.png")
	grass2 = loadTexture("assets/images/world/grass2.png")
	grass3 = loadTexture("assets/images/world/grass3.png")
	grass4 = loadTexture("assets/images/world/grass4.png")
	grass5 = loadTexture("assets/images/world/grass5.png")
	grass6 = loadTexture("assets/images/world/grass6.png")
	barrier = loadTexture("assets/images/blocks/barrier.png")

	rl.SetTextureFilter(smallTree, rl.TextureFilterNearest)
	rl.SetTextureFilter(stone1, rl.TextureFilterNearest)
	rl.SetTextureFilter(stone2, rl.TextureFilterNearest)
	rl.SetTextureFilter(stone3, rl.TextureFilterNearest)
	rl.SetTextureFilter(stone4, rl.TextureFilterNearest)
	rl.SetTextureFilter(normalTree, rl.TextureFilterNearest)
	rl.SetTextureFilter(bigTree, rl.TextureFilterNearest)
	rl.SetTextureFilter(grass1, rl.TextureFilterNearest)
	rl.SetTextureFilter(grass2, rl.TextureFilterNearest)
	rl.SetTextureFilter(grass3, rl.TextureFilterNearest)
	rl.SetTextureFilter(grass4, rl.TextureFilterNearest)
	rl.SetTextureFilter(grass5, rl.TextureFilterNearest)
	rl.SetTextureFilter(grass6, rl.TextureFilterNearest)
	rl.SetTextureFilter(barrier, rl.TextureFilterNearest)

	// Установка id для блоков
	loadID()
}

func unloadWorld() {
	rl.UnloadTexture(wall)
	rl.UnloadTexture(floor)
	rl.UnloadTexture(chest)
	rl.UnloadTexture(smallTree)
	rl.UnloadTexture(stone1)
	rl.UnloadTexture(stone2)
	rl.UnloadTexture(stone3)
	rl.UnloadTexture(stone4)
	rl.UnloadTexture(normalTree)
	rl.UnloadTexture(bigTree)
	rl.UnloadTexture(grass1)
	rl.UnloadTexture(grass2)
	rl.UnloadTexture(grass3)
	rl.UnloadTexture(grass4)
	rl.UnloadTexture(grass5)
	rl.UnloadTexture(grass6)
	rl.UnloadTexture(barrier)
}

func addBlock(img rl.Texture2D, x, y float32, passable bool) {
	block := Block{
		img:      img,
		rec:      rl.NewRectangle(x*TILE_SIZE, y*TILE_SIZE, TILE_SIZE, TILE_SIZE),
		passable: passable,
	}
	world[block.rec] = block
}

func removeBlock(x, y float32) {
	delete(world, rl.NewRectangle(x*TILE_SIZE, y*TILE_SIZE, TILE_SIZE, TILE_SIZE))
}

func generateBarrier() {
	// Генерация верхней и нижней границы
	for x := -WORLD_SIZE / 2; x <= WORLD_SIZE/2; x++ {
		addBlock(barrier, float32(x), float32(-WORLD_SIZE/2), false) // Верхняя граница
		addBlock(barrier, float32(x), float32(WORLD_SIZE/2), false)  // Нижняя граница
	}
	// Генерация левой и правой границы
	for y := -WORLD_SIZE / 2; y <= WORLD_SIZE/2; y++ {
		addBlock(barrier, float32(-WORLD_SIZE/2), float32(y), false) // Левая граница
		addBlock(barrier, float32(WORLD_SIZE/2), float32(y), false)  // Правая граница
	}
}

func generateStructure(x, y, structure int) {
	switch structure {
	case 1:
		addBlock(wall, float32(x), float32(y), false)
		addBlock(wall, float32(x+1), float32(y), false)
		addBlock(wall, float32(x), float32(y+1), false)
		addBlock(floor, float32(x+1), float32(y+1), true)
		addBlock(floor, float32(x+1), float32(y+2), true)
		addBlock(chest, float32(x+2), float32(y+2), false)
		addBlock(floor, float32(x+3), float32(y+2), true)
		addBlock(wall, float32(x+4), float32(y+1), false)
		addBlock(floor, float32(x+3), float32(y+4), true)
	case 2:
		addBlock(wall, float32(x), float32(y), false)
		addBlock(door, float32(x+1), float32(y), true)
		addBlock(wall, float32(x+2), float32(y), false)
		addBlock(wall, float32(x), float32(y+1), false)
		addBlock(floor, float32(x+1), float32(y+1), true)
		addBlock(floor, float32(x+2), float32(y+1), true)
		addBlock(floor, float32(x), float32(y+2), true)
		addBlock(floor, float32(x+1), float32(y+2), true)
		addBlock(chest, float32(x+2), float32(y+2), false)
		addBlock(wall, float32(x+1), float32(y+3), false)
		addBlock(floor, float32(x+2), float32(y+3), true)
	case 3:
		addBlock(floor, float32(x), float32(y), true)
		addBlock(floor, float32(x+1), float32(y), true)
		addBlock(floor, float32(x+3), float32(y), true)
		addBlock(wall, float32(x), float32(y+1), false)
		addBlock(floor, float32(x+1), float32(y+1), true)
		addBlock(floor, float32(x+2), float32(y+1), true)
		addBlock(wall, float32(x+3), float32(y+1), false)
		addBlock(wall, float32(x), float32(y+2), false)
		addBlock(floor, float32(x+1), float32(y+2), true)
		addBlock(chest, float32(x+2), float32(y+2), false)
		addBlock(door, float32(x+3), float32(y+2), true)
		addBlock(floor, float32(x+4), float32(y+2), true)
		addBlock(wall, float32(x+1), float32(y+3), false)
		addBlock(wall, float32(x+2), float32(y+3), false)
		addBlock(wall, float32(x+3), float32(y+3), false)

	}
}

func generateTree(x, y float32) {
	// Генерация случайного номера изображения дерева
	treeImg := rand.Intn(3) + 1
	// Постановка дерева на карту в зависимости от номера текстуры
	switch treeImg {
	case 1:
		addBlock(smallTree, float32(x), float32(y), false)
	case 2:
		addBlock(normalTree, float32(x), float32(y), false)
	case 3:
		addBlock(bigTree, float32(x), float32(y), false)
	}
}

func distanceInBlocks(playerX, playerY, blockX, blockY float32, distance float32) bool {
	dx := math.Floor(math.Abs(float64(blockX-playerX)) / float64(TILE_SIZE))
	dy := math.Floor(math.Abs(float64(blockY-playerY)) / float64(TILE_SIZE))

	return dx <= float64(distance) && dy <= float64(distance)
}

func generateStone(x, y float32) {
	// Генерация случайного номера изображения камня
	stoneImg := rand.Intn(4) + 1
	// Постановка камня на карту в зависимости от номера текстуры
	switch stoneImg {
	case 1:
		addBlock(stone1, float32(x), float32(y), false)
	case 2:
		addBlock(stone2, float32(x), float32(y), false)
	case 3:
		addBlock(stone3, float32(x), float32(y), false)
	case 4:
		addBlock(stone4, float32(x), float32(y), false)
	}
}

func generateGrass(x, y float32) {
	// Генерация шанса спавна травы
	chance := rand.Intn(100) + 1

	// Генерация случайного номера изображения травы
	if chance < 20 {
		grassImage := rand.Intn(6) + 1
		switch grassImage {
		case 1:
			addBlock(grass1, x, y, true)
		case 2:
			addBlock(grass2, x, y, true)
		case 3:
			addBlock(grass3, x, y, true)
		case 4:
			addBlock(grass4, x, y, true)
		case 5:
			addBlock(grass5, x, y, true)
		case 6:
			addBlock(grass6, x, y, true)
		}
	}
	// Поставновка травы на карту в зависимости от номера текстуры

}

func generateWorld() {
	// Генерация травы
	for x := -WORLD_SIZE / 2; x <= WORLD_SIZE/2; x++ {
		for y := -WORLD_SIZE / 2; y <= WORLD_SIZE/2; y++ {
			generateGrass(float32(x), float32(y))
		}
	}

	// Генерация данжа1
	x1 := rand.Intn(WORLD_SIZE+1) - WORLD_SIZE/2
	y1 := rand.Intn(WORLD_SIZE+1) - WORLD_SIZE/2
	generateStructure(x1, y1, 1)

	//Генерация данжа2
	x2 := rand.Intn(WORLD_SIZE+1) - WORLD_SIZE/2
	y2 := rand.Intn(WORLD_SIZE+1) - WORLD_SIZE/2
	generateStructure(x2, y2, 2)

	// Генерация данжа3
	x3 := rand.Intn(WORLD_SIZE+1) - WORLD_SIZE/2
	y3 := rand.Intn(WORLD_SIZE+1) - WORLD_SIZE/2
	generateStructure(x3, y3, 3)

	// Генерация деревьев
	for i := 0; i < WORLD_SIZE*OBJECT_SPAWN_MULTIPLIER; i++ {
		x := rand.Intn(WORLD_SIZE+1) - WORLD_SIZE/2
		y := rand.Intn(WORLD_SIZE+1) - WORLD_SIZE/2
		generateTree(float32(x), float32(y))
	}

	// Генерация камней
	for i := 0; i < WORLD_SIZE*OBJECT_SPAWN_MULTIPLIER; i++ {
		x := rand.Intn(WORLD_SIZE+1) - WORLD_SIZE/2
		y := rand.Intn(WORLD_SIZE+1) - WORLD_SIZE/2
		generateStone(float32(x), float32(y))
	}

	// Генерация барьера вокруг мира
	generateBarrier()

	// Установка флага о завершении генерации мира в значение true
	worldGenerated = true
}

func updateVisibleArea(cam rl.Camera2D, screenWidth, screenHeight int) (left, right, top, bottom float32) {
	zoomFactor := 1 / cam.Zoom
	halfScreenWidth := float32(screenWidth) * 0.5 * zoomFactor
	halfScreenHeight := float32(screenHeight) * 0.5 * zoomFactor

	left = cam.Target.X - halfScreenWidth
	right = cam.Target.X + halfScreenWidth
	top = cam.Target.Y - halfScreenHeight
	bottom = cam.Target.Y + halfScreenHeight
	return
}

func drawWorld() {
	screenWidth, screenHeight := rl.GetScreenWidth(), rl.GetScreenHeight()
	left, right, top, bottom := updateVisibleArea(cam, screenWidth, screenHeight)

	for _, block := range world {
		blockRight := block.rec.X + block.rec.Width
		blockBottom := block.rec.Y + block.rec.Height

		if block.rec.X < right && blockRight > left &&
			block.rec.Y < bottom && blockBottom > top {
			rl.DrawTextureRec(block.img, block.rec, rl.NewVector2(block.rec.X, block.rec.Y), rl.White)
		}
	}
}
