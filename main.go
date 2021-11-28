package main

import (
	"blockcoin/blockchain"
	"blockcoin/blockchain/transactions"
	"blockcoin/blockchain/wallets"
	"fmt"
)

func main() {
	from := wallets.New()
	to := wallets.New()

	fmt.Printf("new wallet: %s\n", from.Address().String())
	fmt.Printf("new wallet: %s\n", to.Address().String())

	bc := blockchain.NewBlockChain(from.Address())

	tx := blockchain.NewTransaction(from, to.Address(), 25, bc)

	bc.Mine([]*transactions.Transaction{tx})

	tx = blockchain.NewTransaction(to, from.Address(), 5, bc)

	bc.Mine([]*transactions.Transaction{tx})

	fmt.Println(bc)
	if bc.Validate() {
		balanceOf(from, bc)
		balanceOf(to, bc)
	} else {
		fmt.Println("corrupt blockchain.")
	}
}

func balanceOf(wallet *wallets.Wallet, bc blockchain.BlockChain) {
	outpus := bc.FindUTXO(wallet.Address().PublicKeyHash())
	var balance float64

	for _, output := range outpus {
		if output.BelongsTo(wallet.Address().PublicKeyHash()) {
			balance += output.Amount
		}
	}

	fmt.Printf("balance of %s: %.8f\n", wallet.Address().String(), balance)
}
