package main

import (
	// "fmt"
	"time"

	"github.com/Ozodimgba/mini-solana/consensus"
)

func main() {
// Initialize PoH with a seed
pohGenerator := poh.NewPoH("initial-seed")

// Start generating hashes in a separate goroutine
go pohGenerator.StartGenerating()

// Allow some time for hashes to be generated
time.Sleep(2 * time.Second)

// // Verify the first few hashes
// for i := 1; i < len(pohGenerator.Hashes); i++ {
// 	if pohGenerator.VerifyHash(i) {
// 		fmt.Printf("Hash at sequence %d is valid.\n", i)
// 	} else {
// 		fmt.Printf("Hash at sequence %d is INVALID!\n", i)
// 	}
// }
}
// 5QBS5XCu4YeSgq6UkPV2XZhx1DHc38Zf11ApQkX7HUiv5nJtwjr3Wa3AE4PWKqMgsuAvU3oq6AQC9W2j457viQSc
// 2gShtyERw1caHAGm3pBWycGQQ562ZCgKGBBHy2CtEUS1rtkcFjzJcdiL35qAWQHCmkpszG2Yhv76wXXwBRrMWgpj