
State Proof Verification
====================


A Go implementation of functionality required to verify Algorand state proofs.
The functions exported in stateproof.verifier provide the verification interface.

# Install

```bash
go get github.com/algorand/go-stateproof-verification
```
Alternatively the same can be achieved if you use import in a package:

```bash
import "github.com/algorand/go-stateproof-verification"
```
and run go get without parameters.

# Usage

Create a verifier and verify state proof messages using the appropriate state proofs.

```go
package main

import (
	"fmt"
	"github.com/algorand/go-stateproof-verification/stateproof"
	"github.com/algorand/go-stateproof-verification/stateproofcrypto"

	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/types"
)

func main() {
	// verifiedVotersCommitment is the VotersCommitment extracted from the previously verified state proof message.
	var verifiedVotersCommitment stateproofcrypto.GenericDigest
	// verifiedVotersCommitment is the LnProvenWeight extracted from the previously verified state proof message.
	var verifiedLnProvenWeight uint64
	
	// We create a verifier using the aforementioned previously verified data.
	verifier := stateproof.MkVerifierWithLnProvenWeight(verifiedVotersCommitment, verifiedLnProvenWeight)

	// stateProof is the proof used in verification, retrieved from the Algorand blockchain using the API.
	var stateProof stateproof.StateProof
	// stateProofMessage is the message the proof attests to, retrieved from the Algorand blockchain using the API.
	var stateProofMessage types.Message

	// We hash the state proof message using the Algorand SDK. the resulting hash is of the form
	// sha256("spm" || msgpack(stateProofMessage)).
	messageHash := stateproofcrypto.MessageHash(crypto.HashStateProofMessage(&stateProofMessage))

	// We verify the message using the message hash and the state proof.
	err := verifier.Verify(stateProofMessage.LastAttestedRound, messageHash, &stateProof)
	if err != nil {
		fmt.Printf("State proof verification failed: %s\n", err)
	}
}


```

# Testing

```go
go test ./test
```