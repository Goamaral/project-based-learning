package blockchain

import (
	"blockchain/pb"
	"fmt"
	"math"
	"math/big"
)

type Block struct {
	BlockHeader
}

func newBlock(prevBlockHash Hash, difficulty uint32) Block {
	return Block{newBlockHeader(prevBlockHash, difficulty)}
}

func (b Block) Proto() *pb.Block {
	return &pb.Block{
		Header: b.BlockHeader.Proto(),
	}
}

// A block is valid if the header hash is below the difficulty target or if is the genesis block
func (b Block) IsValid() (bool, error) {
	if b.IsGenesisBlock() {
		return true, nil
	}

	proofHash, err := b.Hash()
	if err != nil {
		return false, fmt.Errorf("failed to get block header hash")
	}

	var proofHashInt big.Int
	proofHashInt.SetBytes(proofHash[:])

	difficultyTarget := big.NewInt(1)
	difficultyTarget.Lsh(difficultyTarget, uint(math.MaxUint32-b.Difficulty))

	return proofHashInt.Cmp(difficultyTarget) == -1, nil
}

func (b Block) IsGenesisBlock() bool {
	return b.PrevBlockHash == nil
}
