package main

import (
	"encoding/json"
	"fmt"
	"time"
)

type RainStartEvent struct {
	Time time.Time
}

func (r *RainStartEvent) String() string {
	return r.Time.String()
}

func HandleRainStartEvent(b []byte) (*RainStartEvent, error) {
	o := new(RainStartEvent)
	err := json.Unmarshal(b, &o)
	if err != nil {
		return nil, fmt.Errorf("ERROR %w", err)
	}
	return o, nil
}
