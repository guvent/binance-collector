package main

import (
	"binance_collector/binance"
	"binance_collector/utils"
	"log"
)

var env = new(utils.GetEnvironment).Init()

func main() {
	runInflux()
}

func runBoth() {
	bin := new(binance.BinanceBoth).Init(env)

	if anyErr := bin.Start(env.CommitMS); anyErr != nil {
		log.Print(anyErr)
	}

	bin.CloseWs()

	log.Println("terminated.")
}

func runInflux() {

	bin := new(binance.Binance).Init(env)

	//bin.QueryTest()

	if anyErr := bin.Start(env.CommitMS); anyErr != nil {
		log.Print(anyErr)
	}

	bin.CloseWs()

	log.Println("terminated.")
}

func runMemory() {
	bin := new(binance.BinanceMemory).Init(env)

	if err := bin.Start(); err != nil {
		log.Fatal(err)
	}

	bin.CloseWs()

	log.Println("terminated.")
}
