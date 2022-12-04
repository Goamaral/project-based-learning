package blockchain

import (
	"errors"
	"math"
)

type Blockchain struct {
	difficulty uint8
	blocks     []*Block
}

func NewBlockchain(difficulty uint8) *Blockchain {
	return &Blockchain{
		difficulty: difficulty,
		blocks:     []*Block{},
	}
}

func (bc *Blockchain) AddGenisisBlock() error {
	if len(bc.blocks) != 0 {
		return errors.New("Already has blocks")
	}

	bc.blocks = append(bc.blocks, newBlock(Hash{}, bc.difficulty))
	return nil
}

func (bc Blockchain) NewBlock() *Block {
	return newBlock(bc.GetLastHash(), bc.difficulty)
}

func (bc *Blockchain) AddBlock(block *Block) error {
	if !block.Valid() {
		return errors.New("invalid block")
	}

	bc.blocks = append(bc.blocks, block)
	return nil
}

func (bc *Blockchain) MineBlock(block *Block) error {
	for {
		if block.nonce == math.MaxInt32 {
			return errors.New("reached max nonce")
		}
		block.nonce++
		if block.Valid() {
			break
		}
	}

	return nil
}

func (bc Blockchain) GetLastHash() Hash {
	return bc.blocks[len(bc.blocks)-1].GetHash()
}
