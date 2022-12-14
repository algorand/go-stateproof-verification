package merklesignature

import (
	"encoding/binary"

	"github.com/algorand/go-stateproof-verification/stateproofcrypto"
)

const KeysInMSS stateproofcrypto.HashID = "KP"

type (
	// CommittablePublicKey  is used to create a binary representation of public keys in the merkle
	// signature scheme.
	CommittablePublicKey struct {
		VerifyingKey stateproofcrypto.FalconVerifier
		Round        uint64
	}
)

// ToBeHashed returns the sequence of bytes that would be used as an input for the hash function when creating a merkle tree.
// In order to create a more SNARK-friendly commitment we must avoid using the msgpack infrastructure.
// msgpack creates a compressed representation of the struct which might be varied in length, this will
// be bad for creating SNARK
func (e *CommittablePublicKey) ToBeHashed() (stateproofcrypto.HashID, []byte) {
	verifyingRawKey := e.VerifyingKey.GetFixedLengthHashableRepresentation()

	var roundAsBytes [8]byte
	binary.LittleEndian.PutUint64(roundAsBytes[:], e.Round)

	var schemeAsBytes [2]byte
	binary.LittleEndian.PutUint16(schemeAsBytes[:], CryptoPrimitivesID)

	keyCommitment := make([]byte, 0, len(schemeAsBytes)+len(verifyingRawKey)+len(roundAsBytes))
	keyCommitment = append(keyCommitment, schemeAsBytes[:]...)
	keyCommitment = append(keyCommitment, roundAsBytes[:]...)
	keyCommitment = append(keyCommitment, verifyingRawKey...)

	return KeysInMSS, keyCommitment
}
