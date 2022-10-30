package miner

import (
	"errors"

	"time"
)

func OnCleanUp(m *Miner) {
	_ = removeIfZombieMiner(m)
}

//Removal 

func removeIfZombieMiner(m *Miner) bool {
	if time.Now().Sub(m.Shares.LastNShares[len(m.Shares.LastNShares) -1].Timestamp) > time.Hour / 2 {
		m.Channels.OutgoingChans.SendTimeOut <- TimeOutEvent{Conn: m.Connection.Conn, End: time.Now().Add(time.Second * 14), Reason: errors.New("Miner not seen in last Hour")}
		return true
	}
	return false
}