package main

import (
	"embed"
	"log"
	"math/rand"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
)

//go:embed assets/audio/*.ogg
var sounds embed.FS

//go:embed assets/music/*.ogg
var music embed.FS

var (
	blockAction        rl.Sound
	pickupAction       rl.Sound
	musicPaused        bool
	musicActive        bool
	currentSoundTrack  rl.Music
	previousSoundTrack rl.Music
	soundTrack1        rl.Music
	soundTrack2        rl.Music
	soundTrack3        rl.Music
)

const (
	SOUNDTRACK_VOLUME float32 = 0.1
)

func loadSound(fileName string) rl.Sound {
	// Создание временного файла
	tmpFile, err := os.CreateTemp("", "sound-*.ogg")
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

func loadSoundTrack(fileName string) rl.Music {
	// Создание временного файла
	tmpFile, err := os.CreateTemp("", "music-*.ogg")
	if err != nil {
		log.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // Удаление файла после выхода из функции

	// Чтение данных музыки из встроенной файловой системы
	musicBytes, err := music.ReadFile(fileName)
	if err != nil {
		log.Fatalf("Failed to read embedded music file: %v", err)
	}

	// Запись данных во временный файл
	if _, err := tmpFile.Write(musicBytes); err != nil {
		log.Fatalf("Failed to write music to temp file: %v", err)
	}

	// Обязательно закрываем файл, чтобы данные были записаны
	if err := tmpFile.Close(); err != nil {
		log.Fatalf("Failed to close temp file: %v", err)
	}

	// Загрузка музыки из временного файла
	music := rl.LoadMusicStream(tmpFile.Name())

	return music
}

func loadAudio() {
	blockAction = loadSound("assets/audio/block_action.ogg")
	pickupAction = loadSound("assets/audio/pick_up.ogg")
}

func unloadAudio() {
	rl.UnloadSound(blockAction)
	rl.UnloadSound(pickupAction)
}

func loadMusic() {
	soundTrack1 = loadSoundTrack("assets/music/001_OdeToLoneliness_v02.ogg")
	soundTrack2 = loadSoundTrack("assets/music/002_InsatiableCuriosity_v02.ogg")
	soundTrack3 = loadSoundTrack("assets/music/003_TreesandRocks_v01.ogg")
}

func unloadMusic() {
	rl.StopMusicStream(currentSoundTrack)
	rl.UnloadMusicStream(soundTrack3)
}

func soundBlockAction() {
	rl.PlaySound(blockAction)
}

func pickupResourceSound() {
	rl.PlaySound(pickupAction)
}

func playSoundTrack() {
	if !musicActive || currentSoundTrack != previousSoundTrack {
		if currentSoundTrack != previousSoundTrack {
			rl.StopMusicStream(previousSoundTrack)
			previousSoundTrack = currentSoundTrack
		}

		rl.SetMusicVolume(currentSoundTrack, SOUNDTRACK_VOLUME)
		rl.PlayMusicStream(currentSoundTrack)
		musicActive = true
	}
}

func updateMusic() {
	// Проверяем, изменилась ли сцена
	if lastScene != currentScene {
		lastScene = currentScene // Обновляем последнюю сцену

		switch currentScene {
		case TITLE:
			currentSoundTrack = soundTrack3
		case GAME:
			// Генерируем случайное число только между 1 и 2
			randomSoundTrack := rand.Intn(2) + 1
			switch randomSoundTrack {
			case 1:
				currentSoundTrack = soundTrack1
			case 2:
				currentSoundTrack = soundTrack2
			}
		}

		// Сбрасываем флаг активности музыки, чтобы playSoundTrack мог сработать
		musicActive = false
	}

	rl.UpdateMusicStream(currentSoundTrack) // Обновляем поток музыки
	playSoundTrack()                        // Воспроизводим саундтрек, если нужно
}
