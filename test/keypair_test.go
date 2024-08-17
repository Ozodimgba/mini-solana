package test

import (
    "log"
    "testing"
    "github.com/Ozodimgba/mini-solana/keypair"
)

func TestGenerateKeypair(t *testing.T) {
    kp, err := keypair.GenerateKeypair()
    if err != nil {
        t.Fatalf("Failed to generate keypair: %v", err)
    }

    log.Printf("Example Public Key (Base58): %s", kp.PublicKey)
    log.Printf("Example Private Key (Base58): %s", kp.GetPrivateKeyBase58())

    if len(kp.PublicKey) == 0 {
        t.Fatal("Public key is empty")
    }
}

func TestSignAndVerify(t *testing.T) {
    kp, _ := keypair.GenerateKeypair()
    message := []byte("Hello, World!")

    signature := kp.Sign(message)
    log.Printf("Example Signature (Base58): %s", signature)

    if !kp.Verify(message, signature) {
        t.Fatal("Signature verification failed")
    }

    // Test with incorrect message
    if kp.Verify([]byte("Wrong message"), signature) {
        t.Fatal("Signature verified for incorrect message")
    }

    // Test with incorrect signature
    if kp.Verify(message, "incorrect signature") {
        t.Fatal("Incorrect signature was verified")
    }
}