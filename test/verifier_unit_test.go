package test

import (
	"encoding/json"
	"io/ioutil"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
	
	"github.com/algorand/go-stateproof-verification/stateproof"
	"github.com/algorand/go-stateproof-verification/stateproofcrypto"
)

func readJsonMarshaledFile(filePath string, target interface{}, assertions *require.Assertions) {
	contents, err := ioutil.ReadFile(filePath)
	assertions.NoError(err)

	err = json.Unmarshal(contents, &target)
	assertions.NoError(err)
}

func TestVerifier_Verify(t *testing.T) {
	a := require.New(t)

	// Generated from betanet.
	strengthTarget := uint64(256)
	previousLnProvenWeight := uint64(2334949)
	var previousVotersCommitment stateproofcrypto.GenericDigest

	readJsonMarshaledFile(path.Join("resources", "previousVotersCommitment.json"), &previousVotersCommitment, a)

	verifier := stateproof.MkVerifierWithLnProvenWeight(previousVotersCommitment, previousLnProvenWeight, strengthTarget)

	// Generated from betanet.
	lastAttestedRound := uint64(20134400)
	var stateProofMessageHash stateproofcrypto.MessageHash
	var stateProof stateproof.StateProof

	readJsonMarshaledFile(path.Join("resources", "stateProofMessageHash.json"), &stateProofMessageHash, a)
	readJsonMarshaledFile(path.Join("resources", "stateProof.json"), &stateProof, a)

	err := verifier.Verify(lastAttestedRound, stateProofMessageHash, &stateProof)
	a.NoError(err)
}
