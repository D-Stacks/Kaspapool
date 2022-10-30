package miner

import (
	"KPool/command_client"
	"KPool/jobs"
)



type Channels struct {
	IncommingChans *incommingChans
	OutgoingChans  *outgoingChans
}




type incommingChans struct {
	RecvJob		    	chan jobs.Job
	RecvStratumRequest	chan StratumRequestMessage
	RecvCleanup         	chan CleanUpEvent
	RecvDataGather        	chan command_client.GatherMinerInfo
}

type outgoingChans struct {
	SendTimeOut		chan TimeOutEvent
	SendStats		chan MinerStats
}

func NewChannels(sendTimeOut chan TimeOutEvent) *Channels{
	return &Channels{
		IncommingChans: &incommingChans{
			RecvJob:       		make(chan jobs.Job),
			RecvStratumRequest: 	make(chan StratumRequestMessage),
			RecvCleanup:      	make(chan CleanUpEvent),
		},
		OutgoingChans: &outgoingChans{
			SendTimeOut: sendTimeOut,
		},
	}
}