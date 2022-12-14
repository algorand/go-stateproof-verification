package stateproof

import (
	"encoding/binary"
	"golang.org/x/crypto/sha3"
	"math/big"

	"github.com/algorand/go-stateproof-verification/stateproofcrypto"
)

const StateProofCoin stateproofcrypto.HashID = "spc"

// The coinChoiceSeed defines the randomness seed that will be given to an XOF function. This will be used for choosing
// the index of the coin to reveal as part of the state proof.
type coinChoiceSeed struct {
	// the ToBeHashed function should be updated when fields are added to this structure
	version        byte
	partCommitment stateproofcrypto.GenericDigest
	lnProvenWeight uint64
	sigCommitment  stateproofcrypto.GenericDigest
	signedWeight   uint64
	data           stateproofcrypto.MessageHash
}

// ToBeHashed returns a binary representation of the coinChoiceSeed structure.
// Since this code is also implemented as a circuit in the stateproof SNARK prover we can't use
// msgpack encoding since it may result in a variable length byte slice.
// Alternatively, we serialize the fields in the structure in a specific format.
func (cc *coinChoiceSeed) ToBeHashed() (stateproofcrypto.HashID, []byte) {
	var signedWtAsBytes [8]byte
	binary.LittleEndian.PutUint64(signedWtAsBytes[:], cc.signedWeight)

	var lnProvenWtAsBytes [8]byte
	binary.LittleEndian.PutUint64(lnProvenWtAsBytes[:], cc.lnProvenWeight)

	coinChoiceBytes := make([]byte, 0, 1+len(cc.partCommitment)+len(lnProvenWtAsBytes)+len(cc.sigCommitment)+len(signedWtAsBytes)+len(cc.data))
	coinChoiceBytes = append(coinChoiceBytes, cc.version)
	coinChoiceBytes = append(coinChoiceBytes, cc.partCommitment...)
	coinChoiceBytes = append(coinChoiceBytes, lnProvenWtAsBytes[:]...)
	coinChoiceBytes = append(coinChoiceBytes, cc.sigCommitment...)
	coinChoiceBytes = append(coinChoiceBytes, signedWtAsBytes[:]...)
	coinChoiceBytes = append(coinChoiceBytes, cc.data[:]...)

	return StateProofCoin, coinChoiceBytes
}

// coinGenerator is used for extracting "randomized" 64 bits for coin flips
type coinGenerator struct {
	shkContext   sha3.ShakeHash
	signedWeight uint64
	threshold    *big.Int
}

// makeCoinGenerator creates a new CoinHash context.
// it is used for squeezing 64 bits for coin flips.
// the function inits the XOF function in the following manner
// Shake(coinChoiceSeed)
// we extract 64 bits from shake for each coin flip and divide it by signedWeight
func makeCoinGenerator(choice *coinChoiceSeed) coinGenerator {
	choice.version = VersionForCoinGenerator
	rep := stateproofcrypto.HashRep(choice)
	shk := sha3.NewShake256()
	//nolint:errcheck // ShakeHash.Write never returns error
	shk.Write(rep)

	threshold := prepareRejectionSamplingThreshold(choice.signedWeight)
	return coinGenerator{shkContext: shk, signedWeight: choice.signedWeight, threshold: threshold}

}

func prepareRejectionSamplingThreshold(signedWeight uint64) *big.Int {
	// we use rejection sampling in order to have a uniform random coin in [0,signedWeight).
	// use b bits (b=64) per attempt.
	// define k = roundDown( 2^b / signedWeight )  implemented as (2^b div signedWeight)
	// and threshold = k*signedWeight
	threshold := &big.Int{}
	threshold.SetUint64(1)

	const numberOfBitsPerAttempt = 64
	threshold.Lsh(threshold, numberOfBitsPerAttempt)

	signedWt := &big.Int{}
	signedWt.SetUint64(signedWeight)

	// k = 2^b / signedWeight
	threshold.Div(threshold, signedWt)

	threshold.Mul(threshold, signedWt)
	return threshold
}

// getNextCoin returns the next 64bits integer which represents a number between [0,signedWeight)
func (cg *coinGenerator) getNextCoin() uint64 {
	// take b bits from the XOF and generate an integer z.
	// we accept the sample if z < threshold
	// else, we reject the sample and repeat the process.
	var randNumFromXof uint64
	for {
		var shakeDigest [8]byte
		//nolint:errcheck // ShakeHash.Read never returns error
		cg.shkContext.Read(shakeDigest[:])
		randNumFromXof = binary.LittleEndian.Uint64(shakeDigest[:])

		z := &big.Int{}
		z.SetUint64(randNumFromXof)
		if z.Cmp(cg.threshold) == -1 {
			break
		}
	}

	return randNumFromXof % cg.signedWeight
}
