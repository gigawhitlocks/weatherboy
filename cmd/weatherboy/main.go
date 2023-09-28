package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type RawUDPEvent struct {
	Type string `json:"type"`
}

type Event interface {
	String() string
}

type Dashboard struct {
	LastObservation *Observation
}

func (d Dashboard) Init() tea.Cmd {
	return nil
}

func (d Dashboard) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// switch t := msg.(type) {
	// case *Observation:
	// 	d.LastObservation = t
	// }
	return d, nil
}

func (d Dashboard) View() string {
	if d.LastObservation == nil {
		return "no data yet"
	}

	return d.LastObservation.String()
}

func main() {
	dash := &Dashboard{}

	p := tea.NewProgram(dash)
	go collector(p)
	if _, err := p.Run(); err != nil {
		fmt.Printf("fatal gui error: %v", err)
		os.Exit(1)
	}
}

func collector(p *tea.Program) {
	if _, err := os.Stat("/tmp/weatherboy.log"); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			os.Create("/tmp/weatherboy.log")
		}
	}

	logfile, err := os.OpenFile("/tmp/weatherboy.log", os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("opening logfile: %s", err)
	}
	defer logfile.Close()

	log := func(msg string) {
		_, err := fmt.Fprintf(logfile, "%s\n", msg)
		if err != nil {
			fmt.Printf("ERROR %s", err)
			os.Exit(1)
		}
	}

	ln, err := net.ListenUDP("udp", &net.UDPAddr{Port: 50222})
	if err != nil {
		log(fmt.Sprintf("opening socket: %s", err))
		os.Exit(1)
	}

	for {
		outb := make([]byte, 1024)
		n, err := ln.Read(outb)
		if err != nil {
			log(fmt.Sprintf("accepting cxn: %s", err))
			continue
		}
		outb = outb[:n]

		encodedMessageType := new(RawUDPEvent)
		err = json.Unmarshal(outb, &encodedMessageType)
		if err != nil {
			log(fmt.Sprintf("failed unmarshal: %s", err))
			continue
		}
		var ev Event
		switch messageType := encodedMessageType.Type; messageType {
		case "evt_precip":
			ev, err = HandleRainStartEvent(outb)
		case "evt_strike":
			ev, err = HandleLightningStrikeEvent(outb)
		case "rapid_wind":
			ev, err = HandleRapidWindEvent(outb)
		case "obs_st":
			o, err := HandleObservation(outb)
			if err != nil {
				log(fmt.Sprintf("ERROR: %s\n", err))
				continue
			}
			p.Send(o)
			ev = o
		case "device_status":
			ev, err = HandleDeviceStatusEvent(outb)
		case "hub_status":
			ev, err = HandleHubStatusEvent(outb)
		default:
			log(fmt.Sprintf("UNKNOWN MESSAGE TYPE %s", string(outb)))
		}

		if err != nil {
			log(fmt.Sprintf("ERROR: %s\n", err))
			continue
		}

		log(ev.String())
	}
}
