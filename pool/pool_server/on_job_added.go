package pool_server

import (
	"KPool/jobs"
)

func (p *Pool) onJobAdded(jobAdded jobs.Job) {
	for i, _ := range p.Miners {
		if p.Miners[i].State.IsAuthorized && p.Miners[i].State.IsSubscribed {
			p.Miners[i].AddJob(jobAdded)
	}
	}
}
