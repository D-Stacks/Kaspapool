package extranonce

import (
	"math/rand"
	"time"
)

const maxUint16 = 65535

type ExtraNonces struct {
	unusedExtraNonces	[]uint16
	usedExtraNonces	   	map[uint16]bool
}

func NewExtraNonces() *ExtraNonces{
	unusedExtraNonces := make([]uint16, maxUint16)
	for i := 0; i < maxUint16; i++ {
		unusedExtraNonces[i] = uint16(i)
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(
		len(unusedExtraNonces), func(i, j int) { 
			unusedExtraNonces[i], unusedExtraNonces[j] = unusedExtraNonces[j], unusedExtraNonces[i] 
		},
	)
	    
	return &ExtraNonces{
		unusedExtraNonces: unusedExtraNonces,
		usedExtraNonces: make(map[uint16]bool),
	}
}

func (en *ExtraNonces) GetExtraNonce() uint16 {
	extraNonce := en.unusedExtraNonces[0]
	en.unusedExtraNonces = en.unusedExtraNonces[1:]
	en.usedExtraNonces[extraNonce] = true
	return extraNonce
}

func (en *ExtraNonces) RemoveExtraNonce(extraNonce uint16) {
	draw := rand.Intn(len(en.unusedExtraNonces)-1)
	en.unusedExtraNonces = append(append(en.unusedExtraNonces[0:draw], extraNonce), en.unusedExtraNonces[draw:len(en.unusedExtraNonces)-1]...)
	delete(en.usedExtraNonces, extraNonce)
}