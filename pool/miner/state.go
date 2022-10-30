package miner

type MinerState struct {
	IsSubscribed       bool
	IsAuthorized       bool
	CurrentJobId       uint8
	ExtraNonce         uint16
	MiningAgent	   string
	MiningAgentVersion string
	NotifyCounter      int
}