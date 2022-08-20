package main

import (
	"binance_collector/binance"
	"binance_collector/utils"
	"log"
)

func main() {
	env := new(utils.GetEnvironment).Init()

	bin := new(binance.Binance)

	//bin.Init(env).QueryTest()

	if anyErr := bin.Start(env.CommitMS); anyErr != nil {
		log.Print(anyErr)
	}

	bin.CloseWs()

	log.Println("terminated.")
}
