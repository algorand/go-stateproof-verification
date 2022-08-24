package merklearray

import (
	"bytes"
	"errors"
	"fmt"
	"hash"
	"sort"

	"github.com/algorand/go-stateproof-verification/transactionverification"
	"github.com/algorand/go-stateproof-verification/types"
)

const (
	// MaxEncodedTreeDepth is the maximum tree depth (root only depth 0) for a tree which
	// is being encoded (either by msbpack or by the fixed length encoding)
	MaxEncodedTreeDepth = 16

	// MaxNumLeavesOnEncodedTree is the maximum number of leaves allowed for a tree which
	// is being encoded (either by msbpack or by the fixed length encoding)
	MaxNumLeavesOnEncodedTree = 1 << MaxEncodedTreeDepth
)

// Merkle tree errors
var (
	ErrRootMismatch                  = errors.New("root mismatch")
	ErrProofIsNil                    = errors.New("proof should not be nil")
	ErrNonEmptyProofForEmptyElements = errors.New("non-empty proof for empty set of elements")
	ErrPosOutOfBound                 = errors.New("pos out of bound")
)

// Tree is a Merkle tree, represented by layers of nodes (hashes) in the tree
// at each height.
type Tree struct {
	_struct struct{} `codec:",omitempty,omitemptyarray"`

	// Levels represents the tree in layers. layer[0] contains the leaves.
	Levels []Layer `codec:"lvls,allocbound=stateproofcrypto.MaxEncodedTreeDepth+1"`

	// NumOfElements represents the number of the elements in the array which the tree is built on.
	// notice that the number of leaves might be larger in case of a vector commitment
	// In addition, the code will not generate proofs on indexes larger than NumOfElements.
	NumOfElements uint64 `codec:"nl"`

	// Hash represents the hash function which is being used on elements in this tree.
	Hash transactionverification.HashFactory `codec:"hsh"`

	// IsVectorCommitment determines whether the tree was built as a vector commitment
	IsVectorCommitment bool `codec:"vc"`
}

func convertIndexes(elems map[uint64]transactionverification.Hashable, treeDepth uint8) (map[uint64]transactionverification.Hashable, error) {
	msbIndexedElements := make(map[uint64]transactionverification.Hashable, len(elems))
	for i, e := range elems {
		idx, err := merkleTreeToVectorCommitmentIndex(i, treeDepth)
		if err != nil {
			return nil, err
		}
		msbIndexedElements[idx] = e
	}
	return msbIndexedElements, nil
}

func hashLeaves(elems map[uint64]transactionverification.Hashable, treeDepth uint8, hash hash.Hash) (map[uint64]types.GenericDigest, error) {
	hashedLeaves := make(map[uint64]types.GenericDigest, len(elems))
	for i, element := range elems {
		if i >= (1 << treeDepth) {
			return nil, fmt.Errorf("pos %d >= 1^treeDepth %d: %w", i, 1<<treeDepth, ErrPosOutOfBound)
		}
		hashedLeaves[i] = transactionverification.GenericHashObj(hash, element)
	}

	return hashedLeaves, nil
}

func buildFirstPartialLayer(elems map[uint64]types.GenericDigest) partialLayer {
	pl := make(partialLayer, 0, len(elems))
	for pos, elem := range elems {
		pl = append(pl, layerItem{
			pos:  pos,
			hash: elem.ToSlice(),
		})
	}

	sort.Slice(pl, func(i, j int) bool { return pl[i].pos < pl[j].pos })
	return pl
}

func inspectRoot(root types.GenericDigest, pl partialLayer) error {
	computedroot := pl[0]
	if computedroot.pos != 0 || !bytes.Equal(computedroot.hash, root) {
		return ErrRootMismatch
	}
	return nil
}

func verifyPath(root types.GenericDigest, proof *Proof, pl partialLayer) error {
	hints := proof.Path

	s := &siblings{
		hints: hints,
	}

	hsh := proof.HashFactory.NewHash()
	var err error
	for l := uint64(0); len(s.hints) > 0 || len(pl) > 1; l++ {
		if pl, err = pl.up(s, l, true, hsh); err != nil {
			return err
		}
	}

	return inspectRoot(root, pl)
}

// Verify ensures that the positions in elems correspond to the respective hashes
// in a tree with the given root hash.  The proof is expected to be the proof
// returned by Prove().
func Verify(root types.GenericDigest, elems map[uint64]transactionverification.Hashable, proof *Proof) error {
	if proof == nil {
		return ErrProofIsNil
	}

	if len(elems) == 0 {
		if len(proof.Path) != 0 {
			return ErrNonEmptyProofForEmptyElements
		}
		return nil
	}

	hashedLeaves, err := hashLeaves(elems, proof.TreeDepth, proof.HashFactory.NewHash())
	if err != nil {
		return err
	}

	pl := buildFirstPartialLayer(hashedLeaves)
	return verifyPath(root, proof, pl)
}

// VerifyVectorCommitment verifies a vector commitment proof against a given root.
func VerifyVectorCommitment(root types.GenericDigest, elems map[uint64]transactionverification.Hashable, proof *Proof) error {
	if proof == nil {
		return ErrProofIsNil
	}

	msbIndexedElements, err := convertIndexes(elems, proof.TreeDepth)
	if err != nil {
		return err
	}

	return Verify(root, msbIndexedElements, proof)
}