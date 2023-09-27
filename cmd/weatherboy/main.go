package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	ln, err := net.ListenUDP("udp", &net.UDPAddr{Port: 50222})
	if err != nil {
		fmt.Printf("opening socket: %s", err)
		os.Exit(1)
	}

	outb := make([]byte, 128)
	for {
		_, err := ln.Read(outb)
		if err != nil {
			fmt.Printf("accepting cxn: %s", err)
		}

		fmt.Printf("%s\n", string(outb))

	}
}
