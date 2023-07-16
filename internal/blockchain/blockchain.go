package blockchain

import (
	"blockchain/pb"
	"errors"
	"fmt"
	"math"

	bolt "go.etcd.io/bbolt"
	"google.golang.org/protobuf/proto"
)

const InfoBucketName = "info"
const BlockBucketName = "block"

type Blockchain struct {
	db         *bolt.DB
	Difficulty uint32
	LastHash   Hash
	Blocks     map[Hash]Block
}

func NewBlockchain(db *bolt.DB, difficulty uint32) (Blockchain, error) {
	bc := Blockchain{
		db:         db,
		Difficulty: difficulty,
		Blocks:     map[Hash]Block{},
	}
	return bc, db.Update(func(tx *bolt.Tx) error {
		// Reset info bucket
		err := tx.DeleteBucket([]byte(InfoBucketName))
		if err != nil && !errors.Is(err, bolt.ErrBucketNotFound) {
			return fmt.Errorf("failed to delete info bucket: %w", err)
		}
		_, err = tx.CreateBucket([]byte(InfoBucketName))
		if err != nil {
			return fmt.Errorf("failed to create info bucket: %w", err)
		}

		// Reset block bucket
		err = tx.DeleteBucket([]byte(BlockBucketName))
		if err != nil && !errors.Is(err, bolt.ErrBucketNotFound) {
			return fmt.Errorf("failed to delete block bucket: %w", err)
		}
		_, err = tx.CreateBucket([]byte(BlockBucketName))
		if err != nil {
			return fmt.Errorf("failed to create block bucket: %w", err)
		}

		// Set blockchain info
		err = bc.updateInfoBucket(tx)
		if err != nil {
			return fmt.Errorf("failed to update info bucket: %w", err)
		}

		// Add genesis block
		genesisBlock := NewBlock(Hash{}, difficulty)
		_, err = bc.addBlock(tx, genesisBlock)
		if err != nil {
			return fmt.Errorf("failed to add genesis block: %w", err)
		}

		return nil
	})
}

func LoadBlockchain(db *bolt.DB) (Blockchain, error) {
	bc := Blockchain{
		db:     db,
		Blocks: map[Hash]Block{},
	}
	return bc, db.View(func(tx *bolt.Tx) error {
		infoBucket := tx.Bucket([]byte(InfoBucketName))

		// Load config
		bts := infoBucket.Get([]byte("blockchain"))
		pBlockchain := &pb.Blockchain{}
		err := proto.Unmarshal(bts, pBlockchain)
		if err != nil {
			return fmt.Errorf("failed to load blockchain info: %w", err)
		}
		bc.hydrateWithProto(pBlockchain)

		// Load blocks
		blocksBucket := tx.Bucket([]byte(BlockBucketName))
		pBlock := &pb.Block{}
		hash := bc.LastHash
		for {
			_, found := bc.Blocks[hash]
			if found {
				return fmt.Errorf("found duplicate block (%s)", hash)
			}

			bts := blocksBucket.Get(hash[:])
			if bts == nil {
				return fmt.Errorf("failed to find block (%s)", hash)
			}
			err := proto.Unmarshal(bts, pBlock)
			if err != nil {
				return fmt.Errorf("failed to unmarshal block (%s): %w", hash, err)
			}
			block := NewBlockFromProto(pBlock)
			bc.Blocks[hash] = block

			if block.IsGenesisBlock() {
				break
			}
			hash = Hash(block.PrevBlockHash)
		}

		return nil
	})
}

func (bc Blockchain) Proto() *pb.Blockchain {
	return &pb.Blockchain{
		LastHash:   bc.LastHash[:],
		Difficulty: bc.Difficulty,
	}
}

func (bc *Blockchain) MineNewBlock() error {
	block, err := bc.newBlock()
	if err != nil {
		return fmt.Errorf("failed to create block: %w", err)
	}

	println("Mining block")
	for {
		if block.Nonce == math.MaxInt32 {
			return errors.New("reached max nonce")
		}
		block.Nonce++

		isValid, err := block.IsValid()
		if err != nil {
			return fmt.Errorf("failed to validate block: %w", err)
		}
		if isValid {
			break
		}
	}
	println("Block mined")

	err = bc.db.Update(func(tx *bolt.Tx) error {
		_, err := bc.addBlock(tx, block)
		return err
	})
	if err != nil {
		return fmt.Errorf("failed to add mined block: %w", err)
	}

	return nil
}

/* PRIVATE */
func (bc Blockchain) newBlock() (Block, error) {
	return NewBlock(bc.LastHash, bc.Difficulty), nil
}

func (bc *Blockchain) addBlock(tx *bolt.Tx, block Block) (Hash, error) {
	isValid, err := block.IsValid()
	if err != nil {
		return Hash{}, fmt.Errorf("failed to validate block: %w", err)
	}
	if !isValid {
		return Hash{}, errors.New("invalid block")
	}

	hash, err := block.Hash()
	if err != nil {
		return Hash{}, fmt.Errorf("failed to hash block: %w", err)
	}
	bts, err := proto.Marshal(block.Proto())
	if err != nil {
		return Hash{}, fmt.Errorf("failed to marshal block: %w", err)
	}
	bucket := tx.Bucket([]byte(BlockBucketName))
	err = bucket.Put(hash[:], bts)
	if err != nil {
		return Hash{}, fmt.Errorf("failed to put block: %w", err)
	}
	bc.Blocks[hash] = block
	bc.LastHash = hash
	err = bc.updateInfoBucket(tx)
	if err != nil {
		return Hash{}, fmt.Errorf("failed to update info bucket: %w", err)
	}

	return hash, nil
}

func (bc *Blockchain) hydrateWithProto(pBlockchain *pb.Blockchain) {
	bc.LastHash = Hash(pBlockchain.LastHash)
	bc.Difficulty = pBlockchain.Difficulty
}

func (bc *Blockchain) updateInfoBucket(tx *bolt.Tx) error {
	infoBucket := tx.Bucket([]byte(InfoBucketName))
	bts, err := proto.Marshal(bc.Proto())
	if err != nil {
		return fmt.Errorf("failed to marshal blockchain info: %w", err)
	}
	err = infoBucket.Put([]byte("blockchain"), bts)
	if err != nil {
		return fmt.Errorf("failed to put blockchain info: %w", err)
	}
	return nil
}
