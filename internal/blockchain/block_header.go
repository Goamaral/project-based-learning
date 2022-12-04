package blockchain

import (
	"bytes"
	"crypto/sha256"
	"time"
)

type BlockVersion uint8

const (
	BlockVersion_1 BlockVersion = 1
)

type BlockHeader struct {
	version       BlockVersion
	prevBlockHash Hash
	timestamp     uint32
	difficulty    uint8
	nonce         uint32
}

func newBlockHeader(prevBlockHash Hash, difficulty uint8) *BlockHeader {
	return &BlockHeader{
		version:       BlockVersion_1,
		prevBlockHash: prevBlockHash,
		timestamp:     uint32(time.Now().Unix()),
		difficulty:    difficulty,
	}
}

func (bh BlockHeader) GetHash() Hash {
	return sha256.Sum256(
		bytes.Join(
			[][]byte{
				UintToBytes(uint64(bh.version)),
				bh.prevBlockHash[:],
				UintToBytes(uint64(bh.timestamp)),
				UintToBytes(uint64(bh.difficulty)),
				UintToBytes(uint64(bh.nonce)),
			},
			[]byte{},
		),
	)
}
