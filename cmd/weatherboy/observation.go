package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"text/template"
	"time"
)

type Observation struct {
	Time                 time.Time
	WindLull             float64
	WindAvg              float64
	WindGust             float64
	WindDirection        int64
	WindSampleInterval   int64
	StationPressure      float64
	AirTemperature       float64
	RelativeHumidity     float64
	Illuminance          float64
	UV                   float64
	SolarRadiation       float64
	RainPrevMin          float64
	PrecipType           float64
	LightningAvgDistance float64
	LightningCount       int64
	Battery              float64
	ReportInterval       int64
}

func HandleObservation(inb []byte) (*Observation, error) {
	type Obs struct {
		Observation [][18]any `json:"obs"`
	}
	o := new(Obs)
	err := json.Unmarshal(inb, &o)
	if err != nil {
		return nil, fmt.Errorf("ERROR %w", err)
	}
	r := o.Observation[0]
	timestamp := time.Unix(int64(r[0].(float64)), 0)
	return &Observation{
		Time:                 timestamp,
		WindLull:             r[1].(float64),
		WindAvg:              r[2].(float64),
		WindGust:             r[3].(float64),
		WindDirection:        int64(r[4].(float64)),
		WindSampleInterval:   int64(r[5].(float64)),
		StationPressure:      r[6].(float64),
		AirTemperature:       r[7].(float64),
		RelativeHumidity:     r[8].(float64),
		Illuminance:          r[9].(float64),
		UV:                   r[10].(float64),
		SolarRadiation:       r[11].(float64),
		RainPrevMin:          r[12].(float64),
		PrecipType:           r[13].(float64),
		LightningAvgDistance: r[14].(float64),
		LightningCount:       int64(r[15].(float64)),
		Battery:              r[16].(float64),
		ReportInterval:       int64(r[17].(float64)),
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
