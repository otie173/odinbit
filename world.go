package main

import (
	"embed"
	"math"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

//go:embed assets/images/blocks/*.png
//go:embed assets/images/world/*.png
//go:embed assets/images/characters/*.png
//go:embed assets/images/gui/*.png
//go:embed assets/images/items/*.png
var assets embed.FS

var (
	world                             map[rl.Rectangle]Block
	trees                             []Tree
	seeds                             []Seed
	visibleBlocks                     []Block
	id                                map[int]rl.Texture2D
	item                              int
	wall                              rl.Texture2D
	wallWindow                        rl.Texture2D
	floor                             rl.Texture2D
	door                              rl.Texture2D
	chest                             rl.Texture2D
	stone1                            rl.Texture2D
	stone2                            rl.Texture2D
	stone3                            rl.Texture2D
	stone4                            rl.Texture2D
	bigStone1                         rl.Texture2D
	bigStone2                         rl.Texture2D
	bigStone3                         rl.Texture2D
	bigStone4                         rl.Texture2D
	bigStone5                         rl.Texture2D
	smallTree                         rl.Texture2D
	normalTree                        rl.Texture2D
	bigTree                           rl.Texture2D
	grass1                            rl.Texture2D
	grass2                            rl.Texture2D
	grass3                            rl.Texture2D
	grass4                            rl.Texture2D
	grass5                            rl.Texture2D
	grass6                            rl.Texture2D
	barrier                           rl.Texture2D
	lootbox                           rl.Texture2D
	bones1                            rl.Texture2D
	bones2                            rl.Texture2D
	bones3                            rl.Texture2D
	bones4                            rl.Texture2D
	bones5                            rl.Texture2D
	worldGenerated                    bool
	worldInfo                         WorldInfo
	doorOpen                          rl.Texture2D
	bigBarrel                         rl.Texture2D
	bookshelf                         rl.Texture2D
	chair                             rl.Texture2D
	closet                            rl.Texture2D
	fence1, fence2                    rl.Texture2D
	floor2, floor4                    rl.Texture2D
	lamp                              rl.Texture2D
	shelf                             rl.Texture2D
	sign                              rl.Texture2D
	smallBarrel                       rl.Texture2D
	table                             rl.Texture2D
	tombstone                         rl.Texture2D
	trash                             rl.Texture2D
	sapling                           rl.Texture2D
	resourceTick                      int = 0
	seed1Normal, seed1Big             rl.Texture2D
	seed2Small, seed2Normal, seed2Big rl.Texture2D
)

const (
	// Описание мира
	TILE_SIZE               float32 = 10.0
	WORLD_SIZE              int     = 320
	OBJECT_SPAWN_MULTIPLIER int     = 6
	GROWTH_TIME             int     = 10
	RESOURCE_SPAWN_TIME     int     = 180
	SEED_TIME               int     = 75
)

const (
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
	STAIRSDOWN
	STAIRSUP
	SAPLING
	SEED1NORMAL
	SEED1BIG
	SEED2SMALL
	SEED2NORMAL
	SEED2BIG
)

type BlockData struct {
	X         float32 `json:"x"`
	Y         float32 `json:"y"`
	TextureID int     `json:"id"`
}

type Block struct {
	img      rl.Texture2D
	rec      rl.Rectangle
	passable bool
}

type Tree struct {
	x    float32
	y    float32
	tick int
}

type Seed struct {
	x    float32
	y    float32
	tick int
}

type WorldInfo struct {
	StructuresGenerated bool `json:"structures_generated"`
	BonesGenerated      bool `json:"bones_generated"`
	BigStonesCount      int  `json:"big_stones_count"`
	SmallStonesCount    int  `json:"small_stones_count"`
	TreesCount          int  `json:"trees_count"`
	SaplingsCount       int  `json:"saplings_count"`
	SeedsCount          int  `json:"seeds_count"`
	PickaxesCount       int  `json:"pickaxes_count"`
	AxesCount           int  `json:"axes_count"`
	ShovelsCount        int  `json:"shovels_count"`
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
	id[SAPLING] = sapling
	id[SEED1NORMAL] = seed1Normal
	id[SEED1BIG] = seed1Big
	id[SEED2SMALL] = seed2Small
	id[SEED2NORMAL] = seed2Normal
	id[SEED2BIG] = seed2Big
}

func loadWorld() {
	world = make(map[rl.Rectangle]Block, 102_400)
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
	sapling = loadTexture("assets/images/world/sapling.png")
	seed1Normal = loadTexture("assets/images/world/seed1_normal.png")
	seed1Big = loadTexture("assets/images/world/seed1_big.png")
	seed2Small = loadTexture("assets/images/world/seed2_small.png")
	seed2Normal = loadTexture("assets/images/world/seed2_normal.png")
	seed2Big = loadTexture("assets/images/world/seed2_big.png")

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
	rl.UnloadTexture(sapling)
	rl.UnloadTexture(seed1Normal)
	rl.UnloadTexture(seed1Big)
	rl.UnloadTexture(seed2Small)
	rl.UnloadTexture(seed2Normal)
	rl.UnloadTexture(seed2Big)
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

func addTree(x, y float32) {
	tree := Tree{
		x:    x,
		y:    y,
		tick: 0,
	}
	trees = append(trees, tree)
	worldInfo.SaplingsCount++
	worldInfo.TreesCount++
}

func addSeed(x, y float32) {
	seed := Seed{
		x:    x,
		y:    y,
		tick: 0,
	}
	seeds = append(seeds, seed)
	worldInfo.SeedsCount++
}

func removeTree(x, y float32) {
	var newTrees []Tree

	for _, tree := range trees {
		if !(tree.x == x && tree.y == y) {
			newTrees = append(newTrees, tree)
		}
	}
	trees = newTrees
	worldInfo.SaplingsCount--
	worldInfo.TreesCount--
}

func removeSeed(x, y float32) {
	var newSeeds []Seed

	for _, seed := range seeds {
		if !(seed.x == x && seed.y == y) {
			newSeeds = append(newSeeds, seed)
		}
	}
	seeds = newSeeds
	worldInfo.SeedsCount--
}

func updateTree() {
	var treesToRemove []Tree
	for i := range trees {
		trees[i].tick++
		if trees[i].tick == GROWTH_TIME {
			saplingToTree(trees[i].x, trees[i].y)
			treesToRemove = append(treesToRemove, trees[i])
		}
	}

	for _, treeToRemove := range treesToRemove {
		removeTree(treeToRemove.x, treeToRemove.y)
	}
}

func updateSeed() {
	var seedsToRemove []Seed
	for i := range seeds {
		seeds[i].tick++
		if seeds[i].tick == SEED_TIME-40 || seeds[i].tick == SEED_TIME-15 || seeds[i].tick == SEED_TIME {
			seedToPlant(seeds[i].x, seeds[i].y, seeds[i].tick)
			if seeds[i].tick == SEED_TIME {
				seedsToRemove = append(seedsToRemove, seeds[i])
			}
		}
	}

	for _, seedToRemove := range seedsToRemove {
		removeSeed(seedToRemove.x, seedToRemove.y)
	}
}

func saplingToTree(x, y float32) {
	var treeTexture rl.Texture2D
	treeTextureNum := rand.Intn(3)

	switch treeTextureNum {
	case 0:
		treeTexture = smallTree
	case 1:
		treeTexture = normalTree
	case 2:
		treeTexture = bigTree
	}

	removeBlock(x, y)
	addBlock(treeTexture, x, y, false)
	worldInfo.TreesCount++
}

func seedToPlant(x, y float32, tick int) {
	var seedTexture rl.Texture2D
	seedTextureNum := rand.Intn(2)

	switch seedTextureNum {
	case 0:
		if tick == SEED_TIME-40 {
			seedTexture = seed2Small
		} else if tick == SEED_TIME-15 {
			seedTexture = seed1Normal
		} else if tick == SEED_TIME {
			seedTexture = seed1Big
		}
	case 1:
		if tick == SEED_TIME-40 {
			seedTexture = seed2Small
		} else if tick == SEED_TIME-15 {
			seedTexture = seed2Normal
		} else if tick == SEED_TIME {
			seedTexture = seed2Big
		}
	}

	removeBlock(x, y)
	addBlock(seedTexture, x, y, false)
}

func distanceInBlocks(playerX, playerY, blockX, blockY float32, distance float32) bool {
	dx := math.Floor(math.Abs(float64(blockX-playerX)) / float64(TILE_SIZE))
	dy := math.Floor(math.Abs(float64(blockY-playerY)) / float64(TILE_SIZE))

	return dx <= float64(distance) && dy <= float64(distance)
}

func updateWorld() {
	if !worldInfo.StructuresGenerated {
		for i := 0; i <= 6; i++ {
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

	// Очищаем slice видимых блоков
	visibleBlocks = visibleBlocks[:0]

	// Заполняем slice видимыми блоками
	for _, block := range world {
		if block.rec.X+block.rec.Width > left && block.rec.X < right &&
			block.rec.Y+block.rec.Height > top && block.rec.Y < bottom {
			visibleBlocks = append(visibleBlocks, block)
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
