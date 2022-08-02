package main

import (
	"binance_collector/binance"
	"log"
	"time"
)

func main() {
	bin := new(binance.Binance).Init("BtcUsdT")

	log.Printf("Started: %s", time.Now().UTC())
	bin.CollectDepth()
	log.Printf("Collected: %s", time.Now().UTC())

	bin.FetchDepth(100)
	log.Printf("Fetch Complete: %s", time.Now().UTC())

}
