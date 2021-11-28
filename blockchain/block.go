package blockchain

import (
	"blockcoin/blockchain/transactions"
	"fmt"
	"strconv"
	"time"
)

type Block struct {
	Index        int
	Hash         []byte
	Nonce        int64
	PreviousHash []byte
	Timestamp    int64
	Transactions []*transactions.Transaction
}

func (b *Block) String() string {
	str := fmt.Sprintf(`
Block         %d
Hash:         %x
Nonce:        %d
PreviousHash: %x
Timestamp:    %s
PoW:          %s
`,
		b.Index,
		b.Hash,
		b.Nonce,
		b.PreviousHash,
		time.Unix(b.Timestamp, 0).String(),
		strconv.FormatBool(NewProofOfWork(b).IsValid()),
	)
	for _, tx := range b.Transactions {
		str += fmt.Sprint(tx.String())
	}
	return str
}
