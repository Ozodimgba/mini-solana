package poh

import (
	"crypto/sha256"
	"encoding/hex"
)

// State is the internal state of the PoH delay function.
type State [32]byte

// Record "mixes in" a 32-byte value using a new PoH iteration.
func (s *State) Record(mixin *[32]byte) {
	var buf [64]byte
	copy(buf[:32], s[:])
	copy(buf[32:], mixin[:])
	*s = sha256.Sum256(buf[:])
}

// Hash executes a number of PoH iterations.
func (s *State) Hash(n uint) {
	for i := uint(0); i < n; i++ {
		*s = sha256.Sum256(s[:])
	}
}

func (s *State) String() string {
	return hex.EncodeToString(s[:])
}