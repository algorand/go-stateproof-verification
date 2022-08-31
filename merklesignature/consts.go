package merklesignature

import (
	"github.com/algorand/go-stateproof-verification/stateproofcrypto"
)

// HashType/ hashSize relate to the type of hash this package uses.
const (
	MerkleSignatureSchemeRootSize = stateproofcrypto.SumhashDigestSize

	// CryptoPrimitivesID is an identification that the Merkle Signature Scheme uses a subset sum hash function
	// and a falcon signature scheme.
	CryptoPrimitivesID = uint16(0)
)
