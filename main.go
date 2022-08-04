package main

import (
	"binance_collector/binance"
	"log"
)

func main() {

	// I hope it more fast ....

	bin := new(binance.BinanceSingle).Init("BtcUsdT", 1000, 100, false, false)

	if anyErr := bin.Any(2); anyErr != nil {
		log.Print(anyErr)
	}

	bin.CloseWS()

}
