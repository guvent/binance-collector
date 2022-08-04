package main

import (
	"binance_collector/binance"
	"log"
)

func main() {

	// I hope it more fast ....

	bin := new(binance.Binance)
	bin.Init("BtcUsdT", 1000, 100, true, true)

	if anyErr := bin.Start(20); anyErr != nil {
		log.Print(anyErr)
	}

	bin.CloseWs()

}
