package wallets

import "blockcoin/blockchain/converter"

type Address struct {
	addr []byte
}

func (a *Address) String() string {
	return string(a.addr)
}

func (a *Address) Bytes() []byte {
	return a.addr
}

func (a *Address) PublicKeyHash() []byte {
	decoded := converter.Base58Decode(a.addr)
	return decoded[1 : len(decoded)-ChecksumLength]
}
