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
	visibleBlocks  map[rl.Rectangle]Block
	id             map[int]rl.Texture2D
	item           int
	wall           rl.Texture2D
	wallWindow     rl.Texture2D
	floor          rl.Texture2D
	door           rl.Texture2D
	chest          rl.Texture2D
	stone1         rl.Texture2D
	stone2         rl.Texture2D
	stone3         rl.Texture2D
	stone4         rl.Texture2D
	bigStone1      rl.Texture2D
	bigStone2      rl.Texture2D
	bigStone3      rl.Texture2D
	bigStone4      rl.Texture2D
	bigStone5      rl.Texture2D
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
	lootbox        rl.Texture2D
	bones1         rl.Texture2D
	bones2         rl.Texture2D
	bones3         rl.Texture2D
	bones4         rl.Texture2D
	bones5         rl.Texture2D
	worldGenerated bool
	worldInfo      WorldInfo
	doorOpen       rl.Texture2D
	bigBarrel      rl.Texture2D
	bookshelf      rl.Texture2D
	chair          rl.Texture2D
	closet         rl.Texture2D
	fence1, fence2 rl.Texture2D
	floor2, floor4 rl.Texture2D
	lamp           rl.Texture2D
	shelf          rl.Texture2D
	sign           rl.Texture2D
	smallBarrel    rl.Texture2D
	table          rl.Texture2D
	tombstone      rl.Texture2D
	trash          rl.Texture2D
)

const (
	// Описание мира
	TILE_SIZE               float32 = 10.0
	WORLD_SIZE              int     = 320
	OBJECT_SPAWN_MULTIPLIER int     = 6

	// Перечисление для строительных блоков
	WALL = iota
	WALLWINDOW
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
	BIGSTONE1
	BIGSTONE2
	BIGSTONE3
	BIGSTONE4
	BIGSTONE5
	GRASS1
	GRASS2
	GRASS3
	GRASS4
	GRASS5
	GRASS6
	BARRIER
	LOOTBOX
	BONES1
	BONES2
	BONES3
	BONES4
	BONES5
	PICKAXE
	AXE
	SHOVEL
	DOOROPEN
	BIGBARREL
	BOOKSHELF
	CHAIR
	CLOSET
	FENCE1
	FENCE2
	FLOOR2
	FLOOR4
	LAMP
	SHELF
	SIGN
	SMALLBARREL
	TABLE
	TOMBSTONE
	TRASH
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
	StructuresGenerated bool `json:"structures_generated"`
	BonesGenerated      bool `json:"bones_generated"`
}

func loadID() {
	id[WALL] = wall
	id[WALLWINDOW] = wallWindow
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
	id[BIGSTONE1] = bigStone1
	id[BIGSTONE2] = bigStone2
	id[BIGSTONE3] = bigStone3
	id[BIGSTONE4] = bigStone4
	id[BIGSTONE5] = bigStone5
	id[GRASS1] = grass1
	id[GRASS2] = grass2
	id[GRASS3] = grass3
	id[GRASS4] = grass4
	id[GRASS5] = grass5
	id[GRASS6] = grass6
	id[BARRIER] = barrier
	id[LOOTBOX] = lootbox
	id[BONES1] = bones1
	id[BONES2] = bones2
	id[BONES3] = bones3
	id[BONES4] = bones4
	id[BONES5] = bones5
	id[PICKAXE] = pickaxe
	id[AXE] = axe
	id[SHOVEL] = shovel
	id[DOOROPEN] = doorOpen
	id[BIGBARREL] = bigBarrel
	id[BOOKSHELF] = bookshelf
	id[CHAIR] = chair
	id[CLOSET] = closet
	id[FENCE1] = fence1
	id[FENCE2] = fence2
	id[FLOOR2] = floor2
	id[FLOOR4] = floor4
	id[LAMP] = lamp
	id[SHELF] = shelf
	id[SIGN] = sign
	id[SMALLBARREL] = smallBarrel
	id[TABLE] = table
	id[TOMBSTONE] = tombstone
	id[TRASH] = trash
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

	worldGenerated = true
	return world
}

