package transactions

import (
	"blockcoin/blockchain/wallets"
	"bytes"
)

type Output struct {
	Amount    float64
	PublicKey []byte
}

func NewOutput(amount float64, address *wallets.Address) *Output {
	return &Output{Amount: amount, PublicKey: address.PublicKeyHash()}
}

func (out *Output) BelongsTo(publicKeyHash []byte) bool {
	return bytes.Compare(out.PublicKey, publicKeyHash) == 0
}
