package main

import (
	"binance_collector/binance"
	"log"
	"time"
)

func main() {
	bin := new(binance.Binance).Init("btcusdt")

	log.Printf("Started: %s", time.Now().UTC())
	bin.CollectDepth()
	log.Printf("Collected: %s", time.Now().UTC())

	bin.FetchDepth(10)
	log.Printf("Fetched: %s", time.Now().UTC())

	log.Print(bin.Asks)
	log.Print(bin.Bids)

}
