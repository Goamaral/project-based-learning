package blockchain

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"math"
	"math/big"
)

const defaultDifficulty = 10

type Blockchain struct {
	blocks      []*Block
	difficulty  int64
	proofTarget *big.Int

	cachedDifficultyBytes []byte
}

func NewBlockchain() *Blockchain {
	proofTarget := big.NewInt(1)
	proofTarget.Lsh(proofTarget, uint(256-defaultDifficulty))

	return &Blockchain{
		blocks:      []*Block{NewBlock("GENESIS", [32]byte{})},
		difficulty:  defaultDifficulty,
		proofTarget: proofTarget,
	}
}

func (bc *Blockchain) AddBlock(block *Block) error {
	// Validate proof
	if !bc.validateNonce(block) {
		return errors.New("nonce generates invalid proof")
	}

	// Add block
	bc.blocks = append(bc.blocks, block)

	return nil
}

func (bc *Blockchain) MineBlock(block *Block) error {
	for block.nonce < math.MaxInt64 {
		if bc.validateNonce(block) {
			break
		} else {
			block.nonce++
		}
	}
	if block.nonce == math.MaxInt64 {
		return errors.New("failed to mine block: reached max nonce")
	}

	return nil
}

func (bc Blockchain) GetLastHash() [32]byte {
	return bc.blocks[len(bc.blocks)-1].hash
}

func (bc *Blockchain) validateNonce(block *Block) bool {
	proofHash := bc.newProofHashCandidate(block)
	return bc.validateProofHash(proofHash)
}

func (bc *Blockchain) newProofHashCandidate(block *Block) [32]byte {
	return sha256.Sum256(bytes.Join(
		[][]byte{
			block.getHeadersBytes(),
			bc.getDifficultyBytes(),
			IntToBytes(block.nonce),
		},
		[]byte{},
	))
}

// Check if hash starts with n zeros (n is the difficulty)
func (bc Blockchain) validateProofHash(proofHash [32]byte) bool {
	var proofHashInt big.Int
	proofHashInt.SetBytes(proofHash[:])
	return proofHashInt.Cmp(bc.proofTarget) == -1
}

func (bc *Blockchain) getDifficultyBytes() []byte {
	if bc.cachedDifficultyBytes == nil {
		bc.cachedDifficultyBytes = IntToBytes(int64(bc.difficulty))
	}
	return bc.cachedDifficultyBytes
}
