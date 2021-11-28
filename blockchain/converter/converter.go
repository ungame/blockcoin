package converter

import (
	"bytes"
	"encoding/binary"
	"github.com/mr-tron/base58"
	"golang.org/x/crypto/ripemd160"
	"log"
)

func Int64ToBytes(i int64) []byte {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, i)
	if err != nil {
		log.Panicln(err)
	}
	return buf.Bytes()
}

func Base58Encode(in []byte) []byte {
	encoded := base58.Encode(in)
	return []byte(encoded)
}

func Base58Decode(encoded []byte) []byte {
	decoded, err := base58.Decode(string(encoded))
	if err != nil {
		log.Panicln(err)
	}
	return decoded
}

func Ripemd160(in []byte) []byte {
	hasher := ripemd160.New()
	_, err := hasher.Write(in)
	if err != nil {
		log.Panicln(err)
	}
	return hasher.Sum(nil)
}
