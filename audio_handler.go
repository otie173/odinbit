package main

import (
	"embed"
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
	pickupAction        rl.Sound
	plantAction         rl.Sound
	musicPaused         bool
	musicActive         bool
	currentSoundTrack   rl.Music
	previousSoundTrack  rl.Music
	lastSoundTrackIndex int = -1
	soundTrack1         rl.Music
	soundTrack2         rl.Music
	soundTrack3         rl.Music
	soundTrack4         rl.Music
	trackEndTime        time.Time
)

const (
	SOUNDTRACK_VOLUME float32 = 0.20
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
	pickupAction = loadSound("assets/audio/pick_up.ogg")
	plantAction = loadSound("assets/audio/plant_seed.ogg")
}

func unloadAudio() {
	rl.UnloadSound(pickupAction)
	rl.UnloadSound(plantAction)
}

func loadMusic() {
	soundTrack1 = loadSoundTrack("assets/music/001_OdeToLoneliness_v02.ogg")
	soundTrack1.Looping = false
	soundTrack2 = loadSoundTrack("assets/music/002_InsatiableCuriosity_v02.ogg")
	soundTrack2.Looping = false
	soundTrack3 = loadSoundTrack("assets/music/003_TreesandRocks_v01.ogg")
	soundTrack3.Looping = false
	soundTrack4 = soundTrack3
}

func unloadMusic() {
	rl.StopMusicStream(currentSoundTrack)
	rl.UnloadMusicStream(soundTrack3)
}

func pickupResourceSound() {
	rl.PlaySound(pickupAction)
}

func plantSeedSound() {
	rl.PlaySound(plantAction)
}

func playSoundTrack() {
	// Это сделано для того чтобы, если в сцене GAME играл soundTrack3, то
	// при выходе в сцену TITLE - soundTrack3 мог звучать снова
	if !musicActive || currentSoundTrack == soundTrack4 && currentScene == TITLE {
		rl.StopMusicStream(previousSoundTrack)
		rl.SetMusicVolume(currentSoundTrack, SOUNDTRACK_VOLUME)
		rl.PlayMusicStream(currentSoundTrack)
		musicPaused = false

	}

	if !musicActive || currentSoundTrack != previousSoundTrack {
		rl.StopMusicStream(previousSoundTrack)
		previousSoundTrack = currentSoundTrack

		rl.SetMusicVolume(currentSoundTrack, SOUNDTRACK_VOLUME)
		rl.PlayMusicStream(currentSoundTrack)
		musicPaused = false
	}
}

func pauseMusic() {
	if !musicPaused {
		rl.PauseMusicStream(currentSoundTrack)
		musicPaused = true
	}
}

func resumeMusic() {
	if musicPaused {
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

	switch newTrackIndex {
	case 0:
		currentSoundTrack = soundTrack1
	case 1:
		currentSoundTrack = soundTrack2
	case 2:
		currentSoundTrack = soundTrack3
	}

	lastSoundTrackIndex = newTrackIndex
	musicActive = false // Сбрасываем флаг активности музыки для запуска нового трека

	// Важно обновить поток музыки перед получением его длительности
	rl.UpdateMusicStream(currentSoundTrack)
	// Получаем длительность трека в секундах
	trackLength := rl.GetMusicTimeLength(currentSoundTrack)
	// Устанавливаем trackEndTime добавляя к текущему времени длительность трека и дополнительные 3 минуты и 30 секунд
	trackEndTime = time.Now().Add(time.Second*time.Duration(trackLength) + (3*time.Minute + 30*time.Second))

}

func updateMusic() {
	now := time.Now()

	if !rl.IsMusicStreamPlaying(currentSoundTrack) && now.After(trackEndTime) && !musicPaused {
		selectNewTrackForGame()
		playSoundTrack()
	} else {
		// Обновляем текущий поток музыки
		rl.UpdateMusicStream(currentSoundTrack)
	}

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
			if lastScene != currentScene {
				// Если мы только что перешли в GAME
				rl.StopMusicStream(currentSoundTrack)
				lastScene = currentScene
			}
		}
	case TITLE:
		if lastScene != currentScene {
			// При переходе в TITLE всегда играем трек для TITLE
			currentSoundTrack = soundTrack4
			playSoundTrack()
			lastScene = currentScene
		}
	}
}
