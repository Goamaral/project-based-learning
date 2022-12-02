package main

import "blockchain/pkg/blockchain"

func main() {
	bc := blockchain.NewBlockchain()

	println("Creating block")
	newBlock := blockchain.NewBlock("Goa", bc.GetLastHash())

	println("Mining block")
	err := bc.MineBlock(newBlock)
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
