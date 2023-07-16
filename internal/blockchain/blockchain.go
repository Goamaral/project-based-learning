package blockchain

import (
	"blockchain/pb"
	"errors"
	"fmt"
	"math"

	bolt "go.etcd.io/bbolt"
	"google.golang.org/protobuf/proto"
)

const BucketName = "blockchain"

type Blockchain struct {
	db     *bolt.DB
	config *pb.Config
	blocks []Block
}

func NewBlockchain(db *bolt.DB, difficulty uint32) (Blockchain, error) {
	bc := Blockchain{
		db:     db,
		config: &pb.Config{Difficulty: difficulty},
	}
	return bc, db.Update(func(tx *bolt.Tx) error {
		// Reset bucket
		err := tx.DeleteBucket([]byte(BucketName))
		if err != nil && !errors.Is(err, bolt.ErrBucketNotFound) {
			return fmt.Errorf("failed to delete reset bucket: %w", err)
		}
		configBucket, err := tx.CreateBucket([]byte(BucketName))
		if err != nil {
			return fmt.Errorf("failed to create config bucket: %w", err)
		}

		// Set config
		bts, err := proto.Marshal(bc.config)
		if err != nil {
			return fmt.Errorf("failed to marshal config: %w", err)
		}
		err = configBucket.Put([]byte("config"), bts)
		if err != nil {
			return fmt.Errorf("failed to put config: %w", err)
		}

		// Add genesis block
		println("Adding genesis block")
		err = bc.addBlock(tx, newBlock(Hash{}, difficulty))
		if err != nil {
			return fmt.Errorf("failed to add genesis block: %w", err)
		}
		println("Genesis block added")

		return nil
	})
}

func (bc *Blockchain) MineNewBlock() error {
	println("Creating block")
	block, err := bc.newBlock()
	if err != nil {
		return fmt.Errorf("failed to create block: %w", err)
	}
	println("Block created")

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

	println("Adding block")
	err = bc.db.Update(func(tx *bolt.Tx) error {
		return bc.addBlock(tx, block)
	})
	if err != nil {
		return fmt.Errorf("failed to add mined block: %w", err)
	}
	println("Block added")

	return nil
}

/* PRIVATE */
func (bc Blockchain) newBlock() (Block, error) {
	hash, err := bc.getLastHash()
	if err != nil {
		return Block{}, fmt.Errorf("failed to get last hash: %w", err)
	}
	return newBlock(hash, bc.config.Difficulty), nil
}

func (bc *Blockchain) addBlock(tx *bolt.Tx, block Block) error {
	isValid, err := block.IsValid()
	if err != nil {
		return fmt.Errorf("failed to validate block: %w", err)
	}
	if !isValid {
		return errors.New("invalid block")
	}

	// Update db
	hash, err := block.Hash()
	if err != nil {
		return fmt.Errorf("failed to hash block: %w", err)
	}
	bts, err := proto.Marshal(block.Proto())
	if err != nil {
		return fmt.Errorf("failed to marshal block: %w", err)
	}
	bucket := tx.Bucket([]byte(BucketName))
	err = bucket.Put(hash[:], bts)
	if err != nil {
		return fmt.Errorf("failed to put block: %w", err)
	}

	bc.blocks = append(bc.blocks, block)

	return nil
}

func (bc Blockchain) getLastHash() (Hash, error) {
	return bc.blocks[len(bc.blocks)-1].Hash()
}
