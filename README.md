wouldn't you like to know, weatherboy

# Weatherboy

Weatherboy is a little tool to display data from a Weatherflow Tempest device.
It only uses the locally available UDP API, so it does not need an API key. It just listens on the network, and displays what it sees.

Run with no arguments, it will show the current observation, once it receives one, and update automatically whenever a new one is received. Run with -daemon it will expose a simple HTMX website to display the data, and also makes it available as JSON. 

```
~/weatherboy (master|✔) ➜ ./bin/weatherboy
2024-10-25 23:30:05 -0500 CDT
Wind Lull                        0.4 mph
Wind Avg                         1.4 mph
Wind Gust                        2.8 mph
Wind Direction                   197 Degrees
Wind Sample Interval             3 s
Pressure                         994.62 mb
Air Temperature                  73.7 F
Relative Humidity                68.11%
Illuminance                      0 Lux
UV Index                         0
Solar Radiation                  0 W/m^2
Rain amount over previous minute 0mm
Precipitation Type               None
Lightning Strike Avg Distance    0 km
Lightning Strike Count           0
Battery Volts                    2.668V
```

Also, it writes events to /tmp/weatherboy.log

```
~/weatherboy (master|✔) ➜ tail -n15 /tmp/weatherboy.log | grep -i temp | jq .
{
  "Time": "2024-10-25T23:30:05-05:00",
  "WindLull": 0.35790980640000003,
  "WindAvg": 1.4092698627,
  "WindGust": 2.7514316367,
  "WindDirection": 197,
  "WindSampleInterval": 3,
  "StationPressure": 994.62,
  "AirTemperature": 73.67,
  "RelativeHumidity": 68.11,
  "Illuminance": 0,
  "UV": 0,
  "SolarRadiation": 0,
  "RainPrevMin": 0,
  "PrecipType": "None",
  "LightningAvgDistance": 0,
  "LightningCount": 0,
  "Battery": 2.668,
  "ReportInterval": 1
}
```

It is rudimentary however it could be used as a good starting point to do something more interesting with this data.
