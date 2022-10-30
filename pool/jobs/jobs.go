package jobs

import (
	"KPool/block"
	"errors"
	"fmt"

	//	"math/big"

	"github.com/kaspanet/kaspad/app/appmessage"
	"github.com/kaspanet/kaspad/domain/consensus/model/externalapi"
	"github.com/kaspanet/kaspad/domain/consensus/utils/consensushashing"
	"github.com/kaspanet/kaspad/infrastructure/network/rpcclient"

	"sync"
)

type Job struct {
	BlockHash		string
	Timestamp  		int64
	PrePowUInts		[4]uint64
	JobID	   		uint8
}

type Jobs struct {
	Cache			*JobCache
	SeenBlocks		map[string]interface{} //for some reason kaspad send duplicates, ignore
	JobRecivedChan		chan Job
	JobSubmitChan		chan externalapi.DomainBlock
	Index        		uint8
	IsSynced     		bool
	KaspaClient		*rpcclient.RPCClient
	Lock         		sync.RWMutex
}

func NewJobs() (*Jobs, error) {
	kaspaClient, err := rpcclient.NewRPCClient("79.120.76.62:16210")
	if err != nil {
		return nil, err
	}
	return &Jobs{
		Cache:     	   NewJobCache(),
		JobRecivedChan:    make(chan Job),
		Index:		   0,
		IsSynced:	   false,
		KaspaClient:       kaspaClient,
	}, nil

}

func (j *Jobs) Start() {
	go func() { j.recvJobs() } ()
}



func (j *Jobs) onNewBlockTemplate(newBlockTemplate *appmessage.NewBlockTemplateNotificationMessage) {
	j.Lock.Lock()
	defer j.Lock.Unlock()

	rpcBlockTemplateResponse, err := j.KaspaClient.GetBlockTemplate("kaspatest:qprmjxmy36dqpnk6eelndcqta5zxsudg8uf690gn8f20qa364u07xlv63nzfx", "KPool.Developers=[\"jwj\"]")
	if err != nil {
		panic(err)
	}

	j.Index = j.Index + 1

	newBlock, err := block.NewBlockFromRPCBlock(rpcBlockTemplateResponse.Block, j.Index)
	if err != nil {
		panic(err)
	}
	if j.Cache.CheckBlock(newBlock) {
		fmt.Println("FOUND DUPLICATE")
		return
	}

	fmt.Println(consensushashing.BlockHash(newBlock.DomainBlock).String())


	j.Cache.Put(j.Index, newBlock)

	go func() {
		j.JobRecivedChan <- Job{
			BlockHash: consensushashing.BlockHash(newBlock.DomainBlock).String(), 
			JobID: newBlock.JobId, 
			PrePowUInts: newBlock.FourUint64s, 
			Timestamp: newBlock.Timestamp}
		}()

	if j.IsSynced != rpcBlockTemplateResponse.IsSynced {
		j.IsSynced = rpcBlockTemplateResponse.IsSynced
	}
}

func (j *Jobs) SubmitBlock(submitBlock *block.Block) (rejected bool, err error) {
	j.Lock.Lock()
	defer j.Lock.Unlock()

	fmt.Println("submiting block")

	rejectReason, err := j.KaspaClient.SubmitBlock(submitBlock.DomainBlock)
	if err != nil {
		panic(err)
		return false, err
	}
	if rejectReason != 0 {
		panic(errors.New(rejectReason.String()))
		return false, errors.New(rejectReason.String())
	}
	fmt.Println(rejectReason.String())
	fmt.Println("submitted block")
	fmt.Println("block, ", block.GetBlockhash(submitBlock.DomainBlock))

	return true, nil
}

func (j *Jobs) recvJobs() error {

	err := j.KaspaClient.RegisterForNewBlockTemplateNotifications(
		func(notification *appmessage.NewBlockTemplateNotificationMessage) {
		 j.onNewBlockTemplate(notification)
		},
	)
	if err != nil {
		return err
	}
	return nil
}