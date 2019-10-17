package chat

import (
	"crypto/rand"

	"golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/crypto/curve25519"
)

func generateX25519KeyPair() (pubKey [32]byte, priKey [32]byte) {
	rand.Read(priKey[:])
	pubKey[0] &= 248
	pubKey[31] &= 127
	pubKey[31] |= 64
	curve25519.ScalarBaseMult(&pubKey, &priKey)
	return
}

func generateX25519SharedKey(peerKey, priKey [32]byte) (sharedKey [32]byte) {
	curve25519.ScalarMult(&sharedKey, &priKey, &peerKey)
	return
}

func encryptData(plaintext []byte, sharedKey [32]byte) (ciphertext []byte, nonce [24]byte) {
	rand.Read(nonce[:])
	aead, _ := chacha20poly1305.NewX(sharedKey[:])
	ciphertext = aead.Seal(nil, nonce[:], plaintext, nil)
	return
}

func decryptData(ciphertext []byte, sharedKey [32]byte, nonce [24]byte) (plaintext []byte, err error) {
	aead, _ := chacha20poly1305.NewX(sharedKey[:])
	plaintext, err = aead.Open(nil, nonce[:], ciphertext, nil)
	return
}
