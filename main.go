package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/d2r2/go-dht"
	logger "github.com/d2r2/go-logger"
	"github.com/influxdata/influxdb/client/v2"
)

var lastHumid float32
var lastTemp float32

var lg = logger.NewPackageLogger("main",
	logger.InfoLevel,
)

func main() {
	logger.ChangePackageLogLevel("dht", logger.ErrorLevel)
	go outdoorLoop()
	for {
		t, h := getSensorData()
		fmt.Printf("Temp is %v, humidity is: %v\n", t, h)
		writeInflux(t, h, "indoor")
		time.Sleep(10 * time.Second)
	}
	return
}

func getSensorData() (float32, float32) {
	temperature, humidity, _, err :=
		dht.ReadDHTxxWithRetry(dht.DHT22, 22, false, 10)
	if err != nil {
		fmt.Println(err)
		return lastTemp, lastHumid
	}
	lastHumid = humidity
	lastTemp = temperature
	return temperature, humidity
}

func writeInflux(temp float32, hum float32, loc string) {
	// Make client
	tfixed := ((temp * 9) / 5) + 32
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr: "http://localhost:8086",
	})
	if err != nil {
		fmt.Println("Error creating InfluxDB Client: ", err.Error())
		return
	}
	defer c.Close()

	// Create a new point batch
	bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "pistats",
		Precision: "s",
	})

	// Create a point and add to batch
	tags := map[string]string{"location": loc}
	fields := map[string]interface{}{
		"temp":     temp,
		"humidity": hum,
		"tempf":    tfixed,
	}
	pt, err := client.NewPoint("weather", tags, fields, time.Now())
	if err != nil {
		fmt.Println("Error: ", err.Error())
		return
	}
	bp.AddPoint(pt)

	// Write the batch
	c.Write(bp)
	return
}

func outdoorLoop() {
	for {
		w := getOutdoorStats()
		if w.Main.Temp == 0 {
			continue
		}
		fmt.Printf("outdoor temp: %v hum: %v\n", w.Main.Temp, w.Main.Humidity)
		writeInflux(float32(w.Main.Temp), float32(w.Main.Humidity), "outdoor")
		time.Sleep(10 * time.Minute)
	}
}

func getOutdoorStats() WeatherResponse {
	myKey := os.Getenv("WEATHER_KEY")
	url := "http://api.openweathermap.org/data/2.5/weather?q=Chicago&APPID=" + myKey + "&units=metric"
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("cache-control", "no-cache")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))
	w := WeatherResponse{}
	json.Unmarshal(body, &w)
	fmt.Println(w)
	return w
}

type WeatherResponse struct {
	Coord struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"coord"`
	Weather []struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Base string `json:"base"`
	Main struct {
		Temp     float64 `json:"temp"`
		Pressure int     `json:"pressure"`
		Humidity int     `json:"humidity"`
		TempMin  float64 `json:"temp_min"`
		TempMax  int     `json:"temp_max"`
	} `json:"main"`
	Visibility int `json:"visibility"`
	Wind       struct {
		Speed float64 `json:"speed"`
		Deg   int     `json:"deg"`
	} `json:"wind"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Dt  int `json:"dt"`
	Sys struct {
		Type    int     `json:"type"`
		ID      int     `json:"id"`
		Message float64 `json:"message"`
		Country string  `json:"country"`
		Sunrise int     `json:"sunrise"`
		Sunset  int     `json:"sunset"`
	} `json:"sys"`
	ID   int    `json:"id"`
	Name string `json:"name"`
	Cod  int    `json:"cod"`
}
