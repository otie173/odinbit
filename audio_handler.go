package main

import rl "github.com/gen2brain/raylib-go/raylib"

var (
	blockAction rl.Sound
)

func loadAudio() {
	rl.InitAudioDevice()
	blockAction = rl.LoadSound("assets/audio/block_action.mp3")
}

func unloadAudio() {
	rl.UnloadSound(blockAction)
}

func soundBlockAdd() {
	rl.PlaySound(blockAction)
}

func soundBlockRemove() {
	rl.PlaySound(blockAction)
}
