package blockchain

import (
	"blockchain/pb"
	"crypto/sha256"
	"time"

	"google.golang.org/protobuf/proto"
)

type Block_Header struct {
	Version       pb.Block_Header_Version
	PrevBlockHash Hash
	Timestamp     uint32
	Difficulty    uint32
	Nonce         uint32
}

func NewBlock_Header(prevBlockHash Hash, difficulty uint32) Block_Header {
	return Block_Header{
		Version:       pb.Block_Header_VERSION_1,
		PrevBlockHash: prevBlockHash,
		Timestamp:     uint32(time.Now().Unix()),
		Difficulty:    difficulty,
	}
}

func NewBlock_HeaderFromProto(pBlockheader *pb.Block_Header) Block_Header {
	return Block_Header{
		Version:       pBlockheader.Version,
		PrevBlockHash: Hash(pBlockheader.PrevBlockHash),
		Timestamp:     pBlockheader.Timestamp,
		Difficulty:    pBlockheader.Difficulty,
		Nonce:         pBlockheader.Nonce,
	}
}

func (bh *Block_Header) Proto() *pb.Block_Header {
	return &pb.Block_Header{
		Version:       bh.Version,
		PrevBlockHash: bh.PrevBlockHash[:],
		Timestamp:     bh.Timestamp,
		Difficulty:    bh.Difficulty,
		Nonce:         bh.Nonce,
	}
}

func (bh *Block_Header) Hash() (Hash, error) {
	bts, err := proto.Marshal(bh.Proto())
	if err != nil {
		return Hash{}, err
	}
	return Hash(sha256.Sum256(bts)), nil
}
