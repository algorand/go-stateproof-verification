package basics

import (
	"bytes"
	"crypto/sha512"
)

//msgp:allocbound GenericDigest MaxHashDigestSize
type GenericDigest []byte

// ToSlice is used inside the Tree itself when interacting with TreeDigest
func (d GenericDigest) ToSlice() []byte { return d }

// IsEqual compare two digests
func (d GenericDigest) IsEqual(other GenericDigest) bool {
	return bytes.Equal(d, other)
}

// IsEmpty checks wether the generic digest is an empty one or not
func (d GenericDigest) IsEmpty() bool {
	return len(d) == 0
}

// DigestSize is the number of bytes in the preferred hash Digest used here.
const DigestSize = sha512.Size256

// Digest is a SHA512_256 hash
type Digest [DigestSize]byte

// Round represents a round of the Algorand consensus protocol
type Round uint64
