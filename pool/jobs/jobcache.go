package jobs

import (
	"KPool/block"
	"sync/atomic"
)

//optimized atomic job cache: for mostly read int-key, same value.

const cacheSize = uint8(255)

type cache map[uint8]*block.Block


type JobCache struct {
	atmoizer 	atomic.Value
	filter 		map[string]bool //filter duplicates
	//currentWorkID	uint8
}

func NewJobCache() *JobCache {
	var cache cache
	var atomizer atomic.Value
	atomizer.Store(cache)
	self :=  &JobCache{
		atmoizer: atomizer,
		filter: make(map[string]bool),
	}
	
	return self
}

func (jc *JobCache) Get(workID uint8) (*block.Block, bool) {
	block, found := jc.atmoizer.Load().(cache)[workID]
	//block, found := cache[workID]
	if !found {
		return nil, found
	}
	return block.Clone(), found
}

func (jc *JobCache) Put(workID uint8, insertBlock *block.Block) {
	jc.filter[block.GetBlockhash(insertBlock.DomainBlock).String()] = true
	oldCache, _ := jc.atmoizer.Load().(cache)
	oldBlock, found := oldCache[workID]
	if found {
		delete(jc.filter, block.GetBlockhash(oldBlock.DomainBlock).String())
	}
	newCache := make(cache)
	for i, block := range oldCache{
		newCache[i] = block
	}
	newCache[workID] = insertBlock

	jc.atmoizer.Store(newCache)
}

func (jc *JobCache) CheckBlock(checkBlock *block.Block) bool {
	_, found := jc.filter[block.GetBlockhash(checkBlock.DomainBlock).String()]
	return found
}
