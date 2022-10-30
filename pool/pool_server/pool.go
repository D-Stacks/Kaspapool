package pool_server

import (
	"KPool/extranonce"
	"KPool/jobs"
	"KPool/miner"

	"net"
	"sync"
	"time"
)

type connectionStatus struct {
	Strikes        uint8
	TimedOut       miner.TimeOutEvent
	portsConnected map[string]bool
}

type Pool struct {
	MinersByConnection map[net.Conn]*miner.Miner
	Miners             []*miner.Miner
	ConnectionStatus   map[string]connectionStatus

	Jobs		   *jobs.Jobs

	IncommingConnectionChan    	chan net.Conn
	IncommingDisconnectionChan 	chan net.Conn
	IncommingTimeOutChan       	chan miner.TimeOutEvent
	IncommingCleanUpChan 		chan time.Time
	IncommingJobAddedChan		chan jobs.Job

	Listener net.Listener

	ExtraNonces *extranonce.ExtraNonces

	Lock sync.Mutex
}

func NewPool() (*Pool, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:5550")
	if err != nil {
		return nil, err
	}

	extraNonces := extranonce.NewExtraNonces()

	js, err := jobs.NewJobs()
	if err != nil {
		return nil, err
	}

	return &Pool{
		Miners:             make([]*miner.Miner, 0),
		MinersByConnection: make(map[net.Conn]*miner.Miner),
		ConnectionStatus:   make(map[string]connectionStatus),

		Jobs: js,

		IncommingConnectionChan:    make(chan net.Conn),
		IncommingDisconnectionChan: make(chan net.Conn),
		IncommingTimeOutChan:       make(chan miner.TimeOutEvent),

		IncommingCleanUpChan: 	     make(chan time.Time),
		IncommingJobAddedChan:       js.JobRecivedChan,

		Listener: listener,

		ExtraNonces: extraNonces,
	}, nil
}

//## Run Pool ##

func (p *Pool) RunForever() {
	p.Jobs.Start()
	go func() { p.connectListen() }()
	p.processChannels()
	//go func() { p.tickerEvents() }()
	select {} //halt main process i.e. run forever
}

//## Channel Processing Functions ##

func (p *Pool) processChannels() { // use hierarchical select model to prioritize important pool processing
	var eval bool
	for {
		eval = p.processTier1Priority()
		if eval {
			continue
		}
		eval = p.processTier2Priority()
		if eval {
			continue
		}
		eval = p.processTier3Priority()
		if eval {
			continue
		}
		p.processAnything()
	}
}

func (p *Pool) processTier1Priority() bool {
	select {
	case jobAdded := <- p.IncommingJobAddedChan:
		p.onJobAdded(jobAdded) //prioritze jobs
		return true
	default:
		return false
	}
}

func (p *Pool) processTier2Priority() bool {
	select {
	case timestamp := <- p.IncommingCleanUpChan:
		p.onCleanUp(timestamp)
		return true
	default:
		return false
	}
}

func (p *Pool) processTier3Priority() bool {
	select {
	case disconnect := <- p.IncommingDisconnectionChan:
		p.onDisconnection(disconnect)
		return true 
	case timeOut := <- p.IncommingTimeOutChan:
		p.onTimeOut(timeOut)
		return true 
	default:
		return false
	}
}

func (p *Pool) processAnything() {
	select {
	case jobAdded := <- p.IncommingJobAddedChan:
		p.onJobAdded(jobAdded)
	case timestamp := <- p.IncommingCleanUpChan:
		p.onCleanUp(timestamp)
	case disconnection := <-p.IncommingDisconnectionChan:
		p.onDisconnection(disconnection)
	case timeOut := <- p.IncommingTimeOutChan:
		p.onTimeOut(timeOut)
	case connection := <- p.IncommingConnectionChan:
		p.onConnection(connection)
	}
}