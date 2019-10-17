package chat

// User is a user of a chat server that connects through a chat client. It
// represents the public and private key pair used in end to end encryption.
// The public keys' of users are used as addresses when messaging.
type User struct {
	PublicKey  [32]byte `json:"publicKey"`
	PrivateKey [32]byte `json:"privateKey"`
}

// NewUser creates a new user with a new X25519 key pair.
func NewUser() User {
	pubKey, priKey := generateX25519KeyPair()
	return User{pubKey, priKey}
}
