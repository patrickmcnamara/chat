package main

import (
	"fmt"
	"os"

	"github.com/patrickmcnamara/chat"
)

func main() {
	srv := chat.NewServer()
	err := srv.ListenAndServe(":6969")
	chk(err)
	defer srv.Close()
}

func chk(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "chat:", err)
		os.Exit(1)
	}
}
