package miner

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"time"
	"KPool/placeholders"
)


type CleanUpEvent struct{}

type TimeOutEvent struct {
	Conn		net.Conn
	End		time.Time
	Reason		error
}

type MinerConnection struct {
	Conn         		net.Conn
	jsonEncoder  		*json.Encoder
	jsonDecoder  		*json.Decoder
	notifyCounter 		int
	Channels    		*Channels
}


func NewMinerConnection(conn net.Conn, channels *Channels) *MinerConnection {
	return &MinerConnection{
		Conn:        conn,
		jsonEncoder: json.NewEncoder(bufio.NewWriterSize(conn, placeholders.MaxConnectionWriteSizeInByes)),
		jsonDecoder: json.NewDecoder(bufio.NewReaderSize(conn, placeholders.MaxConnectionReadSizeInBytes)),
		Channels: channels,
		notifyCounter: 10,
	}
}


func (mc *MinerConnection) Send(message StratumMessage) error {
	switch message.(type) {
	case StratumRequestMessage:
		fmt.Println(message)
		msg := message.(StratumRequestMessage)
		err := mc.jsonEncoder.Encode(msg)
		if err != nil {
			panic(err)
		}
	case StratumAnswerMessage:
		fmt.Println(message)
		msg := message.(StratumAnswerMessage)
		err := mc.jsonEncoder.Encode(msg)
		if err != nil {
			panic(err)
		}
	}
	return nil
}

func (mc *MinerConnection) RecvLoop() error {
	var err error
	for {
		var stratumRequest StratumRequestMessage
		err = mc.jsonDecoder.Decode(&stratumRequest)
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		go func() { mc.Channels.IncommingChans.RecvStratumRequest <- stratumRequest }()
	}
	mc.Channels.OutgoingChans.SendTimeOut <- TimeOutEvent{Conn: mc.Conn, End: time.Now().Add(time.Second * 3), Reason: err}
	return nil
}
