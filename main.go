package main

import (
	"binance_collector/binance"
	"log"
	"time"
)

func main() {
	// I hope it more fast ....

	bin := new(binance.BinanceSingle).Init("BtcUsdT", "1000ms", 5000)

	if anyErr := bin.Any(10000); anyErr != nil {
		log.Print(anyErr)
	}

	bin.CloseWS()

}

func main2() {
	bin := new(binance.Binance).Init("BtcUsdT")

	log.Printf("Started: %s", time.Now().UTC())
	bin.CollectDepth()
	log.Printf("Collected: %s", time.Now().UTC())

	bin.FetchDepth(100)
	log.Printf("Collect Complete: %s", time.Now().UTC())

	bin.CloseWS()
}
