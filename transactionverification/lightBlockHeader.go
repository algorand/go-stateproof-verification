package transactionverification

import (
	"github.com/algorand/go-stateproof-verification/basics"
	"github.com/algorand/go-stateproof-verification/msgpack"
)

// A Seed contains cryptographic entropy which can be used to determine a
// committee.
type Seed [32]byte

// LightBlockHeader represents a minimal block header. It contains all the necessary fields
// for verifying proofs on transactions.
// In addition, this struct is designed to be used on environments where only SHA256 function exists
type LightBlockHeader struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	Seed                Seed          `codec:"0"`
	RoundNumber         basics.Round  `codec:"r"`
	GenesisHash         basics.Digest `codec:"gh"`
	Sha256TxnCommitment basics.Digest `codec:"tc,allocbound=Sha256Size"`
}

// ToBeHashed implements the crypto.Hashable interface
func (bh LightBlockHeader) ToBeHashed() []byte {
	return append([]byte(BlockHeader256), msgpack.Encode(&bh)...)
}
