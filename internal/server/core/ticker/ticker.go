package ticker

import (
	"log"
	"time"
)

type Ticker struct {
	ticker   *time.Ticker
	duration time.Duration
}

func New(tps int) *Ticker {
	tickDuration := time.Second / time.Duration(tps)
	ticker := time.NewTicker(tickDuration)

	return &Ticker{
		ticker:   ticker,
		duration: tickDuration,
	}
}

func (t *Ticker) Run(onTick func()) {
	for range t.ticker.C {
		start := time.Now()
		onTick()
		elapsed := time.Since(start)

		if elapsed > t.duration {
			log.Printf("Warning! Tick took too long: %v\n", elapsed)
		}
		log.Printf("Info! Tick took only: %v\n", elapsed)
	}
}
