package blockchain

import (
	"blockcoin/blockchain/transactions"
	"blockcoin/blockchain/txo"
	"blockcoin/blockchain/wallets"
	"crypto/ecdsa"
	"fmt"
)

type BlockChain interface {
	Mine(txs []*transactions.Transaction)
	ForEach(next func(block *Block) bool)
	FindUnspentTransactions(publicKeyHash []byte) []*transactions.Transaction
	FindUTXO(publicKeyHash []byte) []*transactions.Output
	FindSpendableOutputs(publicKeyHash []byte, amount float64) (float64, *txo.Map)
	FindTransaction(ID []byte) *transactions.Transaction
	SignTransaction(tx *transactions.Transaction, privateKey ecdsa.PrivateKey)
	VerifyTransaction(tx *transactions.Transaction) bool
	Validate() bool
	String() string
}

type blockChain struct {
	blocks []*Block
}

func NewBlockChain(address *wallets.Address) BlockChain {
	return &blockChain{blocks: []*Block{Genesis(address)}}
}

func (bc *blockChain) Mine(txs []*transactions.Transaction) {
	height := len(bc.blocks)
	previousHash := bc.blocks[height-1].PreviousHash
	candidate := &Block{Index: height, PreviousHash: previousHash, Transactions: txs}
	pow := NewProofOfWork(candidate)
	bc.blocks = append(bc.blocks, pow.Mine())
}

func (bc *blockChain) Validate() bool {
	var isValid = true

	bc.ForEach(func(block *Block) bool {

		pow := NewProofOfWork(block)
		isValid = pow.IsValid()
		if !isValid {
			return false
		}

		transactions.ForEach(block.Transactions, func(tx *transactions.Transaction) bool {
			isValid = bc.VerifyTransaction(tx)
			return isValid
		})

		return isValid
	})

	return isValid
}

func (bc *blockChain) ForEach(next func(block *Block) bool) {
	for index := len(bc.blocks) - 1; index >= 0; index-- {
		if !next(bc.blocks[index]) {
			break
		}
	}
}

func Genesis(address *wallets.Address) *Block {
	coinbase := transactions.NewCoinbase(address)
	genesis := &Block{Transactions: []*transactions.Transaction{coinbase}}
	pow := NewProofOfWork(genesis)
	return pow.Mine()
}

func (bc *blockChain) String() string {
	var str string

	for _, block := range bc.blocks {
		str += fmt.Sprint(block)
		str += "\n"
	}

	return str
}
