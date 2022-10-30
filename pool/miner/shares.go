package miner

import (
	"KPool/placeholders"
	"KPool/shares"
	"fmt"
	"math/big"
	"time"
)

type MinerShares struct {
	CurrentVarDiff 			*big.Rat
	VarDiffCache 			[256]*big.Rat
	CurrentJob 			uint8
	LastNShares     		[]shares.Share
	ShareBehaviour 			int
}


func NewMinerShares() MinerShares {
	var cache [256]*big.Rat

	for i := 0; i < 256; i++ {
		cache[i] = big.NewRat(1, 1)
	}

	return MinerShares{
		CurrentVarDiff: big.NewRat(1,1).Set(placeholders.StartDiff),
		LastNShares: make([]shares.Share, 0),
		VarDiffCache: cache,
	}
}

func (msh *MinerShares) CalibrateCurrentVarDiff() {
	numOfShares := len(msh.LastNShares)

	currentTime := time.Now()
	
	if numOfShares < 11 {
		if numOfShares < 2 {
			msh.CurrentVarDiff.Set(msh.CurrentVarDiff.Quo(msh.CurrentVarDiff, big.NewRat(5,1)))
		} else if currentTime.Sub(msh.LastNShares[numOfShares -1].Timestamp).Seconds() > placeholders.TargetShareSubmit.Seconds() {
			msh.CurrentVarDiff.Set(msh.CurrentVarDiff.Quo(msh.CurrentVarDiff, big.NewRat(5,1)))
		} else {
			msh.CurrentVarDiff.Set(msh.CurrentVarDiff.Quo(msh.CurrentVarDiff, big.NewRat(1,5)))
		}
		//fmt.Println(msh.CurrentVarDiff)
		return
	}

	cutOff := 0

	//var sharesPerSecond float64

	totalVarDiff := &big.Rat{}
	
	for i, share := range msh.LastNShares[:numOfShares - 1] {
		if currentTime.Sub(share.Timestamp).Seconds() > placeholders.ShareDurationCutOff.Seconds() {
			cutOff = i
			continue
		}
		totalVarDiff.Add(totalVarDiff, share.VarDiff)
	}

	msh.LastNShares = msh.LastNShares[cutOff:]

	numOfShares = len(msh.LastNShares)
	
	//var sharesPerSecond *big.Rat
	calibration := &big.Rat{}
	
	sharesPerSecond := float64(currentTime.Sub(msh.LastNShares[0].Timestamp).Seconds()) / float64(numOfShares)
	fmt.Printf("shares per sec: %f", sharesPerSecond)

	avVardiffPerSecond := totalVarDiff.Quo(totalVarDiff, big.NewRat(int64(numOfShares), 1))
	calibration.SetFloat64(placeholders.TargetShareSubmit.Seconds() / sharesPerSecond)
	msh.CurrentVarDiff.Set(avVardiffPerSecond.Mul(avVardiffPerSecond, calibration))
}
