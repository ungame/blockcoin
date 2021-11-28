package wallets

import (
	"blockcoin/blockchain/converter"
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"log"
)

const (
	ChecksumLength = 4
	Version        = byte(0x00)
)

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

func (w *Wallet) Address() *Address {
	publicKeyHash := NewPublicKeyHash(w.PublicKey)
	versionedHash := append([]byte{Version}, publicKeyHash...)
	checksum := Checksum(versionedHash)
	fullHash := append(versionedHash, checksum...)
	address := converter.Base58Encode(fullHash)
	return &Address{address}
}

func IsAddress(address *Address) bool {
	decoded := converter.Base58Decode(address.Bytes())
	checksum := decoded[len(decoded)-ChecksumLength:]
	version := decoded[0]
	publicKeyHash := decoded[1 : len(decoded)-ChecksumLength]
	targetChecksum := Checksum(append([]byte{version}, publicKeyHash...))
	return bytes.Compare(checksum, targetChecksum) == 0
}

func New() *Wallet {
	privateKey, publicKey := NewKeyPair()

	return &Wallet{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}
}

func NewKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()

	privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panicln(err)
	}

	publicKey := append(privateKey.X.Bytes(), privateKey.Y.Bytes()...)

	return *privateKey, publicKey
}

func NewPublicKeyHash(publicKey []byte) []byte {
	hash := sha256.Sum256(publicKey)
	return converter.Ripemd160(hash[:])
}

func Checksum(payload []byte) []byte {
	hash256 := sha256.Sum256(payload)
	hash256 = sha256.Sum256(hash256[:])
	return hash256[:ChecksumLength]
}
