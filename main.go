package main

import (
	"fmt"

	"github.com/sachinpuranik/blockchain/block"
	"github.com/sachinpuranik/blockchain/wallet"
)

func main() {
	// key, _ := GeneratePrivateKey()
	// fmt.Println(key)

	w := wallet.NewWallet()
	fmt.Println("Public Key 	:", w.PublicKeyStr())
	fmt.Println("Private  Key 	:", w.PrivateKeyStr())
	fmt.Println("Address 		:", w.BlockChainAddress())

	// block0 and chain
	myAddress := "my bc address"
	bc := block.NewBlockChain(myAddress)
	bc.Print()

	// block2
	t := block.NewTransaction("sachin", "suvarna", 20.0)
	bc.AddTransaction(t)
	bc.Mining()
	bc.Print()

	// block3
	t = block.NewTransaction("suvarna", "chetan", 10.0)
	bc.AddTransaction(t)

	t = block.NewTransaction("suvarna", "sachin", 10.0)
	bc.AddTransaction(t)
	bc.Mining()
	bc.Print()
}
