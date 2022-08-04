package main

import (
	"binance_collector/binance"
	"log"
)

func main() {

	// I hope it more fast ....

	bin := new(binance.BinanceSingle)
	bin.Init("BtcUsdT", 5000, 100, true, true)

	if anyErr := bin.Start(10); anyErr != nil {
		log.Print(anyErr)
	}

	bin.CloseWs()

}