func loadWorldInfo() WorldInfo {
	odinbitPath := getOdinbitPath()
	worldInfoPath := filepath.Join(odinbitPath, "world_info.json")

	jsonData, err := os.ReadFile(worldInfoPath)
	if err != nil {
		log.Fatalf("Ошибка при чтении файла: %v", err)
	}

	var worldInfoFile WorldInfo
	err = json.Unmarshal(jsonData, &worldInfoFile)
	if err != nil {
		log.Fatalf("Ошибка при десериализации данных: %v", err)
	}

	return worldInfoFile
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
	wallWindow = loadTexture("assets/images/blocks/wall_window.png")
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
	bigStone1 = loadTexture("assets/images/world/big_stone1.png")
	bigStone2 = loadTexture("assets/images/world/big_stone2.png")
	bigStone3 = loadTexture("assets/images/world/big_stone3.png")
	bigStone4 = loadTexture("assets/images/world/big_stone4.png")
	bigStone5 = loadTexture("assets/images/world/big_stone5.png")
	lootbox = loadTexture("assets/images/blocks/lootbox.png")
	bones1 = loadTexture("assets/images/world/bones1.png")
	bones2 = loadTexture("assets/images/world/bones2.png")
	bones3 = loadTexture("assets/images/world/bones3.png")
	bones4 = loadTexture("assets/images/world/bones4.png")
	bones5 = loadTexture("assets/images/world/bones5.png")
	doorOpen = loadTexture("assets/images/blocks/door_open.png")
	bigBarrel = loadTexture("assets/images/blocks/big_barrel1.png")
	bookshelf = loadTexture("assets/images/blocks/bookshelf1.png")
	chair = loadTexture("assets/images/blocks/chair.png")
	closet = loadTexture("assets/images/blocks/closet.png")
	fence1 = loadTexture("assets/images/blocks/fence1.png")
	fence2 = loadTexture("assets/images/blocks/fence2.png")
	floor2 = loadTexture("assets/images/blocks/floor2.png")
	floor4 = loadTexture("assets/images/blocks/floor4.png")
	lamp = loadTexture("assets/images/blocks/lamp.png")
	sign = loadTexture("assets/images/blocks/sign.png")
	smallBarrel = loadTexture("assets/images/blocks/small_barrel1.png")
	table = loadTexture("assets/images/blocks/table.png")
	tombstone = loadTexture("assets/images/blocks/tombstone1.png")
	trash = loadTexture("assets/images/blocks/trash.png")
	shelf = loadTexture("assets/images/blocks/shelf.png")

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
	rl.UnloadTexture(bigStone1)
	rl.UnloadTexture(bigStone2)
	rl.UnloadTexture(bigStone3)
	rl.UnloadTexture(bigStone4)
	rl.UnloadTexture(bigStone5)
	rl.UnloadTexture(bones1)
	rl.UnloadTexture(bones2)
	rl.UnloadTexture(bones3)
	rl.UnloadTexture(bones4)
	rl.UnloadTexture(bones5)
	rl.UnloadTexture(doorOpen)
	rl.UnloadTexture(bigBarrel)
	rl.UnloadTexture(bookshelf)
	rl.UnloadTexture(chair)
	rl.UnloadTexture(closet)
	rl.UnloadTexture(fence1)
	rl.UnloadTexture(fence2)
	rl.UnloadTexture(floor2)
	rl.UnloadTexture(floor4)
	rl.UnloadTexture(lamp)
	rl.UnloadTexture(shelf)
	rl.UnloadTexture(sign)
	rl.UnloadTexture(smallBarrel)
	rl.UnloadTexture(table)
	rl.UnloadTexture(tombstone)
	rl.UnloadTexture(trash)
}

func addBlock(img rl.Texture2D, x, y float32, passable bool) {
	block := Block{
		img:      img,
		rec:      rl.NewRectangle(x*TILE_SIZE, y*TILE_SIZE, TILE_SIZE, TILE_SIZE),
		passable: passable,
	}
	world[block.rec] = block

	if worldGenerated {
		updateVisibleBlocks(cam)
	}
}

