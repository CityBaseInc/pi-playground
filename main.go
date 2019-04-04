package main

import (
	"fmt"
	"log"
	"time"

	"github.com/d2r2/go-dht"
	logger "github.com/d2r2/go-logger"
)

var lg = logger.NewPackageLogger("main",
	logger.InfoLevel,
)

func main() {
	logger.ChangePackageLogLevel("dht", logger.ErrorLevel)
	// temperature, humidity, retried, err :=
	// 	dht.ReadDHTxxWithRetry(dht.DHT22, 22, false, 10)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// // Print temperature and humidity
	// fmt.Printf("Temperature = %v*C, Humidity = %v%% (retried %d times)\n",
	// 	temperature, humidity, retried)
	for {
		t, h := getSensorData()
		fmt.Printf("Temp is %v, humidity is: %v\n", t, h)
		time.Sleep(10 * time.Second)
	}
	return
}

func getSensorData() (float32, float32) {
	temperature, humidity, _, err :=
		dht.ReadDHTxxWithRetry(dht.DHT22, 22, false, 10)
	if err != nil {
		log.Fatal(err)
	}
	return temperature, humidity
}
