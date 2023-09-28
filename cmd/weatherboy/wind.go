package main

import (
	"encoding/json"
	"fmt"
	"time"
)

type RapidWindEvent struct {
	Time      time.Time
	Speed     float64
	Direction int64
}

func (r *RapidWindEvent) String() string {
	return fmt.Sprintf("%s Wind %.001f m/s, %d deg",
		r.Time.String(), r.Speed, r.Direction)
}

func HandleRapidWindEvent(outb []byte) (*RapidWindEvent, error) {
	type RapidWindMsg struct {
		Observation [3]any `json:"ob"`
	}
	r := new(RapidWindMsg)
	err := json.Unmarshal(outb, &r)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling rapid wind message: %w", err)
	}

	timestamp := time.Unix(int64(r.Observation[0].(float64)), 0)
	return &RapidWindEvent{
		Time:      timestamp,
		Speed:     r.Observation[1].(float64),
		Direction: int64(r.Observation[2].(float64)),
	}, nil

}
