package miner

import (
	"math/big"
)

func OnAuthorize(m *Miner, authorizeMessage *AuthorizeMessage) {
	m.State.IsAuthorized = true
	m.Workers[authorizeMessage.Worker] = big.NewFloat(0)
	m.Stats.KaspaAddress = authorizeMessage.KaspaAddress.String()
	//m.Connection.Send(StratumAnswerMessage{
	//	ID: authorizeMessage.ID,
	//	Result: true,
	//	Error: nil,
	//})
	//m.Connection.notifyCounter = m.Connection.notifyCounter + 1
}