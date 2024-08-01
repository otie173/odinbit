package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"runtime"

	rl "github.com/gen2brain/raylib-go/raylib"
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
			Width:  TILE_SIZE,
			Height: TILE_SIZE,
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
