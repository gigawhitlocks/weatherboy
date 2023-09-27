package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type Message struct {
	Type string `json:"type"`
}

func RainStartEvent(b []byte) string {
	return fmt.Sprintf("RainStartEvent %s", string(b))
}

func LightningStrikeEvent(outb []byte) string {
	return fmt.Sprintf("LightningStrikeEvent %s", string(outb))
}
func RapidWind(outb []byte) string {
	type RapidWindMsg struct {
		Observation [3]any `json:"ob"`
	}
	r := new(RapidWindMsg)
	err := json.Unmarshal(outb, &r)
	if err != nil {
		return fmt.Sprintf("error unmarshaling rapid wind message: %s", err)
	}
	return fmt.Sprintf("Wind: Timestamp %d, %.01f m/s, %.01f deg",
		int(r.Observation[0].(float64)),
		r.Observation[1],
		r.Observation[2])
}

func ObservationAir(outb []byte) string {
	return fmt.Sprintf("ObservationAir %s", string(outb))
}

func ObservationSky(outb []byte) string {
	return fmt.Sprintf("ObservationSky %s", string(outb))
}

func ObservationTempest(outb []byte) string {
	type Obs struct {
		Observation [][17]any `json:"obs"`
	}
	o := new(Obs)
	err := json.Unmarshal(outb, &o)
	if err != nil {
		return fmt.Sprintf("ERROR %s", err)
	}
	r := o.Observation[0]

	return fmt.Sprintf(`
	Time Epoch %d s
	Wind Lull %d m/s
	Wind Avg %d m/s
	Wind Direction	%d Degrees
	Wind Sample Interval %d	seconds
	Station Pressure	MB
	Air Temperature	C
	Relative Humidity	%0.1f
	Illuminance	%d Lux
	UV	Index %d
	Solar Radiation	%d W/m^2
	Rain amount over previous minute	%d mm
	Precipitation Type	0 = none, 1 = rain, 2 = hail, 3 = rain + hail (experimental)
	Lightning Strike Avg Distance	km
	Lightning Strike Count	
	Battery	Volts
	Report Interval	Minutes
	`, int(r[0]),
	int(r[1]),
	int(r[2]),
	int(r[5]),
	int(r[6]),
	int(r[7]),
	int(r[8]),
	int(r[9]),
	int(r[13]),
	int(r[14]),
	int(r[15]),
	int(r[16]),
	int(r[17]))
}

func DeviceStatus(outb []byte) string {
	return fmt.Sprintf("DeviceStatus %s", string(outb))
}

func HubStatus(outb []byte) string {
	return fmt.Sprintf("HubStatus %s", string(outb))
}

func main() {
	ln, err := net.ListenUDP("udp", &net.UDPAddr{Port: 50222})
	if err != nil {
		fmt.Printf("opening socket: %s", err)
		os.Exit(1)
	}

	for {
		outb := make([]byte, 1024)
		n, err := ln.Read(outb)
		if err != nil {
			fmt.Printf("accepting cxn: %s", err)
			continue
		}
		outb = outb[:n]

		encodedMessageType := new(Message)
		err = json.Unmarshal(outb, &encodedMessageType)
		if err != nil {
			fmt.Printf("failed unmarshal: %s", err)
			continue
		}
		var outs string
		switch messageType := encodedMessageType.Type; messageType {
		case "evt_precip":
			outs = RainStartEvent(outb)
		case "evt_strike":
			outs = LightningStrikeEvent(outb)
		case "rapid_wind":
			outs = RapidWind(outb)
		case "obs_air":
			outs = ObservationAir(outb)
		case "obs_sky":
			outs = ObservationSky(outb)
		case "obs_st":
			outs = ObservationTempest(outb)
		case "device_status":
			outs = DeviceStatus(outb)
		case "hub_status":
			outs = HubStatus(outb)
		default:
			fmt.Printf("UNKNOWN MESSAGE TYPE %s", string(outb))
		}

		fmt.Println(outs)
	}
}
