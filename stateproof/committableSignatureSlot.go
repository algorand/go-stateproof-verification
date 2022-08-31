package stateproof

import (
	"encoding/binary"
	"fmt"
	"github.com/algorand/go-stateproof-verification/stateproofcrypto"
)

const StateProofSig stateproofcrypto.HashID = "sps"

type committableSignatureSlot struct {
	sigCommit           sigslotCommit
	serializedSignature []byte
	isEmptySlot         bool
}

func buildCommittableSignature(sigCommit sigslotCommit) (*committableSignatureSlot, error) {
	if sigCommit.Sig.MsgIsZero() { // Empty merkle signature
		return &committableSignatureSlot{isEmptySlot: true}, nil
	}
	if sigCommit.Sig.Signature == nil { // Merkle signature is not empty, but falcon signature is (invalid case)
		return nil, fmt.Errorf("buildCommittableSignature: Falcon signature is nil")
	}
	sigBytes, err := sigCommit.Sig.GetFixedLengthHashableRepresentation()
	if err != nil {
		return nil, err
	}
	return &committableSignatureSlot{sigCommit: sigCommit, serializedSignature: sigBytes, isEmptySlot: false}, nil
}

// ToBeHashed returns the sequence of bytes that would be used as an input for the hash function when creating a merkle tree.
// In order to create a more SNARK-friendly commitment we must avoid using the msgpack infrastructure.
// msgpack creates a compressed representation of the struct which might be varied in length, this will
// be bad for creating SNARK
func (cs *committableSignatureSlot) ToBeHashed() (stateproofcrypto.HashID, []byte) {
	if cs.isEmptySlot {
		return StateProofSig, []byte{}
	}
	var binaryLValue [8]byte
	binary.LittleEndian.PutUint64(binaryLValue[:], cs.sigCommit.L)

	sigSlotByteRepresentation := make([]byte, 0, len(binaryLValue)+len(cs.serializedSignature))
	sigSlotByteRepresentation = append(sigSlotByteRepresentation, binaryLValue[:]...)
	sigSlotByteRepresentation = append(sigSlotByteRepresentation, cs.serializedSignature...)

	return StateProofSig, sigSlotByteRepresentation
}
