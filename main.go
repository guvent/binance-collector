package main

import (
	"binance_collector/binance"
	"log"
)

func main() {

	// I hope it more fast ....

	bin := new(binance.BinanceSingle).Init("BtcUsdT", 1000, 1000, true, false)

	if anyErr := bin.Any(20); anyErr != nil {
		log.Print(anyErr)
	}

	bin.CloseWS()

}
