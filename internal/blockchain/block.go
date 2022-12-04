package blockchain

import (
	"blockchain/pb"
	"math/big"

	"google.golang.org/protobuf/proto"
)

type Block struct {
	*BlockHeader
}

func newBlock(prevBlockHash Hash, difficulty uint8) *Block {
	return &Block{BlockHeader: newBlockHeader(prevBlockHash, difficulty)}
}

// A block is valid if the header hash is below the difficulty target
func (b Block) Valid() bool {
	proofHash := b.GetHash()
	var proofHashInt big.Int
	proofHashInt.SetBytes(proofHash[:])

	difficultyTarget := big.NewInt(1)
	difficultyTarget.Lsh(difficultyTarget, uint(b.difficulty))

	return proofHashInt.Cmp(difficultyTarget) == -1
}

func (b Block) Marshal() ([]byte, error) {
	return proto.Marshal(&pb.Block{
		Header: &pb.Block_Header{
			Version:       uint32(b.version),
			PrevBlockHash: b.prevBlockHash[:],
			Timestamp:     b.timestamp,
			Difficulty:    uint32(b.difficulty),
			Nonce:         b.nonce,
		},
	})
}
