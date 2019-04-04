package main

import (
	"fmt"
	"log"

	"github.com/d2r2/go-dht"
	logger "github.com/d2r2/go-logger"
)

var lg = logger.NewPackageLogger("main",
	logger.InfoLevel,
)

func main() {
	logger.ChangePackageLogLevel("dht", logger.InfoLevel)
	// temperature, humidity, retried, err :=
	// 	dht.ReadDHTxxWithRetry(dht.DHT22, 22, false, 10)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// // Print temperature and humidity
	// fmt.Printf("Temperature = %v*C, Humidity = %v%% (retried %d times)\n",
	// 	temperature, humidity, retried)
	t, h := getSensorData()
	fmt.Printf("Temp is %v, humidity is: %v", t, h)
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
