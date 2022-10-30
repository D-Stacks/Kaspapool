package pool_server

import (
	"time"
)

// send internal pool Ticker Events

func (p *Pool) tickerEvents() {
	for {
		ticker := time.NewTicker(time.Second * 70)
		timestamp := <-ticker.C
		p.IncommingCleanUpChan <- timestamp
	}
}
