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
	WindDirection        string
	WindSampleInterval   int64
	StationPressure      float64
	AirTemperature       float64
	RelativeHumidity     float64
	Illuminance          float64
	UV                   float64
	SolarRadiation       float64
	RainPrevMin          float64
	PrecipType           string
	LightningAvgDistance float64
	LightningCount       int64
	Battery              float64
	ReportInterval       int64
}

const (
	north     = 0
	northeast = iota
	east
	southeast
	south
	southwest
	west
	northwest
)

var windDirections [8]string = [8]string{"🡸"," 🡺","🡹 ","🡻" ,"🡼", "🡽"," 🡾","🡿 "}

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
	windLullMPH := r[1].(float64) * 2.23693629
	windAvgMPH := r[2].(float64) * 2.23693629
	windGustMPH := r[3].(float64) * 2.23693629
	windDirection := func(w int64) string {
		// this could be math probably but then I would have to think
		// about it harde
		switch {
		case w < 23:
			return windDirections[north]
		case w >= 23 && w < 67:
			return windDirections[northeast]
		case w >= 67 && w < 112:
			return windDirections[east]
		case w >= 112 && w < 157:
			return windDirections[southeast]
		case w >= 157 && w < 202:
			return windDirections[south]
		case w >= 202 && w < 247:
			return windDirections[southwest]
		case w >= 247 && w < 292:
			return windDirections[west]
		case w >= 292 && w < 337:
			return windDirections[northwest]
		default:
			return windDirections[north]
		}
	}(int64(r[4].(float64)))
	tempF := (r[7].(float64) * 9.0 / 5.0) + 32

	precipType := func(t int) string {
		switch t {
		case 0:
			return "None"
		case 1:
			return "Rain"
		case 2:
			return "Hail"
		case 3:
			return "Hail and Rain"
		}
		return fmt.Sprintf("Unknown precipitation type (%d)", t)
	}(int(r[13].(float64)))

	return &Observation{
		Time:                 timestamp,
		WindLull:             windLullMPH,
		WindAvg:              windAvgMPH,
		WindGust:             windGustMPH,
		WindDirection:        windDirection,
		WindSampleInterval:   int64(r[5].(float64)),
		StationPressure:      r[6].(float64),
		AirTemperature:       tempF,
		RelativeHumidity:     r[8].(float64),
		Illuminance:          r[9].(float64),
		UV:                   r[10].(float64),
		SolarRadiation:       r[11].(float64),
		RainPrevMin:          r[12].(float64),
		PrecipType:           precipType,
		LightningAvgDistance: r[14].(float64),
		LightningCount:       int64(r[15].(float64)),
		Battery:              r[16].(float64),
		ReportInterval:       int64(r[17].(float64)),
	}, nil
}

func (o *Observation) String() string {
	const observation = `{{ .Time }}
Wind Lull                        {{.WindLull | printf "%.01f" }} mph
Wind Avg                         {{ .WindAvg | printf "%.01f" }} mph
Wind Gust                        {{ .WindGust | printf "%.01f" }} mph
Wind Direction                   {{ .WindDirection }}
Wind Sample Interval             {{ .WindSampleInterval }} s
Pressure                         {{ .StationPressure }} mb
Air Temperature                  {{ .AirTemperature | printf "%.1f" }} F
Relative Humidity                {{ .RelativeHumidity }}%
Illuminance                      {{ .Illuminance }} Lux
UV Index                         {{ .UV }}
Solar Radiation                  {{ .SolarRadiation }} W/m^2
Rain amount over previous minute {{ .RainPrevMin }}mm
Precipitation Type               {{ .PrecipType }}
Lightning Strike Avg Distance    {{ .LightningAvgDistance }} km
Lightning Strike Count           {{ .LightningCount }}
Battery	Volts                    {{ .Battery }}V
`

	t := template.Must(template.New("observation").Parse(observation))
	b := new(bytes.Buffer)
	err := t.Execute(b, o)
	if err != nil {
		return fmt.Sprintf("ERROR (output) %s", err)
	}
	return b.String()
}
