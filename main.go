package main

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/Ozodimgba/mini-solana/poh"
)

const (
	benchmarkDuration = 60 * time.Second
	reportInterval    = 400 * time.Millisecond
)

func main() {
	var state poh.State
	rand.Read(state[:])

	fmt.Println("Starting continuous PoH generation and benchmarking...")
	fmt.Printf("Initial state: %x\n", state)

	startTime := time.Now()
	lastReport := startTime
	hashCount := uint64(0)

	for time.Since(startTime) < benchmarkDuration {
		// Perform a batch of hashes
		batchSize := uint(10000)
		state.Hash(batchSize)
		hashCount += uint64(batchSize)

		// Report progress at regular intervals
		if time.Since(lastReport) >= reportInterval {
			reportProgress(hashCount, lastReport)
			lastReport = time.Now()
		}
	}

	// Final report
	reportProgress(hashCount, startTime)
	fmt.Printf("Final state: %x\n", state)
}

func reportProgress(hashCount uint64, since time.Time) {
	duration := time.Since(since)
	rate := float64(hashCount) / duration.Seconds()
	fmt.Printf("Hashes: %d, Duration: %.2fs, Rate: %.2f hashes/second\n",
		hashCount, duration.Seconds(), rate)
}
// 5QBS5XCu4YeSgq6UkPV2XZhx1DHc38Zf11ApQkX7HUiv5nJtwjr3Wa3AE4PWKqMgsuAvU3oq6AQC9W2j457viQSc
// 2gShtyERw1caHAGm3pBWycGQQ562ZCgKGBBHy2CtEUS1rtkcFjzJcdiL35qAWQHCmkpszG2Yhv76wXXwBRrMWgpj