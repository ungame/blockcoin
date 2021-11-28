package blockchain

import (
	"blockcoin/blockchain/converter"
	"blockcoin/blockchain/transactions"
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
	"time"
)

const Difficulty = 12

type ProofOfWork struct {
	block  *Block
	target *big.Int
}

func NewProofOfWork(block *Block) *ProofOfWork {
	var pow ProofOfWork
	pow.target = big.NewInt(1)
	pow.target = pow.target.Lsh(pow.target, 256-Difficulty)
	pow.block = block
	return &pow
}

func (pow *ProofOfWork) Hash() []byte {

	payload := bytes.Join([][]byte{
		pow.block.PreviousHash,
		converter.Int64ToBytes(pow.block.Nonce),
		converter.Int64ToBytes(pow.block.Timestamp),
		converter.Int64ToBytes(Difficulty),
		transactions.HashMulti(pow.block.Transactions),
	}, []byte{})

	hash := sha256.Sum256(payload)

	return hash[:]
}

func (pow *ProofOfWork) Mine() *Block {
	var hash []byte

	for i := int64(0); i < math.MaxInt64; i++ {

		pow.block.Nonce = i
		pow.block.Timestamp = time.Now().Unix()

		hash = pow.Hash()
		fmt.Printf("%x\r", hash)

		if pow.isValid(hash) {
			break
		}
	}

	fmt.Println("")
	pow.block.Hash = hash
	return pow.block
}

func (pow *ProofOfWork) isValid(hash []byte) bool {
	var h big.Int
	h.SetBytes(hash)
	return h.Cmp(pow.target) == -1
}

func (pow *ProofOfWork) IsValid() bool {
	return pow.isValid(pow.Hash())
}
