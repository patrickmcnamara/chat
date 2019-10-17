package chat

import (
	"encoding/json"
)

// Message represents a message sent from one user to another. 'To' is the
// public key of the user the message is for and 'From' is the public key of the
// user sending the message.
type Message struct {
	To   [32]byte `json:"to"`
	From [32]byte `json:"from"`
	Msg  []byte   `json:"msg"`
}

func (msg Message) seal(sharedKey [32]byte) messageSealed {
	data, _ := json.Marshal(msg)
	ciphertext, nonce := encryptData(data, sharedKey)
	return messageSealed{msg.To, msg.From, nonce, ciphertext}
}

type messageSealed struct {
	To    [32]byte `json:"to"`
	From  [32]byte `json:"from"`
	Nonce [24]byte `json:"nonce"`
	Data  []byte   `json:"data"`
}

func (ms messageSealed) open(sharedKey [32]byte) (msg Message, err error) {
	plaintext, err := decryptData(ms.Data, sharedKey, ms.Nonce)
	_ = json.Unmarshal(plaintext, &msg)
	return
}
