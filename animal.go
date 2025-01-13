package main

import rl "github.com/gen2brain/raylib-go/raylib"

func addAnimal(img rl.Texture2D, x, y float32) {
	addBlock(img, x, y, false)
}

func updateAnimals() {

}

func removeAnimal(x, y float32) {
	removeBlock(x, y)
}
