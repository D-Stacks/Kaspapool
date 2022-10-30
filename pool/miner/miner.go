package miner

import (
	"KPool/block"
	"KPool/jobs"
	"math/big"
	"time"

	"net"

)

type blockSubmitter func(block *block.Block) (rejected bool, err error)

type Miner struct {
	Workers	   		map[string]*big.Float
	
	Stats	 	   	MinerStats
	State		  	MinerState
	Shares			MinerShares
	Channels 		*Channels
	Connection		*MinerConnection

	SubmitBlockFunc		blockSubmitter

	JobCache 		*jobs.JobCache
	
}

func NewMiner(conn net.Conn, extraNonce uint16, jobCache *jobs.JobCache, sendTimeOut chan TimeOutEvent, submitBlockFunc blockSubmitter) *Miner {
	channels := NewChannels(sendTimeOut)
	return &Miner{
		Workers:       make(map[string]*big.Float),
		Stats: MinerStats{
			vardiff: big.NewFloat(0.0),
			FirstSeen: time.Now(),
			LastSeen: time.Now(),
		},
		State: MinerState{
			IsSubscribed: false,
			IsAuthorized: false,
			ExtraNonce: extraNonce,
		},
		Shares: NewMinerShares(),
		Channels: channels,

		Connection: NewMinerConnection(conn, channels),

		SubmitBlockFunc: submitBlockFunc,

		JobCache: jobCache,
	}
}

//## Run Miner ##

func (m *Miner) Run() {
	go func() { m.Connection.RecvLoop() }()
	go func() { m.processIncommingEvents() }()
}

//## Channel Processor ##

func (m *Miner) processIncommingEvents() {
	for {
		eval := m.processTier1Events()
		if eval { continue }
		eval = m.processTier2Events()
		if eval { continue }
		eval = m.processTier3Events()
		if eval { continue }
		m.processAnyEvents()
	}
}

func (m *Miner) processTier1Events() bool {
	select {
	case jobAddedEvent := <-m.Channels.IncommingChans.RecvJob:
		err := OnJobAdded(m, jobAddedEvent)
		if err != nil {
			panic(err)
		}
		return true
	case <- m.Channels.IncommingChans.RecvDataGather:
		OnDataGather(m)
		return true
	default:
		return false
	}
}

func (m *Miner) processTier2Events() bool {
	select {
	case StratumRequestEvent := <-m.Channels.IncommingChans.RecvStratumRequest:
		err := OnStratumRequest(m, StratumRequestEvent)
		if err != nil {
			panic(err)
		}
		return true
	default:
		return false
	}
}

func (m *Miner) processTier3Events() bool {
	select {
	case <-m.Channels.IncommingChans.RecvCleanup:
		OnCleanUp(m)
		return true
	default:
		return false
	}
}

func (m *Miner) processAnyEvents() {
	select {
	case jobAddedEvent := <-m.Channels.IncommingChans.RecvJob:
		err := OnJobAdded(m, jobAddedEvent)
		if err != nil {
			panic(err)
		}
	case StratumRequestEvent := <-m.Channels.IncommingChans.RecvStratumRequest:
		err := OnStratumRequest(m, StratumRequestEvent)
		if err != nil {
			panic(err)
		}
	case <-m.Channels.IncommingChans.RecvCleanup:
		OnCleanUp(m)
	case <- m.Channels.IncommingChans.RecvDataGather:
		OnDataGather(m)
	}
}

func (m *Miner) AddJob(job jobs.Job) {
	go func() { m.Channels.IncommingChans.RecvJob <- job }()
}
