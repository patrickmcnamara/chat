package chat

import (
	"encoding/json"
	"net"
)

type client struct {
	user User
	conn *net.TCPConn
}

// Client is a chat client that connects to chat server for a user.
type Client interface {
	Dial(address string) error
	Send(msg Message) error
	Recv() (Message, error)
	Close()
}

// NewClient creates a new chat client for the given user.
func NewClient(u User) Client {
	return &client{user: u}
}

func (cli *client) Dial(address string) error {
	// connect to server
	raddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return err
	}
	conn, err := net.DialTCP("tcp", nil, raddr)
	if err != nil {
		return err
	}
	cli.conn = conn

	// send public key
	enc := json.NewEncoder(conn)
	err = enc.Encode(cli.user.PublicKey)
	return err
}

func (cli *client) Send(msg Message) error {
	// encrypt message using shared key from diffie-hellman key exchange
	sharedKey := generateX25519SharedKey(msg.To, cli.user.PrivateKey)
	ms := msg.seal(sharedKey)

	// send encrypted message to server
	enc := json.NewEncoder(cli.conn)
	err := enc.Encode(ms)
	return err
}

func (cli *client) Recv() (msg Message, err error) {
	// receive encrypted message from server
	dec := json.NewDecoder(cli.conn)
	var ms messageSealed
	err = dec.Decode(&ms)
	if err != nil {
		return
	}

	// decrypt message using shared secret from diffie-hellman key exchange
	sharedKey := generateX25519SharedKey(ms.To, cli.user.PrivateKey)
	msg, err = ms.open(sharedKey)
	return
}

func (cli *client) Close() {
	cli.conn.Close()
}
