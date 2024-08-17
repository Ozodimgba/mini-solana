package keypair

import (
    "crypto/ed25519"
    "crypto/rand"
    "github.com/mr-tron/base58"
)

type Keypair struct {
    PublicKey  string
    privateKey ed25519.PrivateKey
}

// GenerateKeypair creates a new keypair
func GenerateKeypair() (*Keypair, error) {
    publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
    if err != nil {
        return nil, err
    }

    return &Keypair{
        PublicKey:  base58.Encode(publicKey),
        privateKey: privateKey,
    }, nil
}

// Sign signs the given message with the private key
func (kp *Keypair) Sign(message []byte) string {
    signature := ed25519.Sign(kp.privateKey, message)
    return base58.Encode(signature)
}

// Verify checks if the signature is valid for the given message
func (kp *Keypair) Verify(message []byte, signatureBase58 string) bool {
    publicKey, err := base58.Decode(kp.PublicKey)
    if err != nil {
        return false
    }

    signature, err := base58.Decode(signatureBase58)
    if err != nil {
        return false
    }

    return ed25519.Verify(publicKey, message, signature)
}

// GetPrivateKeyBase58 returns the Base58 encoded private key (use cautiously!)
func (kp *Keypair) GetPrivateKeyBase58() string {
    return base58.Encode(kp.privateKey)
}