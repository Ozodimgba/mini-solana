package poh

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

// Constants for PoH and Tower BFT
const (
	BufferSize      = 10000                  // Define the size of the circular buffer
	SlotDuration    = 400 * time.Millisecond // Duration of each slot in milliseconds
	MaxSlots        = 512                      // Maximum number of slots to handle
	MaxTimeoutSlots = MaxSlots               // Align max timeouts with the number of slots
)

// PoH represents a Proof of History generator and Tower BFT implementation.
type PoH struct {
	LastHash     string
	Sequence     int
	Hashes       [BufferSize]string
	CurrentIndex int
	Slot         int
	Timeouts     []int // Timeout for each slot in number of slots
}

// NewPoH creates a new PoH generator with an initial seed and sets up the circular buffer.
func NewPoH(seed string) *PoH {
	poh := &PoH{
		LastHash:     seed,
		Sequence:     0,
		CurrentIndex: 0,
		Slot:         0,
		Timeouts:     make([]int, MaxSlots),
	}

	// Initialize the buffer with the seed
	poh.Hashes[0] = seed

	return poh
}

// GenerateNextHash generates the next hash in the PoH sequence as fast as possible.
func (p *PoH) GenerateNextHash() {
	hash := sha256.New()
	hash.Write([]byte(p.LastHash))
	hashBytes := hash.Sum(nil)
	newHash := hex.EncodeToString(hashBytes)

	// Update PoH state
	p.LastHash = newHash
	p.Sequence++

	// Store in circular buffer, overwriting the oldest hash
	p.CurrentIndex = (p.CurrentIndex + 1) % BufferSize
	p.Hashes[p.CurrentIndex] = newHash

	// Log the hash and handle slot transitions
	if p.Sequence%(BufferSize/MaxSlots) == 0 {
		p.Slot = (p.Slot + 1) % MaxSlots
		p.updateTimeouts()
		fmt.Printf("Slot %d: %s (timeouts: %v)\n", p.Slot, newHash, p.Timeouts)
	} else {
		fmt.Printf("Tower BFT %d: %s\n", p.Sequence, newHash)
	}
}

// updateTimeouts updates the timeouts for Tower BFT.
func (p *PoH) updateTimeouts() {
	// Double the timeouts for votes in previous slots
	for i := range p.Timeouts {
		if p.Timeouts[i] > 0 {
			p.Timeouts[i] *= 2
		}
	}

	// Add a new vote with a base timeout (2 slots initially)
	if len(p.Timeouts) < MaxTimeoutSlots {
		p.Timeouts = append([]int{2}, p.Timeouts...)
	} else {
		p.Timeouts[0] = 2
	}
}

// VerifyHash verifies that a given hash was correctly generated as part of the PoH sequence.
func (p *PoH) VerifyHash(sequence int) bool {
	if sequence <= 0 || sequence >= len(p.Hashes) {
		return false
	}

	previousHash := p.Hashes[(sequence-1)%BufferSize]
	expectedHash := sha256.New()
	expectedHash.Write([]byte(previousHash))
	expectedHashBytes := expectedHash.Sum(nil)
	expectedHashString := hex.EncodeToString(expectedHashBytes)

	// If the hash is verified, drop it by setting it to an empty string (or overwrite on next insert)
	isValid := expectedHashString == p.Hashes[sequence%BufferSize]
	if isValid {
		p.Hashes[sequence%BufferSize] = "" // Mark hash as used by clearing it
	}

	return isValid
}

// StartGenerating runs the PoH generator in a loop as fast as possible.
func (p *PoH) StartGenerating() {
	for {
		p.GenerateNextHash()
	}
}
