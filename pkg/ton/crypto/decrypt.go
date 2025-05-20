package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"
	"golang.org/x/crypto/nacl/box"
)

// DecryptX25519AESCBCMessage decrypts a base64-encoded AES-CBC encrypted message using a shared secret derived
// from the sender's X25519 public key and the recipient's X25519 private key.
//
// It assumes the message was encrypted using the following process:
//  1. The sender computes a shared secret via ECDH (X25519).
//  2. The first 16 bytes of this shared secret are used as an AES-128 key.
//  3. The message is encrypted using AES-CBC with a 16-byte IV prepended to the ciphertext.
//  4. PKCS#7 padding is applied to the plaintext before encryption.
//
// Parameters:
// - encryptedBase64: base64-encoded string of IV (16 bytes) + ciphertext.
// - senderPublicKeyBase64: base64-encoded sender’s X25519 public key (32 bytes).
// - recipientPrivateKeyBase64: base64-encoded recipient’s X25519 private key (32 bytes).
//
// Returns:
// - Decrypted plaintext message as a string.
// - An error if decryption or decoding fails.
//
// Note: This function is compatible with NaCl box-style ECDH using X25519.
func DecryptX25519AESCBCMessage(encryptedBase64, senderPublicKeyBase64 string, recipientPrivateKey []byte) (string, error) {
	// Decode inputs from Base64
	encryptedMessage, err := base64.StdEncoding.DecodeString(encryptedBase64)
	if err != nil {
		return "", fmt.Errorf("failed to decode encrypted message: %w", err)
	}

	senderPublicKey, err := base64.StdEncoding.DecodeString(senderPublicKeyBase64)
	if err != nil {
		return "", fmt.Errorf("failed to decode sender public key: %w", err)
	}

	// Validate lengths
	if len(encryptedMessage) < 16 {
		return "", errors.New("encrypted message too short")
	}
	if len(senderPublicKey) != 32 || len(recipientPrivateKey) != 32 {
		return "", errors.New("invalid key length (must be 32 bytes)")
	}

	// Extract IV and ciphertext
	iv := encryptedMessage[:16]
	ciphertext := encryptedMessage[16:]

	// Compute shared secret (X25519 ECDH)
	var sharedSecret [32]byte
	var senderPub [32]byte
	var recipientPriv [32]byte
	copy(senderPub[:], senderPublicKey)
	copy(recipientPriv[:], recipientPrivateKey)

	box.Precompute(&sharedSecret, &senderPub, &recipientPriv)
	aesKey := sharedSecret[:16] // Use first 16 bytes for AES-128

	// Decrypt ciphertext using AES-CBC
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", fmt.Errorf("failed to create AES cipher: %w", err)
	}

	if len(ciphertext)%aes.BlockSize != 0 {
		return "", errors.New("ciphertext is not a multiple of the AES block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)

	// Remove PKCS#7 padding
	plaintext, err = pkcs7Unpad(plaintext, aes.BlockSize)
	if err != nil {
		return "", fmt.Errorf("unpadding failed: %w", err)
	}

	return string(plaintext), nil
}

// pkcs7Unpad removes PKCS#7 padding from decrypted plaintext
func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	if len(data) == 0 || len(data)%blockSize != 0 {
		return nil, errors.New("invalid padded data")
	}

	paddingLen := int(data[len(data)-1])
	if paddingLen == 0 || paddingLen > blockSize {
		return nil, errors.New("invalid padding length")
	}

	// Check padding bytes
	for i := len(data) - paddingLen; i < len(data); i++ {
		if data[i] != byte(paddingLen) {
			return nil, errors.New("invalid padding bytes")
		}
	}

	return data[:len(data)-paddingLen], nil
}
