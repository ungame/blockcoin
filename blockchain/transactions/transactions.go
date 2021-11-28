package transactions

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strings"
)

type Transaction struct {
	ID         []byte
	Inputs     []*Input
	Outputs    []*Output
	isCoinbase bool
}

func (tx *Transaction) Serialize() []byte {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(tx)
	if err != nil {
		log.Panicln(err)
	}
	return buf.Bytes()
}

func (tx *Transaction) IsCoinbase() bool {
	return tx.isCoinbase
}

func (tx *Transaction) Sign(privateKey ecdsa.PrivateKey, prevTXs *Map) {
	if tx.IsCoinbase() {
		return
	}

	for _, input := range tx.Inputs {
		txID := hex.EncodeToString(input.TransactionID)

		if prevTXs.Get(txID) == nil {
			log.Panicf("invalid transaction: ID=%x\n", input.TransactionID)
		}
	}

	txCopy := Copy(tx)

	for index, input := range tx.Inputs {

		txID := hex.EncodeToString(input.TransactionID)

		txCopy.Inputs[index].Signature = nil
		output := prevTXs.Get(txID).Outputs[input.OutputIndex]
		txCopy.Inputs[index].PublicKey = output.PublicKey
		txCopy.ID = NewID(txCopy)

		txCopy.Inputs[index].PublicKey = nil

		r, s, err := ecdsa.Sign(rand.Reader, &privateKey, txCopy.ID)
		if err != nil {
			log.Panicln(err)
		}

		signature := append(r.Bytes(), s.Bytes()...)

		tx.Inputs[index].Signature = signature
	}
}

func (tx *Transaction) Verify(prevTXs *Map) bool {
	if tx.IsCoinbase() {
		return true
	}

	for _, input := range tx.Inputs {
		txID := hex.EncodeToString(input.TransactionID)
		if prevTXs.Get(txID).ID == nil {
			log.Panicf("invalid transaction: ID=%v\n", prevTXs.Get(txID).ID)
		}
	}

	txCopy := Copy(tx)
	curve := elliptic.P256()

	for index, input := range tx.Inputs {
		txID := hex.EncodeToString(input.TransactionID)

		txCopy.Inputs[index].Signature = nil
		output :=  prevTXs.Get(txID).Outputs[input.OutputIndex]
		txCopy.Inputs[index].PublicKey = output.PublicKey
		txCopy.ID = NewID(txCopy)
		txCopy.Inputs[index].PublicKey = nil

		r := big.Int{}
		s := big.Int{}
		length := len(input.Signature)
		r.SetBytes(input.Signature[:(length/2)])
		s.SetBytes(input.Signature[(length/2):])

		x := big.Int{}
		y := big.Int{}
		length = len(input.PublicKey)
		x.SetBytes(input.PublicKey[:(length/2)])
		y.SetBytes(input.PublicKey[(length/2):])

		publicKey := &ecdsa.PublicKey{Curve:curve, X: &x, Y: &y}
		if ecdsa.Verify(publicKey, txCopy.ID, &r, &s) == false {
			return false
		}
	}

	return true
}

func (tx *Transaction) String() string {
	var lines []string

	lines = append(lines, fmt.Sprintf("-- Transaction %x", tx.ID))
	lines = append(lines, fmt.Sprintf("   Coinbase:   %v", tx.IsCoinbase()))

	for index, input := range tx.Inputs {
		lines = append(lines, fmt.Sprintf("    Input %d:", index))
		lines = append(lines, fmt.Sprintf("      TXID:        %x", input.TransactionID))
		lines = append(lines, fmt.Sprintf("      OutputIndex: %d", input.OutputIndex))
		lines = append(lines, fmt.Sprintf("      Signature:   %x", input.Signature))
		lines = append(lines, fmt.Sprintf("      PublicKey:   %x", input.PublicKey))
	}

	for index, output := range tx.Outputs {
		lines = append(lines, fmt.Sprintf("    Output %d:", index))
		lines = append(lines, fmt.Sprintf("      Amount: %.8f", output.Amount))
		lines = append(lines, fmt.Sprintf("      Script: %x", output.PublicKey))
	}

	return strings.Join(lines, "\n")
}

func NewID(tx *Transaction) []byte {
	txCopy := Copy(tx)
	txCopy.ID = []byte{}
	hash := sha256.Sum256(txCopy.Serialize())
	return hash[:]
}

func Copy(tx *Transaction) *Transaction {

	inputs := make([]*Input, 0, len(tx.Inputs))

	for _, input := range tx.Inputs {
		inputs = append(inputs, &Input{
			TransactionID: input.TransactionID,
			OutputIndex:   input.OutputIndex,
			Signature:     input.Signature,
			PublicKey:     input.PublicKey,
		})
	}

	outputs := make([]*Output, 0, len(tx.Outputs))

	for _, output := range tx.Outputs {
		outputs = append(outputs, &Output{
			Amount:    output.Amount,
			PublicKey: output.PublicKey,
		})
	}

	return &Transaction{
		ID:         tx.ID,
		Inputs:     inputs,
		Outputs:    outputs,
		isCoinbase: tx.IsCoinbase(),
	}
}

func HashMulti(txs []*Transaction) []byte {
	var hashes [][]byte
	for _, tx := range txs {
		hashes = append(hashes, tx.ID)
	}
	hash := sha256.Sum256(bytes.Join(hashes, []byte{}))
	return hash[:]
}

func ForEach(txs []*Transaction, next func(tx *Transaction) bool) {
	for _, tx := range txs {
		if !next(tx) {
			break
		}
	}
}
