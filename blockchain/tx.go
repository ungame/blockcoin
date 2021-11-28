package blockchain

import (
	"blockcoin/blockchain/transactions"
	"blockcoin/blockchain/txo"
	"blockcoin/blockchain/wallets"
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"log"
)

func NewTransaction(from *wallets.Wallet, to *wallets.Address, amount float64, chain BlockChain) *transactions.Transaction {
	var inputs []*transactions.Input
	var outputs []*transactions.Output

	publicKeyHash := from.Address().PublicKeyHash()

	accumulated, utxo := chain.FindSpendableOutputs(publicKeyHash, amount)
	if accumulated < amount {
		log.Panicln("not enough funds")
	}

	for id, outputs := range utxo.GetAll() {
		txID, err := hex.DecodeString(id)
		if err != nil {
			log.Panicln("invalid transaction id:", err)
		}

		for _, output := range outputs {
			input := &transactions.Input{
				TransactionID: txID,
				OutputIndex:   output,
				PublicKey:     from.PublicKey,
			}
			inputs = append(inputs, input)
		}
	}

	outputs = append(outputs, transactions.NewOutput(amount, to))

	if accumulated > amount {
		outputs = append(outputs, transactions.NewOutput(accumulated-amount, from.Address()))
	}

	tx := &transactions.Transaction{Inputs: inputs, Outputs: outputs}
	tx.ID = transactions.NewID(tx)

	chain.SignTransaction(tx, from.PrivateKey)

	return tx
}

func (bc *blockChain) SignTransaction(tx *transactions.Transaction, privateKey ecdsa.PrivateKey) {
	prevTXs := transactions.NewMap()

	for _, input := range tx.Inputs {
		prevTX := bc.FindTransaction(input.TransactionID)
		if prevTX == nil {
			log.Panicf("transaction not found: ID=%x\n", input.TransactionID)
		}
		txID := hex.EncodeToString(input.TransactionID)
		prevTXs.Set(txID, prevTX)
	}

	tx.Sign(privateKey, prevTXs)
}

func (bc *blockChain) VerifyTransaction(tx *transactions.Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}

	prevTXs := transactions.NewMap()

	for _, input := range tx.Inputs {
		prevTX := bc.FindTransaction(input.TransactionID)
		if prevTX == nil {
			log.Panicf("invalid transaction: ID=%x\n", input.TransactionID)
		}
		txID := hex.EncodeToString(prevTX.ID)
		prevTXs.Set(txID, prevTX)
	}

	return tx.Verify(prevTXs)
}

func (bc *blockChain) FindUnspentTransactions(publicKeyHash []byte) []*transactions.Transaction {
	var unspentTXs []*transactions.Transaction

	spentTXOs := txo.NewMap()

	bc.ForEach(func(block *Block) bool {

		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			for index, output := range tx.Outputs {
				for _, outputIndex := range spentTXOs.Get(txID) {
					if index == outputIndex {
						continue Outputs
					}
				}
				if output.BelongsTo(publicKeyHash) {
					unspentTXs = append(unspentTXs, tx)
				}
			}

			if !tx.IsCoinbase() {
				for _, input := range tx.Inputs {
					if input.SignedBy(publicKeyHash) {
						inputTXID := hex.EncodeToString(input.TransactionID)
						spentTXOs.Set(inputTXID, input.OutputIndex)
					}
				}
			}
		}

		return true
	})

	return unspentTXs
}

func (bc *blockChain) FindUTXO(publicKeyHash []byte) []*transactions.Output {
	var utxos []*transactions.Output

	unspentTXs := bc.FindUnspentTransactions(publicKeyHash)

	transactions.ForEach(unspentTXs, func(tx *transactions.Transaction) bool {

		for _, output := range tx.Outputs {
			if output.BelongsTo(publicKeyHash) {
				utxos = append(utxos, &transactions.Output{
					Amount:    output.Amount,
					PublicKey: output.PublicKey,
				})
			}
		}

		return true
	})

	return utxos
}

func (bc *blockChain) FindSpendableOutputs(publicKeyHash []byte, amount float64) (float64, *txo.Map) {
	utxo := txo.NewMap()
	unspentTXs := bc.FindUnspentTransactions(publicKeyHash)

	var accumulated float64

Work:
	for _, tx := range unspentTXs {
		txID := hex.EncodeToString(tx.ID)

		for index, output := range tx.Outputs {

			if output.BelongsTo(publicKeyHash) {
				accumulated += output.Amount
				utxo.Set(txID, index)

				if accumulated >= amount {
					break Work
				}
			}
		}
	}

	return accumulated, utxo
}

func (bc *blockChain) FindTransaction(ID []byte) *transactions.Transaction {
	var transaction *transactions.Transaction

	bc.ForEach(func(block *Block) bool {

		transactions.ForEach(block.Transactions, func(tx *transactions.Transaction) bool {
			if bytes.Compare(tx.ID, ID) == 0 {
				transaction = tx
				return false
			}
			return true
		})

		return transaction == nil
	})

	return transaction
}
