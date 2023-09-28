package main

import "time"

type RainStartEvent struct {
	Time time.Time
}

func (r *RainStartEvent) String() string {
	return r.Time.String()
}

func HandleRainStartEvent(b []byte) (*RainStartEvent, error) {
	return nil, nil
}
