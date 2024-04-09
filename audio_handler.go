package main

import (
	"embed"
	"log"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
)

//go:embed assets/audio/*.wav
var sounds embed.FS

var (
	blockAction  rl.Sound
	pickupAction rl.Sound
)

func loadSound(fileName string) rl.Sound {
	// Создание временного файла
	tmpFile, err := os.CreateTemp("", "sound-*.wav")
	if err != nil {
		log.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // Удаление файла после выхода из функции

	// Чтение данных звука из встроенной файловой системы
	soundBytes, err := sounds.ReadFile(fileName)
	if err != nil {
		log.Fatalf("Failed to read embedded sound file: %v", err)
	}

	// Запись данных во временный файл
	if _, err := tmpFile.Write(soundBytes); err != nil {
		log.Fatalf("Failed to write sound to temp file: %v", err)
	}

	// Обязательно закрываем файл, чтобы данные были записаны
	if err := tmpFile.Close(); err != nil {
		log.Fatalf("Failed to close temp file: %v", err)
	}

	// Загрузка звука из временного файла
	sound := rl.LoadSound(tmpFile.Name())

	return sound
}

func loadAudio() {
	rl.InitAudioDevice()
	blockAction = loadSound("assets/audio/block_action.wav")
	pickupAction = loadSound("assets/audio/pick_up.wav")
}

func unloadAudio() {
	rl.UnloadSound(blockAction)
	rl.UnloadSound(pickupAction)
}

func soundBlockAction() {
	rl.PlaySound(blockAction)
}

func pickupResourceSound() {
	rl.PlaySound(pickupAction)
}
