package blockchain

import (
	"blockchain/pb"
	"crypto/sha256"
	"time"

	"google.golang.org/protobuf/proto"
)

type BlockHeader struct {
	Version       pb.Block_Header_Version
	PrevBlockHash []byte
	Timestamp     uint32
	Difficulty    uint32
	Nonce         uint32
}

func newBlockHeader(prevBlockHash Hash, difficulty uint32) BlockHeader {
	return BlockHeader{
		Version:       pb.Block_Header_VERSION_1,
		PrevBlockHash: prevBlockHash[:],
		Timestamp:     uint32(time.Now().Unix()),
		Difficulty:    difficulty,
	}
}

func (bh *BlockHeader) Hash() (Hash, error) {
	bts, err := proto.Marshal(bh.Proto())
	if err != nil {
		return Hash{}, err
	}
	return Hash(sha256.Sum256(bts)), nil
}

func (bh *BlockHeader) Proto() *pb.Block_Header {
	return &pb.Block_Header{
		Version:       bh.Version,
		PrevBlockHash: bh.PrevBlockHash,
		Timestamp:     bh.Timestamp,
		Difficulty:    bh.Difficulty,
		Nonce:         bh.Nonce,
	}
}
