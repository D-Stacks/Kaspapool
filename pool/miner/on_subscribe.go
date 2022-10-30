package miner

import (
	"KPool/placeholders"
	"strconv"
)

func OnSubscribe(m *Miner, subscribeMessage *SubscribeMessage) {
	m.State.MiningAgent = subscribeMessage.MiningAgent
	m.State.IsSubscribed = true	

	m.Connection.Send(StratumAnswerMessage{
		ID: subscribeMessage.ID,
		Result: true,
		Error: nil,
	})
	m.Connection.notifyCounter = m.Connection.notifyCounter + 1

	m.Connection.Send(
		StratumRequestMessage{
			ID: m.Connection.notifyCounter,
			Method: "set_extranonce",
			Params: []any{strconv.FormatInt(int64(m.State.ExtraNonce), 16), 6},
		},
	)

	m.Connection.notifyCounter = m.Connection.notifyCounter + 1
	startDiff, _ := placeholders.StartDiff.Float64()
	m.Connection.Send(
		StratumRequestMessage{
			ID: m.Connection.notifyCounter,
			Method: "mining.set_difficulty",
			Params: []any{startDiff * 2},
		},
	)
}