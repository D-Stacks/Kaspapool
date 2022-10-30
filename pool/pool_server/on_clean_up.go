package pool_server

import (
	"KPool/miner"
	"sort"
	"time"
)

type SortByVarDiff []*miner.Miner
	
	func (a SortByVarDiff) Len() int           { return len(a) }
	func (a SortByVarDiff) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
	func (a SortByVarDiff) Less(i, j int) bool { return a[i].Shares.CurrentVarDiff.Cmp(a[j].Shares.CurrentVarDiff) == 1 }

func (p *Pool) onCleanUp(timestamp time.Time) {
	var newMiners []*miner.Miner
	sort.Sort(SortByVarDiff(p.Miners))
	for _, m := range p.Miners {
		_, found := p.MinersByConnection[m.Connection.Conn]
		if !found {
			continue
		} else {
			newMiners = append(newMiners, m)
			go func() { m.Channels.IncommingChans.RecvCleanup <- miner.CleanUpEvent{} } ()
		}
	}
	p.Miners = newMiners
}