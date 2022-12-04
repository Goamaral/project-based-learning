package main

import "blockchain/internal/blockchain"

func main() {
	bc := blockchain.NewBlockchain(250)
	err := bc.AddGenisisBlock()
	if err != nil {
		panic(err)
	}

	println("Creating block")
	newBlock := bc.NewBlock()

	println("Mining block")
	err = bc.MineBlock(newBlock)
	if err != nil {
		panic(err)
	}

	println("Adding block")
	err = bc.AddBlock(newBlock)
	if err != nil {
		panic(err)
	}

	println("Block added")
}
