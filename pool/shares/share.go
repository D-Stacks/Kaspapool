package shares

import (
	"encoding/binary"
	"errors"
	"math/big"
	"time"
	"bytes"
	"encoding/gob"

	"github.com/kaspanet/kaspad/domain/consensus/model/externalapi"
	"github.com/kaspanet/kaspad/util"
)

type Share struct {
	WorkerName		string
	PayoutAddress 		util.Address
	BlockHash		externalapi.DomainHash
	Nonce			uint64
	JobId			uint8
	VarDiff			*big.Rat
	ShareDiff		*big.Rat
	KaspaDiff		*big.Rat
	Timestamp		time.Time

	IsValid			bool

	Paid			bool
	ValueInUSD		bool
	Reward			float64

}

func NewShare(kaspaAddress util.Address, workerName string, 
	varDiff *big.Rat, shareDiff *big.Rat, kaspaDiff *big.Rat, blockHash externalapi.DomainHash, 
	jobID uint8, Nonce uint64, miningAgent string) (*Share, error) {

	if len(workerName) > 16 {
		return nil, errors.New("worker name is not allowed to be over 16 bytes long")
	}

	return &Share{
		WorkerName: 	workerName,
		PayoutAddress: 	kaspaAddress,
		Timestamp: 	time.Now(),
		BlockHash: 	blockHash,
		VarDiff:   	varDiff,
		ShareDiff:	shareDiff,
		KaspaDiff:   	kaspaDiff,
	}, nil
}

func (s *Share) Key() []byte { //blockhash + Nonce should ensure share is unique for checks.
	nonceBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(nonceBytes, s.Nonce)
	return append(s.BlockHash.ByteSlice(), nonceBytes...)
}

func (s *Share) Searilize() ([]byte, error) {
	res := new(bytes.Buffer)
	enc := gob.NewEncoder(res)
	err := enc.Encode(s)
	if err != nil {
		return nil, err
	}
	return res.Bytes(), nil
}

func Deseailize(shareBytes []byte) (*Share, error) {
	var share Share
	dec := gob.NewDecoder(bytes.NewReader(shareBytes))
	err := dec.Decode(&share)
	if err != nil {
		return nil, err
	}
	return &share, nil
}