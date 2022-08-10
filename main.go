package main

import (
	"binance_collector/binance"
	"log"
	"os"
	"strings"
)

func main() {
	serverURL := os.Getenv("INFLUX_URL")
	authToken := os.Getenv("INFLUX_TOKEN")
	bucket := os.Getenv("INFLUX_BUCKET")
	symbols := os.Getenv("SYMBOLS")

	serverURL = strings.ReplaceAll(serverURL, "'", "")
	authToken = strings.ReplaceAll(authToken, "'", "")
	bucket = strings.ReplaceAll(bucket, "'", "")
	symbols = strings.ReplaceAll(symbols, "'", "")

	log.Printf("INFLUX_URL : %s", serverURL)
	log.Printf("INFLUX_TOKEN : %s", authToken)
	log.Printf("INFLUX_BUCKET : %s", bucket)
	log.Printf("SYMBOLS: %v", symbols)

	bin := new(binance.BinanceDatabase)
	bin.Init(symbols, 100, serverURL, authToken, bucket)

	if anyErr := bin.Start(); anyErr != nil {
		log.Print(anyErr)
	}

	bin.CloseWs()

	log.Println("Complete....")
}

func main2() {

	// I hope it more fast ....

	bin := new(binance.Binance)
	bin.Init("BtcUsdT", 1000, 100, true, true)

	if anyErr := bin.Start(20); anyErr != nil {
		log.Print(anyErr)
	}

	bin.CloseWs()

}
