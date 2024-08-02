package main

import (
	"fmt"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

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
		// генерация дома 5x5
		addBlock(wall, float32(x), float32(y), false)
		addBlock(floor, float32(x+2), float32(y), true)
		addBlock(floor, float32(x+4), float32(y), true)
		addBlock(floor, float32(x), float32(y+1), true)
		addBlock(wallWindow, float32(x+3), float32(y+1), false)
		addBlock(wall, float32(x), float32(y+2), false)
		addBlock(floor, float32(x+2), float32(y+2), true)
		addBlock(chest, float32(x+4), float32(y+2), false)
		addBlock(floor, float32(x+1), float32(y+3), true)
		addBlock(door, float32(x+3), float32(y+3), true)
		addBlock(fence1, float32(x), float32(y+4), false)
		addBlock(floor, float32(x+2), float32(y+4), true)
		addBlock(floor, float32(x+4), float32(y+4), true)
	case 2:
		// генерация домика 6x6
		addBlock(wall, float32(x), float32(y), false)
		addBlock(wall, float32(x+4), float32(y), false)
		addBlock(floor, float32(x+1), float32(y+1), true)
		addBlock(shelf, float32(x+5), float32(y+1), false)
		addBlock(floor, float32(x+2), float32(y+2), true)
		addBlock(floor, float32(x+4), float32(y+2), true)
		addBlock(table, float32(x), float32(y+3), false)
		addBlock(chair, float32(x+3), float32(y+3), false)
		addBlock(closet, float32(x+5), float32(y+3), false)
		addBlock(floor, float32(x+1), float32(y+4), true)
		addBlock(trash, float32(x+4), float32(y+4), false)
		addBlock(lootbox, float32(x+2), float32(y+5), false)
		addBlock(wall, float32(x+5), float32(y+5), false)
	case 3:
		// генерация кладбища
		addBlock(fence2, float32(x), float32(y), false)
		addBlock(bones1, float32(x+3), float32(y), false)
		addBlock(fence2, float32(x+6), float32(y), false)
		addBlock(tombstone, float32(x+1), float32(y+1), false)
		addBlock(wall, float32(x+4), float32(y+1), false)
		addBlock(bones3, float32(x+2), float32(y+2), false)
		addBlock(sign, float32(x+5), float32(y+2), false)
		addBlock(fence2, float32(x), float32(y+3), false)
		addBlock(smallBarrel, float32(x+3), float32(y+3), false)
		addBlock(fence2, float32(x+6), float32(y+3), false)
		addBlock(bones4, float32(x+1), float32(y+4), false)
		addBlock(lamp, float32(x+4), float32(y+4), false)
		addBlock(shovel, float32(x+2), float32(y+5), false)
		addBlock(tombstone, float32(x+5), float32(y+5), false)
		addBlock(fence2, float32(x), float32(y+6), false)
		addBlock(bones2, float32(x+3), float32(y+6), false)
		addBlock(fence2, float32(x+6), float32(y+6), false)
	}

	worldInfo.StructuresGenerated = true
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
	worldInfo.SmallStonesCount++
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
	worldInfo.BigStonesCount++
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
		worldInfo.PickaxesCount++
	case 2:
		addBlock(bones3, x, y, false)     // 0 0
		addBlock(bones2, x-1, y, false)   // -1 0
		addBlock(bones4, x-1, y+1, false) // -1 -1
		addBlock(bones2, x, y+1, false)   // 0 -1
		addBlock(axe, x+1, y+1, false)    // 1 -1
		worldInfo.AxesCount++
	case 3:
		addBlock(bones3, x, y, false)   // 0 0
		addBlock(bones2, x, y+1, false) // 0 -1
		addBlock(shovel, x-1, y, false) // -1 0
		worldInfo.ShovelsCount++
	}
	fmt.Printf("Сгенерировано на %.0f и %0.f\n", x, y)

	worldInfo.BonesGenerated = true
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
	worldInfo.TreesCount++
}

func generateWorld() {
	for x := -WORLD_SIZE / 2; x <= WORLD_SIZE/2; x++ {
		for y := -WORLD_SIZE / 2; y <= WORLD_SIZE/2; y++ {
			generateGrass(float32(x), float32(y))
		}
	}
	for i := 0; i <= 6; i++ {
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

func generateResource() {
	resourceTick++
	if resourceTick != RESOURCE_SPAWN_TIME {
		return
	}

	for i := 0; i <= 2 && worldInfo.BigStonesCount < 960; i++ {
		generateBigStone(generateRandomPosition())
	}

	for i := 0; i <= 5 && worldInfo.SmallStonesCount < 640; i++ {
		generateStone(generateRandomPosition())
	}

	for i := 0; i <= 21 && worldInfo.TreesCount < 1920; i++ {
		generateTree(generateRandomPosition())
	}

	switch {
	case worldInfo.PickaxesCount < 9:
		needItems := 9 - worldInfo.PickaxesCount
		for i := 0; i < needItems; i++ {
			x, y := generateRandomPosition()
			generateBones(x, y, 1)
		}
	case worldInfo.AxesCount < 9:
		needItems := 9 - worldInfo.AxesCount
		for i := 0; i < needItems; i++ {
			x, y := generateRandomPosition()
			generateBones(x, y, 2)
		}
	case worldInfo.ShovelsCount < 9:
		needItems := 9 - worldInfo.ShovelsCount
		for i := 0; i < needItems; i++ {
			x, y := generateRandomPosition()
			generateBones(x, y, 3)
		}
	}
	resourceTick = 0
}

func generateRandomPosition() (float32, float32) {
	for {
		x, y := float32(rand.Intn(WORLD_SIZE+1)-WORLD_SIZE/2), float32(rand.Intn(WORLD_SIZE+1)-WORLD_SIZE/2)

		value, exist := world[rl.NewRectangle(x, y, TILE_SIZE, TILE_SIZE)]
		if !exist || value.img == grass1 || value.img == grass2 || value.img == grass3 || value.img == grass4 || value.img == grass5 || value.img == grass6 {
			return x, y
		}
	}
}
