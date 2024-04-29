package main

import (
	"embed"
	"log"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
)

//go:embed assets/fonts/pypx/*.ttf
var fonts embed.FS

var (
	bkgColor                  rl.Color = rl.NewColor(0, 0, 0, 255)
	fontBold, font            rl.Font
	screenWidth, screenHeight int32
	monitorFPS, currentFPS    int32
)

func loadFont(fontName string, fontSize int32) rl.Font {
	fontData, err := fonts.ReadFile(fontName)
	if err != nil {
		log.Fatalf("Failed to read embedded font file: %v", err)
	}

	// Создание временного файла для шрифта
	tmpFile, err := os.CreateTemp("", "*.ttf")
	if err != nil {
		log.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // Убедитесь, что временный файл будет удален

	// Запись данных шрифта в временный файл
	_, err = tmpFile.Write(fontData)
	if err != nil {
		log.Fatalf("Failed to write font to temp file: %v", err)
	}
	tmpFile.Close()

	// Загрузка шрифта из временного файла
	font := rl.LoadFontEx(tmpFile.Name(), fontSize, nil, 0)

	return font
}

func update() {
	keyboardHandler()
	mouseHandler()
	updateCameraTarget(&cam, playerPosition, playerRectangle)
	updatePlayerPosition()
	updateMusic()
}

func render() {
	drawScene()
}

func exit() {
	rl.CloseWindow()

	unloadWorld()
	unloadPlayer()
	unloadAudio()
	unloadMusic()
	rl.UnloadFont(fontBold)
	rl.UnloadFont(font)
	unloadInventory()
}

func init() {
	rl.SetConfigFlags(rl.FlagFullscreenMode)
	screenWidth, screenHeight = int32(rl.GetScreenWidth()), int32(rl.GetScreenHeight())
	rl.InitWindow(screenWidth, screenHeight, "Odinbit")
	rl.SetExitKey(0)
	monitorFPS = int32(rl.GetMonitorRefreshRate(rl.GetCurrentMonitor()))
	rl.SetTargetFPS(monitorFPS)
	rl.InitAudioDevice()
	visibleBlocks = make(map[rl.Rectangle]Block)
	prevCamPosition = rl.NewVector2(-1, -1)
	fontBold = loadFont("assets/fonts/pypx/pypx_bold.ttf", 32)
	font = loadFont("assets/fonts/pypx/pypx.ttf", 32)
	loadUI()
	loadWorld()
	loadPlayer()
	loadAudio()
	loadMusic()
	loadInventory()
}

func main() {
	for !rl.WindowShouldClose() {
		update()
		render()
	}
	exit()
}
