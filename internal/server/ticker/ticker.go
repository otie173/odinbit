package ticker

import (
	"log"
	"time"
)

type Ticker struct {
	tps    int
	onTick func()
}

func New(tps int, onTick func()) *Ticker {
	return &Ticker{
		tps:    tps,
		onTick: onTick,
	}
}

func (t *Ticker) Run() {
	tickDuration := time.Second / time.Duration(t.tps)
	ticker := time.NewTicker(tickDuration)

	for range ticker.C {
		start := time.Now()
		t.onTick()
		elapsed := time.Since(start)

		if elapsed > tickDuration {
			log.Printf("Warning! Tick took too long: %v\n", elapsed)
		}
		//log.Printf("Info! Tick too only: %v\n", elapsed)
	}
}
