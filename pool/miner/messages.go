package miner

import (
	"KPool/placeholders"
	"encoding/hex"
	"errors"
	"strconv"
	"strings"

	"github.com/kaspanet/kaspad/util"
)

type StratumMessage interface{}

type StratumRequestMessage struct {
	ID     int           	`json:"id"`
	Method string        	`json:"method"`
	Params []any 		`json:"params"`
}

type StratumAnswerMessage struct {
	Result bool		`json:"result"`
	ID     int    		`json:"id"`
	Error  any		`json:"error,omitempty"`
}



type SubmitMessage struct {
	ID		int
	Method		string
	KaspaAddress	util.Address
	Worker		string
	Nonce		uint64
	JobID		uint8
} 


type AuthorizeMessage struct {
	ID		int
	Method		string
	KaspaAddress	util.Address
	Worker		string
} 

type SubscribeMessage struct {
	ID		int
	Method		string
	MiningAgent	string
	Worker		string
}

func ParseMiningSubscribe(stratumRequest StratumRequestMessage) (*SubscribeMessage, error) {
	miningAgent, ok := stratumRequest.Params[0].(string)
	if !ok {
		return nil, errors.New("could not parse mining agent as string")
	}

	return &SubscribeMessage{
		ID: stratumRequest.ID,
		Method: stratumRequest.Method,
		MiningAgent: miningAgent,

	}, nil
}


func ParseMiningAuthorize(m *Miner, stratumRequest StratumRequestMessage) (*AuthorizeMessage, error) {
	minerInfo, ok := stratumRequest.Params[0].(string)
	if !ok {
		return nil, errors.New("could not parse miner Info as string")
	}

	minerInfoSplit := strings.SplitAfter(minerInfo, ".")

	var stringKaspaAddress string
	var worker	string

	if len(minerInfoSplit) > 1{
		stringKaspaAddress, worker = minerInfoSplit[0], minerInfoSplit[1]
	} else {
		stringKaspaAddress = minerInfoSplit[0]
		worker = strconv.FormatInt(int64(len(m.Workers)), 16)
	}

	if len(worker) > 16 {
		return nil, errors.New("name of worker not allowed to be over 16 chars")
	}

	if len(m.Workers) > 2500 {
		return nil, errors.New("exceeding max amount of workers: 2500")
	}

	kaspaAddress, err := util.DecodeAddress(stringKaspaAddress, placeholders.NetworkPrefix)
	if err != nil {
		return nil, err
	}
	return &AuthorizeMessage{
		ID: stratumRequest.ID,
		Method: stratumRequest.Method,
		KaspaAddress: kaspaAddress,
		Worker: worker,
	}, nil
}


func ParseMingingSubmit(m *Miner, stratumRequest StratumRequestMessage) (*SubmitMessage, error){
	minerInfo, ok := stratumRequest.Params[0].(string)
	if !ok {
		return nil, errors.New("could not parse miner Info as string")
	}
	
	jobIDHex, ok := stratumRequest.Params[1].(string)
	if !ok {
		return nil, errors.New("could not parse JobID as string")
	}

	jobID, err := hex.DecodeString(jobIDHex)
	if err != nil {
		return nil, err
	}
	stringNonce, ok := stratumRequest.Params[2].(string)
	if !ok {
		return nil, errors.New("could not parse Nonce as string")
	}
	nonce, err := strconv.ParseUint(stringNonce[2:], 16, 64)
	if err != nil {
		return nil, err
	}
	
	var stringKaspaAddress string
	var worker	string
	
	minerInfoSplit := strings.SplitAfter(minerInfo, ".")
	if len(minerInfoSplit) > 1{
		stringKaspaAddress, worker = minerInfoSplit[0], minerInfoSplit[1]
	} else {
		stringKaspaAddress = minerInfoSplit[0]
		worker = ""

	}

	if len(worker) > 16 {
		return nil, errors.New("name of worker not allowed to be over 16 chars")
	}

	if len(m.Workers) > 2500 {
		return nil, errors.New("exceeding max amount of workers: 2500")
	}


	kaspaAddress, err := util.DecodeAddress(stringKaspaAddress, placeholders.NetworkPrefix)
	if err != nil {
		return nil, err
	}

	return &SubmitMessage{
		ID: 		stratumRequest.ID,
		Method: 	stratumRequest.Method,
		KaspaAddress: 	kaspaAddress,
		Worker: 	worker,
		JobID: 		uint8(jobID[0]),
		Nonce: 		nonce,
	}, nil
}