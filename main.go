package main

import (
	"fmt"
	"log"
	"time"

	"github.com/d2r2/go-dht"
	logger "github.com/d2r2/go-logger"
	"github.com/influxdata/influxdb/client/v2"
)

var lastHumid float64
var lastTemp float64

var lg = logger.NewPackageLogger("main",
	logger.InfoLevel,
)

func main() {
	logger.ChangePackageLogLevel("dht", logger.ErrorLevel)
	for {
		t, h := getSensorData()
		fmt.Printf("Temp is %v, humidity is: %v\n", t, h)
		writeInflux(t, h)
		time.Sleep(10 * time.Second)
	}
	return
}

func getSensorData() (float32, float32) {
	temperature, humidity, _, err :=
		dht.ReadDHTxxWithRetry(dht.DHT22, 22, false, 10)
	if err != nil {
		log.Error(err)
		return lastTemp, lastHumid
	}
	lastHumid = humidity
	lastTemp = temperature
	return temperature, humidity
}

func writeInflux(temp float32, hum float32) {
	// Make client
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
	tags := map[string]string{"weather": "current"}
	fields := map[string]interface{}{
		"temp":     temp,
		"humidity": hum,
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
