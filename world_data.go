package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync/atomic"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	gameMode byte
)

const (
	SINGLEPLAYER byte = iota
	MULTIPLAYER
)

const (
	WORLD_WIDTH  = 320
	WORLD_HEIGHT = 320
	MAX_BLOCKS   = WORLD_WIDTH * WORLD_HEIGHT
	BLOCK_BITS   = 7
	BLOCK_MASK   = (1 << BLOCK_BITS) - 1
)

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
	blocks := make([]byte, (WORLD_SIZE+1)*(WORLD_SIZE+1))
	index := 0
	for y := -WORLD_SIZE / 2; y <= WORLD_SIZE/2; y++ {
		for x := -WORLD_SIZE / 2; x <= WORLD_SIZE/2; x++ {
			rect := rl.Rectangle{
				X:      float32(x) * TILE_SIZE,
				Y:      float32(y) * TILE_SIZE,
				Width:  TILE_SIZE,
				Height: TILE_SIZE,
			}

			if block, exists := world[rect]; exists {
				var textureID byte
				for id, texture := range id {
					if block.img == texture {
						textureID = byte(id)
						break
					}
				}
				if textureID >= (1 << BLOCK_BITS) {
					log.Printf("Предупреждение: ID блока %d превышает максимальное значение %d", textureID, (1<<BLOCK_BITS)-1)
					textureID = 0 // или какое-то другое значение по умолчанию
				}
				blocks[index] = textureID + 1 // Добавляем 1, чтобы 0 означало пустой блок
			} else {
				blocks[index] = 0 // Пустой блок
			}
			index++
		}
	}

	data := make([]byte, (len(blocks)*BLOCK_BITS+7)/8)
	for i, block := range blocks {
		bitIndex := i * BLOCK_BITS
		byteIndex := bitIndex / 8
		bitOffset := uint(bitIndex % 8)
		data[byteIndex] |= byte(block << bitOffset)
		if bitOffset > 8-BLOCK_BITS && byteIndex+1 < len(data) {
			data[byteIndex+1] |= byte(block >> (8 - bitOffset))
		}
	}

	odinbitPath := getOdinbitPath()
	var worldDataPath string
	if gameMode == SINGLEPLAYER {
		worldDataPath = filepath.Join(odinbitPath, "world.odn")
	}
	if gameMode == MULTIPLAYER {
		switch atomic.LoadInt32(&worldType) {
		case 0:
			worldDataPath = filepath.Join(odinbitPath, "world_send.odn")
		case 1:
			worldDataPath = filepath.Join(odinbitPath, "world_receive.odn")
		}
	}

	err := os.WriteFile(worldDataPath, data, 0644)
	if err != nil {
		log.Fatalf("Не удалось сохранить мир: %v", err)
	}
}

func loadWorldFile() map[rl.Rectangle]Block {
	odinbitPath := getOdinbitPath()
	var worldDataPath string
	if gameMode == SINGLEPLAYER {
		worldDataPath = filepath.Join(odinbitPath, "world.odn")
	}
	if gameMode == MULTIPLAYER {
		switch atomic.LoadInt32(&worldType) {
		case 0:
			worldDataPath = filepath.Join(odinbitPath, "world_send.odn")
		case 1:
			worldDataPath = filepath.Join(odinbitPath, "world_receive.odn")
		}
	}

	data, err := os.ReadFile(worldDataPath)
	if err != nil {
		log.Printf("Ошибка при чтении файла: %v", err)
		return make(map[rl.Rectangle]Block)
	}

	blocks := make([]byte, (WORLD_SIZE+1)*(WORLD_SIZE+1))
	for i := range blocks {
		bitIndex := i * BLOCK_BITS
		byteIndex := bitIndex / 8
		bitOffset := uint(bitIndex % 8)
		if byteIndex+1 < len(data) {
			blocks[i] = byte((uint16(data[byteIndex]) | uint16(data[byteIndex+1])<<8) >> bitOffset & BLOCK_MASK)
		} else if byteIndex < len(data) {
			blocks[i] = byte(uint16(data[byteIndex]) >> bitOffset & BLOCK_MASK)
		} else {
			blocks[i] = 0
		}
	}

	world := make(map[rl.Rectangle]Block)
	index := 0
	for y := -WORLD_SIZE / 2; y <= WORLD_SIZE/2; y++ {
		for x := -WORLD_SIZE / 2; x <= WORLD_SIZE/2; x++ {
			textureID := int(blocks[index]) - 1 // Вычитаем 1, чтобы вернуться к оригинальному ID
			if textureID >= 0 {                 // Загружаем только непустые блоки
				rect := rl.Rectangle{
					X:      float32(x) * TILE_SIZE,
					Y:      float32(y) * TILE_SIZE,
					Width:  TILE_SIZE,
					Height: TILE_SIZE,
				}

				passable := false
				passableBlocks := []int{DOOR, GRASS1, GRASS2, GRASS3, GRASS4, GRASS5, GRASS6, FLOOR, FLOOR2, FLOOR4, DOOROPEN}
				for _, block := range passableBlocks {
					if textureID == block {
						passable = true
						break
					}
				}
				world[rect] = Block{img: id[textureID], rec: rect, passable: passable}

				if textureID == SAPLING {
					treeExists := false
					for _, tree := range trees {
						if tree.x == float32(x) && tree.y == float32(y) {
							treeExists = true
							break
						}
					}
					if !treeExists {
						trees = append(trees, Tree{x: float32(x), y: float32(y)})
					}
				}

				if textureID == SEED1NORMAL || textureID == SEED2NORMAL || textureID == SEED2SMALL {
					seedExists := false
					for _, seed := range seeds {
						if seed.x == float32(x) && seed.y == float32(y) {
							seedExists = true
							break
						}
					}
					if !seedExists {
						seeds = append(seeds, Seed{x: float32(x), y: float32(y)})
					}
				}
			}
			index++
		}
	}

	worldGenerated = true
	if atomic.LoadInt32(&connectedToServer) == 1 && gameMode == MULTIPLAYER {
		if err := os.Remove(worldDataPath); err != nil {
			fmt.Println("Error with remove file: ", err)
		}
	}
	return world
}

func checkWorldFile(worldname string) bool {
	odinbitPath := getOdinbitPath()

	worldDataPath := filepath.Join(odinbitPath, worldname)
	_, err := os.Stat(worldDataPath)
	return !os.IsNotExist(err)
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
