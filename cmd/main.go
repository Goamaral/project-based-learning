package main

import (
	"blockchain/internal/blockchain"

	bolt "go.etcd.io/bbolt"
)

func main() {
	db, closeDb, err := openDb()
	if err != nil {
		panic(err)
	}
	defer closeDb()

	bc, err := blockchain.NewBlockchain(db, 20)
	if err != nil {
		panic(err)
	}

	err = bc.MineNewBlock()
	if err != nil {
		panic(err)
	}
}

func openDb() (*bolt.DB, func(), error) {
	db, err := bolt.Open("blockchain.db", 0600, nil)
	if err != nil {
		return nil, nil, err
	}
	return db, func() { db.Close() }, nil
}
