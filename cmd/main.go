package main

import "blockchain/internal/blockchain"

func main() {
	bc := blockchain.NewBlockchain(20)

	err := bc.MineNewBlock()
	if err != nil {
		panic(err)
	}
}
