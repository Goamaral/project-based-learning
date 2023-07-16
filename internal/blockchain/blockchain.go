package blockchain

import (
	"errors"
	"fmt"
	"math"
)

type Blockchain struct {
	difficulty uint8
	blocks     []Block
}

func NewBlockchain(difficulty uint8) Blockchain {
	return Blockchain{
		difficulty: difficulty,
		blocks: []Block{
			newBlock(Hash{}, difficulty), // Genesis block
		},
	}
}

func (bc *Blockchain) MineNewBlock() error {
	println("Creating block")
	block := bc.newBlock()
	println("Block created")

	println("Mining block")
	for {
		if block.nonce == math.MaxInt32 {
			return errors.New("reached max nonce")
		}
		block.nonce++
		if block.Valid() {
			break
		}
	}
	println("Block mined")

	println("Adding block")
	err := bc.addBlock(block)
	if err != nil {
		return fmt.Errorf("failed to add mined block: %w", err)
	}
	println("Block added")

	return nil
}

/* PRIVATE */
func (bc Blockchain) newBlock() Block {
	return newBlock(bc.getLastHash(), bc.difficulty)
}

func (bc *Blockchain) addBlock(block Block) error {
	if !block.Valid() {
		return errors.New("invalid block")
	}

	bc.blocks = append(bc.blocks, block)
	return nil
}

func (bc Blockchain) getLastHash() Hash {
	return bc.blocks[len(bc.blocks)-1].GetHash()
}
