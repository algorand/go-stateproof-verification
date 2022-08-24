package stateproofverification

import (
	"github.com/algorand/go-stateproof-verification/msgpack"
	"github.com/algorand/go-stateproof-verification/stateproofverification/stateproof"
	"github.com/algorand/go-stateproof-verification/transactionverification"
	"github.com/algorand/go-stateproof-verification/types"
)

const strengthTarget = uint64(256)

type StateProofVerifier struct {
	stateProofVerifier *stateproof.Verifier
}

func InitializeVerifier(votersCommitment types.GenericDigest, lnProvenWeight uint64) *StateProofVerifier {
	return &StateProofVerifier{stateProofVerifier: stateproof.MkVerifierWithLnProvenWeight(votersCommitment,
		lnProvenWeight, strengthTarget)}
}

func (v *StateProofVerifier) VerifyStateProofMessage(stateProof *transactionverification.EncodedStateProof, message transactionverification.Message) error {
	messageHash := message.IntoStateProofMessageHash()

	var decodedStateProof stateproof.StateProof
	err := msgpack.Decode(*stateProof, &decodedStateProof)
	if err != nil {
		return err
	}

	return v.stateProofVerifier.Verify(message.LastAttestedRound, messageHash, &decodedStateProof)
}
