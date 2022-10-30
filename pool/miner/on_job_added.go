package miner

import (
	//"KPool/defines"
	"KPool/jobs"
	"encoding/hex"
	"fmt"
	"math/big"

	//"math/big"
)

func OnJobAdded(m *Miner, jobAdded jobs.Job) error {
	fmt.Println("got job")
	err := SignalJob(m, jobAdded)
	if err != nil {
		panic(err)
	}

	return nil
}

func SignalJob(m *Miner, jobAdded jobs.Job) error {


	m.Shares.CurrentJob = jobAdded.JobID
	//if m.Shares.GetCurrVardiff().Cmp(big.NewFloat(0)) == 0 {
	//	m.Shares.SetCurrVardiff(big.NewFloat(defines.StartDiff))
	//}
	//m.Shares.SetNextVardiff(m.Shares.GetCurrVardiff())
	
	params := make([]interface{}, 3)

	lastDiff := *big.NewRat(1,1)
	lastDiff.Set(m.Shares.CurrentVarDiff)
	m.Shares.CalibrateCurrentVarDiff()

	if lastDiff.Cmp(m.Shares.CurrentVarDiff) != 0 {
		err := SignalVarDiff(m)
		if err != nil {
			panic(err)
		}
	}

	params[0] = hex.EncodeToString([]byte{jobAdded.JobID,})
	params[1] = jobAdded.PrePowUInts
	params[2] = jobAdded.Timestamp

	m.Connection.notifyCounter = m.Connection.notifyCounter + 1

	err := m.Connection.Send(
		StratumRequestMessage{
			ID: m.Connection.notifyCounter,
			Method: "mining.notify", 
			Params: params,
			},
		)

	if err != nil {
		panic(err)
	}

	m.Shares.VarDiffCache[jobAdded.JobID].Set(m.Shares.CurrentVarDiff)

	return nil
}

func SignalVarDiff(m *Miner) error {

	signaledVarDiff, _ := m.Shares.CurrentVarDiff.Float64()

	fmt.Println(signaledVarDiff)
	
	m.Connection.notifyCounter = m.Connection.notifyCounter + 1
	
	err := m.Connection.Send(
		StratumRequestMessage{
			ID: m.Connection.notifyCounter,
			Method: "mining.set_difficulty",
			Params: []any{signaledVarDiff * 2}, //multiply by two for stratum vardiff
		},
	)
	
	if err != nil {
		return err
	}

	return nil
}
