package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"text/template"
)

type Message struct {
	Type string `json:"type"`
}

func ParseRainStartEvent(b []byte) string {
	return fmt.Sprintf("RainStartEvent %s", string(b))
}

func ParseLightningStrikeEvent(outb []byte) string {
	return fmt.Sprintf("LightningStrikeEvent %s", string(outb))
}

func ParseRapidWind(outb []byte) string {
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

func ParseObservationAir(outb []byte) string {
	return fmt.Sprintf("ObservationAir %s", string(outb))
}

func ParseObservationSky(outb []byte) string {
	return fmt.Sprintf("ObservationSky %s", string(outb))
}

type RapidWind struct {
	Time int64
}

type Observation struct {
	Time                 float64
	WindLull             float64
	WindAvg              float64
	WindGust             float64
	WindDirection        float64
	WindSampleInterval   float64
	StationPressure      float64
	AirTemperature       float64
	RelativeHumidity     float64
	Illuminance          float64
	UV                   float64
	SolarRadiation       float64
	RainPrevMin          float64
	PrecipType           float64
	LightningAvgDistance float64
	LightningCount       float64
	Battery              float64
	ReportInterval       float64
}

func ParseObservation(outb []byte) (*Observation, error) {
	type Obs struct {
		Observation [][18]any `json:"obs"`
	}
	o := new(Obs)
	err := json.Unmarshal(outb, &o)
	if err != nil {
		return nil, fmt.Errorf("ERROR %w", err)
	}
	r := o.Observation[0]
	return &Observation{
		Time:                 r[0].(float64),
		WindLull:             r[1].(float64),
		WindAvg:              r[2].(float64),
		WindGust:             r[3].(float64),
		WindDirection:        r[4].(float64),
		WindSampleInterval:   r[5].(float64),
		StationPressure:      r[6].(float64),
		AirTemperature:       r[7].(float64),
		RelativeHumidity:     r[8].(float64),
		Illuminance:          r[9].(float64),
		UV:                   r[10].(float64),
		SolarRadiation:       r[11].(float64),
		RainPrevMin:          r[12].(float64),
		PrecipType:           r[13].(float64),
		LightningAvgDistance: r[14].(float64),
		LightningCount:       r[15].(float64),
		Battery:              r[16].(float64),
		ReportInterval:       r[17].(float64),
	}, nil
}

func (o *Observation) String() string {
	const observation = `
Time Epoch {{ .Time }}s
Wind Lull {{.WindLull}} m/s
Wind Avg {{ .WindAvg }} m/s
Wind Gust {{ .WindGust }} m/s
Wind Direction	{{ .WindDirection }} Degrees
Wind Sample Interval {{ .WindSampleInterval }}s
Station Pressure {{ .StationPressure }}
Air Temperature	{{ .AirTemperature }} C
Relative Humidity	{{ .RelativeHumidity }}%
Illuminance	{{ .Illuminance }} Lux
UV	Index {{ .UV }}
Solar Radiation	{{ .SolarRadiation }} W/m^2
Rain amount over previous minute {{ .RainPrevMin }}mm
Precipitation Type {{.PrecipType}}	0 = none, 1 = rain, 2 = hail, 3 = rain + hail (experimental)
Lightning Strike Avg Distance	{{ .LightningAvgDistance }}km
Lightning Strike Count	{{ .LightningCount }}
Battery	Volts {{ .Battery }}V
Report Interval	{{ .ReportInterval }} Minutes
`

	t := template.Must(template.New("observation").Parse(observation))
	b := new(bytes.Buffer)
	err := t.Execute(b, o)
	if err != nil {
		return fmt.Sprintf("ERROR (output) %s", err)
	}
	return b.String()
}

func ParseDeviceStatus(outb []byte) string {
	return fmt.Sprintf("DeviceStatus")
}

func ParseHubStatus(outb []byte) string {
	return fmt.Sprintf("HubStatus")
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
			outs = ParseRainStartEvent(outb)
		case "evt_strike":
			outs = ParseLightningStrikeEvent(outb)
		case "rapid_wind":
			outs = ParseRapidWind(outb)
		case "obs_st":
			o, err := ParseObservation(outb)
			if err != nil {
				fmt.Printf("error parsing observation: %s", err)
				continue
			}
			outs = o.String()
		case "device_status":
			outs = ParseDeviceStatus(outb)
		case "hub_status":
			outs = ParseHubStatus(outb)
		default:
			fmt.Printf("UNKNOWN MESSAGE TYPE %s", string(outb))
		}

		fmt.Println(outs)
	}
}
