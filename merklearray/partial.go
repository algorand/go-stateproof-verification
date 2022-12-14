package merklearray

import (
	"fmt"
	"hash"

	"github.com/algorand/go-stateproof-verification/stateproofcrypto"
)

// siblings represents the siblings needed to compute the root hash
// given a set of leaf nodes.  This data structure can operate in two
// modes: either build up the set of sibling hints, if tree is not nil,
// or use the set of sibling hints, if tree is nil.
type siblings struct {
	tree  *Tree
	hints []stateproofcrypto.GenericDigest
}

// get returns the sibling from tree level l (0 being the leaves)
// position i.
func (s *siblings) get(l uint64, i uint64) (res stateproofcrypto.GenericDigest, err error) {
	if s.tree == nil {
		if len(s.hints) > 0 {
			res = s.hints[0].ToSlice()
			s.hints = s.hints[1:]
			return
		}

		err = fmt.Errorf("no more sibling hints")
		return
	}

	if l >= uint64(len(s.tree.Levels)) {
		err = fmt.Errorf("level %d beyond tree height %d", l, len(s.tree.Levels))
		return
	}

	if i < uint64(len(s.tree.Levels[l])) {
		res = s.tree.Levels[l][i]
	}

	s.hints = append(s.hints, res)
	return
}

// partialLayer represents a subset of a Layer (i.e., nodes at some
// level in the Merkle tree).  layerItem represents one element in the
// partial Layer.
//msgp:ignore partialLayer
type partialLayer []layerItem

type layerItem struct {
	pos  uint64
	hash stateproofcrypto.GenericDigest
}

// up takes a partial Layer at level l, and returns the next-higher (partial)
// level in the tree.  Since the Layer is partial, up() requires siblings.
//
// The implementation is deterministic to ensure that up() asks for siblings
// in the same order both when generating a proof, as well as when checking
// the proof.
//
// If doHash is false, fill in zero hashes, which suffices for constructing
// a proof.
func (pl partialLayer) up(s *siblings, l uint64, doHash bool, hsh hash.Hash) (partialLayer, error) {
	var res partialLayer
	for i := 0; i < len(pl); i++ {
		item := pl[i]
		pos := item.pos
		posHash := item.hash

		siblingPos := pos ^ 1
		var siblingHash stateproofcrypto.GenericDigest
		if i+1 < len(pl) && pl[i+1].pos == siblingPos {
			// If our sibling is also in the partial Layer, use its
			// hash (and skip over its position).
			siblingHash = pl[i+1].hash
			i++
		} else {
			// Ask for the sibling hash from the tree / proof.
			var err error
			siblingHash, err = s.get(l, siblingPos)
			if err != nil {
				return nil, err
			}
		}

		nextLayerPos := pos / 2
		var nextLayerHash stateproofcrypto.GenericDigest

		if doHash {
			var p pair
			p.hashDigestSize = hsh.Size()
			if pos&1 == 0 {
				// We are left
				p.l = posHash
				p.r = siblingHash
			} else {
				// We are right
				p.l = siblingHash
				p.r = posHash
			}
			nextLayerHash = stateproofcrypto.GenericHashObj(hsh, &p)
		}

		res = append(res, layerItem{
			pos:  nextLayerPos,
			hash: nextLayerHash,
		})
	}

	return res, nil
}
