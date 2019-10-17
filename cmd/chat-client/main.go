package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/patrickmcnamara/chat"
	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	// set up config directory
	configDir, _ := os.UserConfigDir()
	configDir += "/chat"
	_ = os.MkdirAll(configDir, 0700)

	// load or generate user files
	cfg, err := loadConfig(configDir)
	chk(err)
	u := loadUser(configDir)
	cs := loadContacts(configDir)

	// connect to server as user
	cli := chat.NewClient(u)
	err = cli.Dial(cfg)
	chk(err)
	defer cli.Close()

	// setup terminal
	oldState, _ := terminal.MakeRaw(0)
	defer terminal.Restore(0, oldState)
	term := terminal.NewTerminal(os.Stdin, "> ")

	// receive messages
	go func() {
		for {
			// receive message
			msg, err := cli.Recv()
			if err != nil {
				fmt.Fprintln(term, err)
				break
			}

			// check for matching contact
			var contactName string
			for name, pubKey := range cs {
				if bytes.Equal(msg.From[:], pubKey[:]) {
					contactName = name
					continue
				}
			}
			if contactName != "" {
				fmt.Fprintf(term, "[%s] %s: %s\n", time.Now().Format(time.RFC3339), contactName, msg.Msg)
			} else {
				fmt.Fprintf(term, "[%s] %X: %s\n", time.Now().Format(time.RFC3339), msg.From[:8], msg.Msg)
			}
		}
	}()

	// send messages
	var contactName string
	for {
		// read a line
		line, err := term.ReadLine()
		if err != nil {
			break
		}

		// parse possible command
		lr := csv.NewReader(strings.NewReader(line))
		lr.Comma = ' '
		toks, _ := lr.Read()

		// switch contact possibly
		if len(toks) == 2 && toks[0] == "/msg" {
			// check if contact exists
			if _, ok := cs[toks[1]]; !ok {
				fmt.Fprintln(term, "Contact doesn't exist. Add it to your contacts file.")
				continue
			}
			// if it does, set it as the current contact
			contactName = toks[1]
			fmt.Fprintf(term, "Setting contact to %q.\n", contactName)
			continue
		}

		// check if there is a contact
		if contactName == "" {
			fmt.Fprintln(term, "No contact selected. Use \"/msg CONTACT_NAME\".")
			continue
		}

		// otherwise send message to contact
		_ = cli.Send(chat.Message{
			To:   cs[contactName],
			From: u.PublicKey,
			Msg:  []byte(line),
		})
	}
}

func loadConfig(configDir string) (string, error) {
	fn := configDir + "/config"
	f, err := os.Open(fn)
	if err != nil {
		return "", fmt.Errorf("no server specified in %q", fn)
	}
	b, _ := ioutil.ReadAll(f)
	return strings.TrimSpace(string(b)), nil
}

func loadUser(configDir string) (u chat.User) {
	fn := configDir + "/profile"
	f, _ := os.OpenFile(fn, os.O_RDWR|os.O_CREATE, 0600)
	dec := json.NewDecoder(f)
	err := dec.Decode(&u)
	if err != nil {
		u = chat.NewUser()
		enc := json.NewEncoder(f)
		enc.Encode(u)
	}
	return
}

func loadContacts(configDir string) (cs map[string][32]byte) {
	cs = make(map[string][32]byte)
	fn := configDir + "/contacts"
	f, _ := os.OpenFile(fn, os.O_RDWR|os.O_CREATE, 0600)
	dec := json.NewDecoder(f)
	err := dec.Decode(&cs)
	if err != nil {
		enc := json.NewEncoder(f)
		_ = enc.Encode(cs)
	}
	return
}

func chk(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "chat:", err)
		os.Exit(1)
	}
}
