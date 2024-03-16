package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	blockAction  rl.Sound
	pickupAction rl.Sound
)

func loadAudio() {
	rl.InitAudioDevice()
	blockAction = rl.LoadSound("assets/audio/block_action.wav")
	pickupAction = rl.LoadSound("assets/audio/pick_up.wav")
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
