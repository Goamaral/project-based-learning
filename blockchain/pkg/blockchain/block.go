package blockchain

import (
	"bytes"
	"crypto/sha256"
	"time"
)

type Block struct {
	data          string
	timestamp     int64
	prevBlockHash [32]byte
	hash          [32]byte
	nonce         int64

	cachedTimestampBytes []byte
	cachedHeadersBytes   []byte
}

func NewBlock(data string, prevBlockHash [32]byte) *Block {
	block := &Block{data: data, timestamp: time.Now().Unix(), prevBlockHash: prevBlockHash}
	block.UpdateHash()
	return block
}

func (b *Block) UpdateHash() {
	if b.cachedTimestampBytes == nil {
		b.cachedTimestampBytes = IntToBytes(b.timestamp)
	}

	b.hash = sha256.Sum256(b.getHeadersBytes())
}

func (b *Block) getHeadersBytes() []byte {
	if b.cachedHeadersBytes == nil {
		b.cachedHeadersBytes = bytes.Join([][]byte{b.prevBlockHash[:], []byte(b.data), b.cachedTimestampBytes}, []byte{})
	}
	return b.cachedHeadersBytes
}
