package main

import (
	"fmt"

	"github.com/sachinpuranik/blockchain/block"
	"github.com/sachinpuranik/blockchain/wallet"
)

func main() {
	walletM := wallet.NewWallet()
	walletA := wallet.NewWallet()
	walletB := wallet.NewWallet()

	walletM.Print()
	walletA.Print()
	walletB.Print()

	t := wallet.NewTransaction(walletA.PrivateKey(), walletA.PublicKey(), walletA.BlockChainAddress(), walletB.BlockChainAddress(), 10)
	fmt.Printf("signature %s", t.GenerateSignature())

	bc := block.NewBlockChain(walletM.BlockChainAddress())
	isAdded := bc.AddTransaction(walletA.BlockChainAddress(), walletB.BlockChainAddress(), 10, walletA.PublicKey(), t.GenerateSignature())

	fmt.Println("Added?", isAdded)

	bc.Mining()
	bc.Print()

	fmt.Println("Total A : ", bc.CalculateTotalAmount(walletA.BlockChainAddress()))
	fmt.Println("Total B : ", bc.CalculateTotalAmount(walletB.BlockChainAddress()))
	fmt.Println("Total M : ", bc.CalculateTotalAmount(walletM.BlockChainAddress()))
}
