package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type Event struct {
	Type string `json:"type"`
}

func main() {
	ln, err := net.ListenUDP("udp", &net.UDPAddr{Port: 50222})
	if err != nil {
		fmt.Printf("opening socket: %s", err)
		os.Exit(1)
	}

	outb := make([]byte, 1024)
	for {
		n, err := ln.Read(outb)
		if err != nil {
			fmt.Printf("accepting cxn: %s", err)
			continue
		}

		encodedMessageType := new(Event)
		err = json.Unmarshal(outb[:n], &encodedMessageType)
		if err != nil {
			fmt.Printf("failed unmarshal: %s", err)
			continue
		}
		messageType := encodedMessageType.Type
		fmt.Println(messageType)
		fmt.Println(string(outb))
	}
}
