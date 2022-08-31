package test

import (
	_ "embed"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/algorand/go-stateproof-verification/stateproof"
	"github.com/algorand/go-stateproof-verification/stateproofcrypto"
)

//go:embed "resources/previousVotersCommitment.json"
var previousVotersCommitmentData []byte

//go:embed "resources/stateProofMessageHash.json"
var stateProofMessageHashData []byte

//go:embed "resources/stateProof.json"
var stateProofData []byte

func TestVerifier_Verify(t *testing.T) {
	a := require.New(t)

	// Generated from betanet.
	strengthTarget := uint64(256)
	previousLnProvenWeight := uint64(2334949)
	var previousVotersCommitment stateproofcrypto.GenericDigest

	err := json.Unmarshal(previousVotersCommitmentData, &previousVotersCommitment)
	a.NoError(err)

	verifier := stateproof.MkVerifierWithLnProvenWeight(previousVotersCommitment, previousLnProvenWeight, strengthTarget)

	// Generated from betanet.
	lastAttestedRound := uint64(20134400)
	var stateProofMessageHash stateproofcrypto.MessageHash
	var stateProof stateproof.StateProof

	err = json.Unmarshal(stateProofMessageHashData, &stateProofMessageHash)
	a.NoError(err)
	err = json.Unmarshal(stateProofData, &stateProof)
	a.NoError(err)

	err = verifier.Verify(lastAttestedRound, stateProofMessageHash, &stateProof)
	a.NoError(err)
}
