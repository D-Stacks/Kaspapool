package miner

import (
	"fmt"
	"time"
)

func OnStratumRequest(m *Miner, stratumRequest StratumRequestMessage) error {
	fmt.Println(stratumRequest)
	switch stratumRequest.Method {
	case "mining.submit":
		submitMessage, err := ParseMingingSubmit(m, stratumRequest)
		if err != nil {
			return err
		}
		OnSubmit(m, submitMessage)
	case "mining.subscribe":
		SubscribeMessage, err := ParseMiningSubscribe(stratumRequest)
		if err != nil {
			return err
		}
		OnSubscribe(m, SubscribeMessage)
	case "mining.authorize":
		AuthorizeMessage, err := ParseMiningAuthorize(m, stratumRequest)
		if err != nil {
			return err
		}
		OnAuthorize(m, AuthorizeMessage)
	default:
		err := fmt.Errorf("unknown Communication Protocol: %s", stratumRequest.Method)
		m.Channels.OutgoingChans.SendTimeOut <- TimeOutEvent{Conn: m.Connection.Conn, End: time.Now().Add(time.Second * 30), Reason: err}
		return err
	}
	
	return nil
}