package block

import (
	"KPool/placeholders"
	"github.com/kaspanet/kaspad/util"
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/kaspanet/kaspad/app/appmessage"
	"github.com/kaspanet/kaspad/domain/consensus/model/externalapi"
	"github.com/kaspanet/kaspad/domain/consensus/utils/consensushashing"
)

var MaxPoW, _ = new(big.Int).SetString("7EEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEEE", 16)

//var Pow1 = new(big.Int).Quo(MaxPoW, new(big.Int).SetUint64(uint64(2 ^ 31)))

type Block struct {
	DomainBlock 		*externalapi.DomainBlock
	PrePowHash 		*externalapi.DomainHash
	FourUint64s		[4]uint64
	blockHash		*externalapi.DomainHash
	Target			*big.Int
	Timestamp		int64
	Confirmations		uint8
	JobId			uint8
}

func NewBlockFromDomainBlock(domainBlock *externalapi.DomainBlock, jobId uint8) (*Block, error) {
	prePowHash := SearializeHeader(domainBlock, true)
	fourUint64s, err := SearializeHeaderToUint64s(domainBlock, true)
	if err != nil {
		return nil, err
	}
	fmt.Println("here")
	for _, out := range domainBlock.Transactions[0].Outputs {
		fmt.Println(out.ScriptPublicKey.Script)
		addr, err := util.NewAddressPublicKey(out.ScriptPublicKey.Script[1:len(out.ScriptPublicKey.Script)-2], placeholders.NetworkPrefix)
		if err != nil {
			panic(err)
		}
		fmt.Println(addr.String())
	}
	
	return &Block{
		DomainBlock: 		domainBlock,
		PrePowHash: 		prePowHash,
		FourUint64s: 		fourUint64s,
		blockHash: 		GetBlockhash(domainBlock),
		Timestamp: 		domainBlock.Header.TimeInMilliseconds(),
		Confirmations:		0,
		JobId: 			jobId,	
	}, nil
}

func NewBlockFromRPCBlock(rpcBlock *appmessage.RPCBlock,  jobId uint8) (*Block, error) {
	domainBlock, err := appmessage.RPCBlockToDomainBlock(rpcBlock)
	if err != nil {
		return nil, err
	}
	return NewBlockFromDomainBlock(domainBlock, jobId)
}

func SearializeHeader(domainBlock *externalapi.DomainBlock, prePow bool) *externalapi.DomainHash {
	if prePow {
		timestamp, nonce := domainBlock.Header.TimeInMilliseconds(), domainBlock.Header.Nonce()
		blockHeader := domainBlock.Header.ToMutable()
		blockHeader.SetTimeInMilliseconds(0)
		blockHeader.SetNonce(0)
		prePowHash := consensushashing.HeaderHash(blockHeader)
		blockHeader.SetTimeInMilliseconds(timestamp)
		blockHeader.SetNonce(nonce)
		domainBlock.Header = blockHeader.ToImmutable()
		return prePowHash
	}
	powHash := consensushashing.HeaderHash(domainBlock.Header)
	return powHash
}

func SearializeHeaderToUint64s(domainBlock *externalapi.DomainBlock, prePow bool) ([4]uint64, error) {
	var uints [4]uint64
	if prePow {
		timestamp, nonce := domainBlock.Header.TimeInMilliseconds(), domainBlock.Header.Nonce()
		header := domainBlock.Header.ToMutable()
		header.SetTimeInMilliseconds(0)
		header.SetNonce(0)
		prePowHashs := consensushashing.HeaderHash(header)
		fmt.Println(prePowHashs.String())
		prePowHash := consensushashing.HeaderHash(header).ByteSlice()
		//for i, j := 0, len(prePowHash)-1; i < j; i, j = i+1, j-1 {
		//	prePowHash[i], prePowHash[j] = prePowHash[j], prePowHash[i]
		//    }
		header.SetTimeInMilliseconds(timestamp)
		header.SetNonce(nonce)
		domainBlock.Header = header.ToImmutable()
		ids := []uint64{}
		ids = append(ids, uint64(binary.LittleEndian.Uint64(prePowHash[0:])))
		ids = append(ids, uint64(binary.LittleEndian.Uint64(prePowHash[8:])))
		ids = append(ids, uint64(binary.LittleEndian.Uint64(prePowHash[16:])))
		ids = append(ids, uint64(binary.LittleEndian.Uint64(prePowHash[24:])))

		final := []uint64{}
		for _, v := range ids {
			asHex := fmt.Sprintf("%x", v)
			bb := big.Int{}
			bb.SetString(asHex, 16)

			final = append(final, bb.Uint64())
		}
		return [4]uint64{final[0], final[1], final[2], final[3]}, nil
	}
	return uints, nil
	}

func GetBlockhash(domainBlock *externalapi.DomainBlock) *externalapi.DomainHash {
	return consensushashing.BlockHash(domainBlock)
}

func (b *Block) Clone() (block *Block) {
	return &Block{
		DomainBlock: 		b.DomainBlock.Clone(),
		PrePowHash:		externalapi.NewDomainHashFromByteArray(b.PrePowHash.ByteArray()),
		blockHash: 		GetBlockhash(b.DomainBlock),
		FourUint64s:		b.FourUint64s,
		Timestamp:		b.Timestamp,
		JobId:			b.JobId,
	}
}