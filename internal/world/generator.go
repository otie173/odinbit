package world

func generateBarrier() {
	for x := -WorldSize / 2; x <= WorldSize/2; x++ {
		AddBlock(float32(x), float32(-WorldSize/2), Barrier) // Верхняя граница
		AddBlock(float32(x), float32(WorldSize/2), Barrier)  // Нижняя граница
	}
	// Генерация левой и правой границы
	for y := -WorldSize / 2; y <= WorldSize/2; y++ {
		AddBlock(float32(-WorldSize/2), float32(y), Barrier) // Левая граница
		AddBlock(float32(WorldSize/2), float32(y), Barrier)  // Правая граница
	}
}

func generateTree() {

}

func generateStone() {

}

func GenerateWorld() {
	generateBarrier()
}
