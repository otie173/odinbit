package main

import (
	"time"
)

var (
	lastTickTime time.Time = time.Now()
)

func doTick() {
	if time.Since(lastTickTime) >= 1*time.Second {
		lastTickTime = time.Now()
		updateTree()
		updateSeed()
		generateResource()
	}
}
