package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

// Transaction represents a simple transaction in the blockchain.
type Transaction struct {
	Sender    string
	Receiver  string
	Amount    int
	Signature string
}

// NewTransaction creates a new transaction.
func NewTransaction(sender, receiver string, amount int, privateKey string) (*Transaction, error) {
	tx := &Transaction{
		Sender:   sender,
		Receiver: receiver,
		Amount:   amount,
	}

	// Generate a signature using the private key (this is a placeholder function)
	signature, err := signTransaction(tx, privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign transaction: %v", err)
	}

	tx.Signature = signature
	return tx, nil
}

// signTransaction simulates signing a transaction (you'd replace this with real cryptography).
func signTransaction(tx *Transaction, privateKey string) (string, error) {
	// Simplified signing process
	txData := fmt.Sprintf("%s:%s:%d", tx.Sender, tx.Receiver, tx.Amount)
	hash := sha256.Sum256([]byte(txData + privateKey))
	return hex.EncodeToString(hash[:]), nil
}

// VerifyTransactionSignature verifies that the transaction's signature is valid.
func VerifyTransactionSignature(tx *Transaction, publicKey string) bool {
	// Simplified verification (replace with real public key verification)
	txData := fmt.Sprintf("%s:%s:%d", tx.Sender, tx.Receiver, tx.Amount)
	expectedHash := sha256.Sum256([]byte(txData + publicKey))
	expectedSignature := hex.EncodeToString(expectedHash[:])
	return tx.Signature == expectedSignature
}

// Hash returns the hash of the transaction.
func (tx *Transaction) Hash() string {
	txData := fmt.Sprintf("%s:%s:%d", tx.Sender, tx.Receiver, tx.Amount)
	hash := sha256.Sum256([]byte(txData))
	return hex.EncodeToString(hash[:])
}

// SerializeTransaction serializes the transaction to a JSON string.
func SerializeTransaction(tx *Transaction) (string, error) {
	data, err := json.Marshal(tx)
	if err != nil {
		return "", fmt.Errorf("failed to serialize transaction: %v", err)
	}
	return string(data), nil
}

// DeserializeTransaction deserializes a JSON string into a Transaction.
func DeserializeTransaction(data string) (*Transaction, error) {
	var tx Transaction
	err := json.Unmarshal([]byte(data), &tx)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize transaction: %v", err)
	}
	return &tx, nil
}

// ValidateTransaction checks if a transaction is valid.
func ValidateTransaction(tx *Transaction, publicKey string) bool {
	// For simplicity, we're only verifying the signature here.
	return VerifyTransactionSignature(tx, publicKey)
}
