package main

import "fmt"

type LightningStrikeEvent struct {
	Body []byte
}

func (l *LightningStrikeEvent) String() string {
	return fmt.Sprintf("LightningStrikeEvent %s", string(l.Body))
}

func HandleLightningStrikeEvent(outb []byte) (*LightningStrikeEvent, error) {
	return &LightningStrikeEvent{Body: outb}, nil
}
