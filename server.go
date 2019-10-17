package chat

import (
	"encoding/json"
	"log"
	"net"
	"sync"
)

type server struct {
	conns map[[32]byte]*net.TCPConn
	lock  sync.Mutex
}

// Server is a chat server that routes messages between chat clients.
type Server interface {
	ListenAndServe(address string) error
	Close()
}

// NewServer creates a new chat server.
func NewServer() Server {
	return &server{conns: make(map[[32]byte]*net.TCPConn)}
}

func (srv *server) ListenAndServe(address string) error {
	// set up server
	laddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return err
	}
	lst, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		return err
	}
	defer lst.Close()

	// listen for connections
	for {
		conn, err := lst.AcceptTCP()
		if err != nil {
			return err
		}
		go func() {
			srv.handle(conn)
			defer conn.Close()
		}()
	}
}

func (srv *server) Close() {
	for _, conn := range srv.conns {
		conn.Close()
	}
}

func (srv *server) handle(conn *net.TCPConn) {
	// get public key of client
	var id [32]byte
	enc := json.NewDecoder(conn)
	err := enc.Decode(&id)
	if err != nil {
		return
	}

	// add public key to server's map of connections
	func() {
		srv.lock.Lock()
		srv.conns[id] = conn
		srv.lock.Unlock()
	}()
	defer func() {
		srv.lock.Lock()
		delete(srv.conns, id)
		srv.lock.Unlock()
	}()

	// log connection status
	log.Printf("%X connected", id[:8])
	defer log.Printf("%X disconnected", id[:8])

	// receive messages
	for {
		ms, err := srv.recv(conn)
		if err != nil {
			log.Printf("%s from %X", err, id[:8])
			return
		}
		err = srv.send(ms)
		if err != nil {
			log.Printf("%s from %X", err, id[:8])
			return
		}
	}
}

func (srv *server) recv(conn *net.TCPConn) (ms messageSealed, err error) {
	// receive encrypted message
	dec := json.NewDecoder(conn)
	err = dec.Decode(&ms)
	if err == nil {
		log.Printf("%X sent message to %X", ms.From[:8], ms.To[:8])
	}
	return
}

func (srv *server) send(ms messageSealed) error {
	// check if recipient is online
	conn, ok := srv.conns[ms.To]
	if !ok {
		return ErrNoSuchUser
	}

	// send encrypted message to client
	enc := json.NewEncoder(conn)
	err := enc.Encode(ms)
	if err == nil {
		log.Printf("%X received message from %X", ms.To[:8], ms.From[:8])
	}
	return err
}