func removeBlock(x, y float32) {
	delete(world, rl.NewRectangle(x*TILE_SIZE, y*TILE_SIZE, TILE_SIZE, TILE_SIZE))

	if worldGenerated {
		updateVisibleBlocks(cam)
	}
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

	worldInfo.StructuresGenerated = true
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

func generateBigStone(x, y float32) {
	stoneImg := rand.Intn(5) + 1
	switch stoneImg {
	case 1:
		addBlock(bigStone1, float32(x), float32(y), false)
	case 2:
		addBlock(bigStone2, float32(x), float32(y), false)
	case 3:
		addBlock(bigStone3, float32(x), float32(y), false)
	case 4:
		addBlock(bigStone4, float32(x), float32(y), false)
	case 5:
		addBlock(bigStone5, float32(x), float32(y), false)
	}
}

func generateGrass(x, y float32) {
	chance := rand.Intn(100) + 1
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

}

func generateBones(x, y float32, bonesPattern int) {
	switch bonesPattern {
	case 1:
		addBlock(bones1, x, y, false)      // 0 0
		addBlock(bones2, x-1, y, false)    // -1 0
		addBlock(bones2, x-1, y+1, false)  // -1 -1
		addBlock(bones4, x-2, y-1, false)  // -2 1
		addBlock(pickaxe, x-1, y-1, false) // -1 1
	case 2:
		addBlock(bones3, x, y, false)     // 0 0
		addBlock(bones2, x-1, y, false)   // -1 0
		addBlock(bones4, x-1, y+1, false) // -1 -1
		addBlock(bones2, x, y+1, false)   // 0 -1
		addBlock(axe, x+1, y+1, false)    // 1 -1
	case 3:
		addBlock(bones3, x, y, false)   // 0 0
		addBlock(bones2, x, y+1, false) // 0 -1
		addBlock(shovel, x-1, y, false) // -1 0
	}

	worldInfo.BonesGenerated = true
}

func generateWorld() {
	for x := -WORLD_SIZE / 2; x <= WORLD_SIZE/2; x++ {
		for y := -WORLD_SIZE / 2; y <= WORLD_SIZE/2; y++ {
			generateGrass(float32(x), float32(y))
		}
	}
	for i := 0; i <= 8; i++ {
		generateStructure(rand.Intn(WORLD_SIZE+1)-WORLD_SIZE/2, rand.Intn(WORLD_SIZE+1)-WORLD_SIZE/2, 1)
		generateStructure(rand.Intn(WORLD_SIZE+1)-WORLD_SIZE/2, rand.Intn(WORLD_SIZE+1)-WORLD_SIZE/2, 2)
		generateStructure(rand.Intn(WORLD_SIZE+1)-WORLD_SIZE/2, rand.Intn(WORLD_SIZE+1)-WORLD_SIZE/2, 3)
	}
	for i := 0; i < WORLD_SIZE*OBJECT_SPAWN_MULTIPLIER; i++ {
		generateTree(float32(rand.Intn(WORLD_SIZE+1)-WORLD_SIZE/2), float32(rand.Intn(WORLD_SIZE+1)-WORLD_SIZE/2))
	}
	for i := 0; i < WORLD_SIZE*(OBJECT_SPAWN_MULTIPLIER-4); i++ {
		generateStone(float32(rand.Intn(WORLD_SIZE+1)-WORLD_SIZE/2), float32(rand.Intn(WORLD_SIZE+1)-WORLD_SIZE/2))
	}
	for i := 0; i < WORLD_SIZE*(OBJECT_SPAWN_MULTIPLIER-3); i++ {
		generateBigStone(float32(rand.Intn(WORLD_SIZE+1)-WORLD_SIZE/2), float32(rand.Intn(WORLD_SIZE+1)-WORLD_SIZE/2))
	}
	for i := 0; i <= 8; i++ {
		generateBones(float32(rand.Intn(WORLD_SIZE+1)-WORLD_SIZE/2), float32(rand.Intn(WORLD_SIZE+1)-WORLD_SIZE/2), 1)
		generateBones(float32(rand.Intn(WORLD_SIZE+1)-WORLD_SIZE/2), float32(rand.Intn(WORLD_SIZE+1)-WORLD_SIZE/2), 2)
		generateBones(float32(rand.Intn(WORLD_SIZE+1)-WORLD_SIZE/2), float32(rand.Intn(WORLD_SIZE+1)-WORLD_SIZE/2), 3)
	}
	generateBarrier()
	worldGenerated = true
}

func updateWorld() {
	if !worldInfo.StructuresGenerated {
		for i := 0; i <= 8; i++ {
			generateStructure(rand.Intn(WORLD_SIZE+1)-WORLD_SIZE/2, rand.Intn(WORLD_SIZE+1)-WORLD_SIZE/2, 1)
			generateStructure(rand.Intn(WORLD_SIZE+1)-WORLD_SIZE/2, rand.Intn(WORLD_SIZE+1)-WORLD_SIZE/2, 2)
			generateStructure(rand.Intn(WORLD_SIZE+1)-WORLD_SIZE/2, rand.Intn(WORLD_SIZE+1)-WORLD_SIZE/2, 3)
		}
	}
	if !worldInfo.BonesGenerated {
		for i := 0; i <= 8; i++ {
			generateBones(float32(rand.Intn(WORLD_SIZE+1)-WORLD_SIZE/2), float32(rand.Intn(WORLD_SIZE+1)-WORLD_SIZE/2), 1)
			generateBones(float32(rand.Intn(WORLD_SIZE+1)-WORLD_SIZE/2), float32(rand.Intn(WORLD_SIZE+1)-WORLD_SIZE/2), 2)
			generateBones(float32(rand.Intn(WORLD_SIZE+1)-WORLD_SIZE/2), float32(rand.Intn(WORLD_SIZE+1)-WORLD_SIZE/2), 3)
		}
	}
}

func isCameraMoved(cam rl.Camera2D) bool {
	moved := cam.Target.X != prevCamPosition.X || cam.Target.Y != prevCamPosition.Y
	if moved {
		prevCamPosition = cam.Target
	}
	return moved
}

func updateVisibleBlocks(cam rl.Camera2D) {
	screenWidth, screenHeight := rl.GetScreenWidth(), rl.GetScreenHeight()
	left, right, top, bottom := updateVisibleArea(cam, screenWidth, screenHeight)

	// Очистите мапу видимых блоков перед заполнением
	for k := range visibleBlocks {
		delete(visibleBlocks, k)
	}

	// Заполните мапу видимыми блоками
	for _, block := range world {
		if block.rec.X+block.rec.Width > left && block.rec.X < right &&
			block.rec.Y+block.rec.Height > top && block.rec.Y < bottom {
			visibleBlocks[block.rec] = block
		}
	}
}

func updateVisibleArea(cam rl.Camera2D, screenWidth, screenHeight int) (left, right, top, bottom float32) {
	// Вычисляем центральную точку камеры
	camCenter := cam.Target

	// Вычисляем размеры видимой области на основе зума камеры и размеров экрана
	// Предполагаем, что cam.Zoom равен 1.0 при стандартном масштабе
	halfWidth := float32(screenWidth) / 2.0 / cam.Zoom
	halfHeight := float32(screenHeight) / 2.0 / cam.Zoom

	// Определяем границы видимой области
	left = camCenter.X - halfWidth
	right = camCenter.X + halfWidth
	top = camCenter.Y - halfHeight
	bottom = camCenter.Y + halfHeight

	return left, right, top, bottom
}

func drawWorld(cam rl.Camera2D) {
	if isCameraMoved(cam) {
		updateVisibleBlocks(cam)
	}

	for _, block := range visibleBlocks {
		rl.DrawTextureRec(block.img, block.rec, rl.NewVector2(block.rec.X, block.rec.Y), rl.White)
	}
}
