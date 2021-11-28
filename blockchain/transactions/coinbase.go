package transactions

import (
	"blockcoin/blockchain/wallets"
	"crypto/sha256"
	"fmt"
	"time"
)

const coinbase = 50

func NewCoinbase(address *wallets.Address) *Transaction {

	timestamp := fmt.Sprint(time.Now().Unix())
	signature := sha256.Sum256([]byte(timestamp))
	publicKey := sha256.Sum256([]byte("coinbase"))

	tx := &Transaction{
		Inputs: []*Input{
			{
				TransactionID: []byte{},
				OutputIndex:   -1,
				Signature:     signature[:],
				PublicKey:     publicKey[:],
			},
		},
		Outputs: []*Output{
			NewOutput(coinbase, address),
		},
		isCoinbase: true,
	}

	tx.ID = NewID(tx)

	return tx
}
