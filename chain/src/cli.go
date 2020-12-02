package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"

	"localhost/chain/lib/blockchain"
)

type CLI struct{}

func (cli *CLI) createBlockchain(address string) {
	chain := blockchain.InitBlockChain(address)
	chain.Database.Close()
	fmt.Println("Done!")
}

func (cli *CLI) getBalance(address string) {
	bc := blockchain.ContinueBlockChain(address)

}

func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println(" getbalance -address ADDRESS - get the balance for an address")
	fmt.Println(" createblockchain -address ADDRESS creates a blockchain and sends genesis reward to address")
	fmt.Println(" printchain - Prints the blocks in the chain")
	fmt.Println(" send -from FROM -to TO -amount AMOUNT - Send amount of coins")
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		runtime.Goexit()
	}
}

func (cli *CLI) printChain() {
	chain := blockchain.ContinueBlockChain("")
	defer chain.Database.Close()
	iter := chain.Iterator()

	for {
		block := iter.Next()

		fmt.Printf("Prev. hash: %x\n", block.PrevHash)
		fmt.Printf("Hash: %x\n", block.Hash)
		pow := blockchain.NewProof(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevHash) == 0 {
			break
		}
	}
}
