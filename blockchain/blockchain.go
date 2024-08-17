package blockchain

import (
	"crypto/sha256"
	"fmt"
	"strconv"
	"time"
)

type Block struct {
    Index        int
    Timestamp    int64
    PreviousHash string
    Hash         string
    Transactions []Transaction
}

type Blockchain struct {
    Blocks []Block
}

func NewBlockchain() *Blockchain {
    return &Blockchain{
        Blocks: []Block{genesisBlock()},
    }
}

func genesisBlock() Block {
    return Block{
        Index:     0,
        Timestamp: time.Now().Unix(),
        Hash:      "genesis",
    }
}

func (bc *Blockchain) AddBlock(transactions []Transaction) {
    lastBlock := bc.Blocks[len(bc.Blocks)-1]
    newBlock := Block{
        Index:        lastBlock.Index + 1,
        Timestamp:    time.Now().Unix(),
        PreviousHash: lastBlock.Hash,
        Transactions: transactions,
    }
    newBlock.Hash = calculateHash(newBlock)
    bc.Blocks = append(bc.Blocks, newBlock)
}

func calculateHash(block Block) string {
    // Simplified hash calculation
    record := strconv.Itoa(block.Index) + strconv.FormatInt(block.Timestamp, 10) + block.PreviousHash
    hash := sha256.Sum256([]byte(record))
    return fmt.Sprintf("%x", hash)
}
