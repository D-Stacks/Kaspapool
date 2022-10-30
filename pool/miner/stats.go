package miner

import (
	"math/big"
	"time"
)

type MinerStats struct {
	vardiff 	*big.Float
	FirstSeen	time.Time
	LastSeen	time.Time
	KaspaAddress	string
	WorkerStats	[]*WorkerStats
}

type WorkerStats struct {
	vardiff 	*big.Float
	FirstSeen	time.Time
	LastSeen	time.Time
}
