package miner

import (
	"KPool/block"
	"KPool/placeholders"
	"KPool/shares"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/kaspanet/kaspad/domain/consensus/utils/consensushashing"
	"github.com/kaspanet/kaspad/domain/consensus/utils/pow"
	"github.com/kaspanet/kaspad/util/difficulty"
)

func OnSubmit(m *Miner, submitMessage *SubmitMessage) error {
	blockIfValid, shareIfValid, err := evaluateSubmitMessage(m, submitMessage)
	if err != nil {
		m.Connection.Send(
			StratumAnswerMessage{
				ID:     submitMessage.ID,
				Result: false,
				Error:  err,
			})
	}
	if blockIfValid != nil {
		_, err := m.SubmitBlockFunc(blockIfValid)
		if err != nil {
			m.Connection.Send(
				StratumAnswerMessage{
					ID:     submitMessage.ID,
					Result: false,
					Error:  err,
				})
		}
	}
	if shareIfValid != nil {
		err = m.Connection.Send(
			StratumAnswerMessage{
				ID:     submitMessage.ID,
				Result: true,
				Error:  nil,
			})
		if err != nil {
			return err
		}
		logValidShare(m, shareIfValid)
		pruneShares(m)
	} else {
		logInvalidShare(m)
	}

	return nil

}

func evaluateSubmitMessage(m *Miner, submitMessage *SubmitMessage) (blockIfValid *block.Block, shareIfValid *shares.Share, err error) {

	jobBlock, found := m.JobCache.Get(submitMessage.JobID)
	if !found { //stale job
		return nil, nil, errors.New("submitted job for stale block")
	}

	state := pow.NewState(jobBlock.DomainBlock.Header.ToMutable())
	state.Nonce = submitMessage.Nonce
	shareDiff := new(big.Rat).SetFrac(block.MaxPoW, state.CalculateProofOfWorkValue())
	KaspaDiff := new(big.Rat).SetFrac(block.MaxPoW, difficulty.CompactToBig(jobBlock.DomainBlock.Header.Bits()))
	if shareDiff.Cmp(KaspaDiff) > -1 {
		fmt.Println()
		fmt.Println("PASSED CHECK")
		fmt.Println()
		fmt.Println("pre: ",block.GetBlockhash(jobBlock.DomainBlock))
		mutHeader :=  jobBlock.DomainBlock.Header.ToMutable()
		mutHeader.SetNonce(state.Nonce)
		jobBlock.DomainBlock.Header = mutHeader.ToImmutable()
		fmt.Println("post: ",block.GetBlockhash(jobBlock.DomainBlock))
		fmt.Println()
		shareIfValid, err = shares.NewShare(
			submitMessage.KaspaAddress,
			submitMessage.Worker,
			m.Shares.VarDiffCache[submitMessage.JobID],
			shareDiff,
			KaspaDiff,
			*consensushashing.BlockHash(jobBlock.DomainBlock),
			submitMessage.JobID,
			submitMessage.Nonce,
			m.State.MiningAgent,
		)
		if err != nil {
			return nil, nil, err
		}

		return jobBlock, shareIfValid, nil

	//work := difficulty.CalcWork(jobBlock.DomainBlock.Header.Bits())

	} else if shareDiff.Cmp(m.Shares.VarDiffCache[submitMessage.JobID]) > -1 {
		gotten, _ := shareDiff.Float64()
		needed, _ := m.Shares.VarDiffCache[submitMessage.JobID].Float64()
		fmt.Println("%s vs %s", gotten, needed)
		shareIfValid, err = shares.NewShare(
			submitMessage.KaspaAddress,
			submitMessage.Worker,
			m.Shares.VarDiffCache[submitMessage.JobID],
			shareDiff,
			KaspaDiff,
			*consensushashing.BlockHash(jobBlock.DomainBlock),
			submitMessage.JobID,
			submitMessage.Nonce,
			m.State.MiningAgent,
		)
		if err != nil {
			return nil, nil, err
		}

		return nil, shareIfValid, nil
	}

	gotten, _ := shareDiff.Float64()
	needed, _ := m.Shares.VarDiffCache[submitMessage.JobID].Float64()

	return nil, nil, fmt.Errorf("share with JobID %d had diff of %s did not satify vardiff of of %s", submitMessage.JobID, gotten, needed)
}

func logValidShare(m *Miner, share *shares.Share) {
	if m.Shares.ShareBehaviour < 99 {
		m.Shares.ShareBehaviour = m.Shares.ShareBehaviour + 1
	}

	if len(m.Shares.LastNShares) > 99 {
		m.Shares.LastNShares = append(m.Shares.LastNShares[1:], *share)
	} else {
		m.Shares.LastNShares = append(m.Shares.LastNShares, *share)
	}

}

func logInvalidShare(m *Miner) {
	m.Shares.ShareBehaviour = m.Shares.ShareBehaviour - 1
	if m.Shares.ShareBehaviour < -100 {
		m.Channels.OutgoingChans.SendTimeOut <- TimeOutEvent{
			Conn:   m.Connection.Conn,
			End:    time.Now().Add(time.Second * 30),
			Reason: errors.New("supplied too many bad shares"),
		}
	}
}

func pruneShares(m *Miner) {
	var cutOff int
	currentTime := time.Now()

	for i, share := range m.Shares.LastNShares {
		if currentTime.Sub(share.Timestamp).Nanoseconds() < placeholders.ShareDurationCutOff.Nanoseconds() {
			cutOff = i
			break
		}
	}
	m.Shares.LastNShares = m.Shares.LastNShares[cutOff:]
}
