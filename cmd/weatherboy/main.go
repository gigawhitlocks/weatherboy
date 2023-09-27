package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	ln, err := net.Listen("udp", ":50222")
	if err != nil {
		fmt.Printf("opening socket: %s", err)
		os.Exit(1)
	}

	for {

		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("accepting cxn: %s", err)
		}

		fmt.Printf("Got %s", conn)

	}
}
