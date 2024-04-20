package main

import (
	"embed"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

//go:embed assets/audio/*.ogg
var sounds embed.FS

//go:embed assets/music/*.ogg
var music embed.FS

var (
	blockAction         rl.Sound
	pickupAction        rl.Sound
	musicPaused         bool
	musicActive         bool
	currentSoundTrack   rl.Music
	previousSoundTrack  rl.Music
	lastSoundTrackIndex int = -1
	lastTrackEndTime    time.Time
	pauseDuration       time.Duration
	soundTrack1         rl.Music
	soundTrack2         rl.Music
	soundTrack3         rl.Music
)

const (
	SOUNDTRACK_VOLUME float32 = 0.05
	MENU_TARCK_INDEX  int     = 2
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
	// Это сделано для того чтобы, если в сцене GAME играл soundTrack3, то
	// при выходе в сцену TITLE - soundTrack3 мог звучать снова
	if !musicActive || currentSoundTrack == soundTrack3 && currentScene == TITLE {
		fmt.Println("Я здесь")
		if currentSoundTrack != previousSoundTrack {
			rl.StopMusicStream(previousSoundTrack)
		}
		previousSoundTrack = currentSoundTrack

		rl.SetMusicVolume(currentSoundTrack, SOUNDTRACK_VOLUME)
		rl.PlayMusicStream(currentSoundTrack)
		musicActive = true
		musicPaused = false
		// Обновляем время окончания трека, учитывая длительность трека
		lastTrackEndTime = time.Now().Add(time.Second * time.Duration(rl.GetMusicTimeLength(currentSoundTrack)))
	}

	if !musicActive || currentSoundTrack != previousSoundTrack {
		if currentSoundTrack != previousSoundTrack {
			rl.StopMusicStream(previousSoundTrack)
			previousSoundTrack = currentSoundTrack
		}

		rl.SetMusicVolume(currentSoundTrack, SOUNDTRACK_VOLUME)
		rl.PlayMusicStream(currentSoundTrack)
		musicActive = true
		musicPaused = false
		// Обновляем время окончания трека, учитывая длительность трека
		lastTrackEndTime = time.Now().Add(time.Second * time.Duration(rl.GetMusicTimeLength(currentSoundTrack)))
	}
}

func pauseMusic() {
	if musicActive && !musicPaused {
		rl.PauseMusicStream(currentSoundTrack)
		musicPaused = true
	}
}

func resumeMusic() {
	if musicActive && musicPaused {
		rl.ResumeMusicStream(currentSoundTrack)
		musicPaused = false
	}
}

func selectNewTrackForGame() {
	var newTrackIndex int
	for {
		newTrackIndex = rand.Intn(3) // Предполагается, что индексы треков начинаются с 0
		if newTrackIndex != lastSoundTrackIndex {
			break
		}
	}

	fmt.Printf("Новый трек: %d\n", newTrackIndex)
	switch newTrackIndex {
	case 0:
		currentSoundTrack = soundTrack1
	case 1:
		currentSoundTrack = soundTrack2
	case 2:
		currentSoundTrack = soundTrack3
	}

	lastSoundTrackIndex = newTrackIndex
	musicActive = false           // Сбрасываем флаг активности музыки для запуска нового трека
	lastTrackEndTime = time.Now() // Сбрасываем время, чтобы начать отсчет нового трека
}

func updateMusic() {
	now := time.Now()

	switch currentScene {
	case MENU, INVENTORY:
		if lastScene != currentScene {
			// При первом переходе в MENU или INVENTORY ставим музыку на паузу
			if !musicPaused {
				pauseMusic()
			}
			lastScene = currentScene
		}
	case GAME:
		if lastScene == MENU || lastScene == INVENTORY {
			// Возвращаемся из MENU или INVENTORY, возобновляем музыку, не выбирая новый трек
			if musicPaused {
				resumeMusic()
			}
			lastScene = currentScene
		} else {
			// Обрабатываем логику для сцены GAME
			if lastScene != currentScene || now.After(lastTrackEndTime.Add(pauseDuration)) {
				// Если мы только что перешли в GAME или прошло достаточно времени после последнего трека
				selectNewTrackForGame()
				playSoundTrack()
				lastScene = currentScene
				pauseDuration = 5 * time.Minute
			}
		}
	case TITLE:
		if lastScene != currentScene {
			// При переходе в TITLE всегда играем трек для TITLE
			currentSoundTrack = soundTrack3
			playSoundTrack()
			lastScene = currentScene
			pauseDuration = 30 * time.Second
		} else if now.After(lastTrackEndTime.Add(pauseDuration)) {
			// Повторяем трек TITLE каждые 30 секунд
			playSoundTrack()
		}
	}

	rl.UpdateMusicStream(currentSoundTrack) // Обновляем поток музыки
}
