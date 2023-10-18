package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"text/template"
	"time"
)

//go:embed observation.html.tmpl
var observationHTMLTemplate string

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
	PrecipType           string
	LightningAvgDistance float64
	LightningCount       int64
	Battery              float64
	ReportInterval       int64

	html string
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
	windLullMPH := r[1].(float64) * 2.23693629
	windAvgMPH := r[2].(float64) * 2.23693629
	windGustMPH := r[3].(float64) * 2.23693629
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
		WindDirection:        int64(r[4].(float64)),
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

func (o *Observation) HTML() string {
	if o == nil {
		return fmt.Sprintf(`<p style="color: %s" id="loading" class="loading" hx-swap="outerHTML" hx-trigger="every 1s" hx-get="/update">No observation has yet been recorded.</p>`,
			func() string {
				if time.Now().Unix()%2 == 0 {
					return "grey"
				}
				return "black"
			}())
	}

	if o.html == "" {
		t := template.Must(template.New("observation").Parse(observationHTMLTemplate))
		b := new(bytes.Buffer)
		err := t.Execute(b, o)
		if err != nil {
			return fmt.Sprintf("ERROR (output) %s", err)
		}

		o.html = b.String()
	}

	return o.html
}

func (o *Observation) String() string {
	const observation = `{{ .Time }}
Wind Lull                        {{.WindLull | printf "%.01f" }} mph
Wind Avg                         {{ .WindAvg | printf "%.01f" }} mph
Wind Gust                        {{ .WindGust | printf "%.01f" }} mph
Wind Direction                   {{ .WindDirection }} Degrees
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
