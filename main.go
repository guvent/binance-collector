package main

import (
	"binance_collector/binance"
	"log"
)

func main() {
	test3()
}

func test1() {
	if market, err := new(binance.Market).Init().GetDepth("BTCUSDT"); err != nil {
		log.Fatal(err)
	} else {
		log.Print(market)
	}
}

func test2() {
	if market, err := new(binance.Market).Init().GetTicker("24hr"); err != nil {
		log.Fatal(err)
	} else {
		log.Print(market)
	}
}

func test3() {
	if err := new(binance.Stream).Init().Connect(
		[]string{
			"btcusdt@aggTrade",
			"btcusdt@ticker",
			"btcusdt@depth@1000ms",
		},
		func(message interface{}) {
			log.Print(message)
		},
	); err != nil {
		log.Fatal(err)
	}
}
