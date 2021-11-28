package transactions

import (
	"blockcoin/blockchain/wallets"
	"bytes"
)

type Input struct {
	TransactionID []byte
	OutputIndex   int
	Signature     []byte
	PublicKey     []byte
}

func (in *Input) SignedBy(publicKeyHash []byte) bool {
	signedBy := wallets.NewPublicKeyHash(in.PublicKey)
	return bytes.Compare(signedBy, publicKeyHash) == 0
}
