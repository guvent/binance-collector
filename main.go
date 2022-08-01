package main

import (
	"binance_collector/binance"
	"binance_collector/binance/binance_models"
	"binance_collector/queue"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
	"time"
)

func main() {
	test4()
	test4()
	test4()
	test4()
	test4()
	time.Sleep(1 * time.Second)
	test5()
}

func test5() {
	kf := new(queue.KafkaConnection).Init("binance-2")
	kf.Receiver("localhost")

	if err := kf.Receive("test", func(msg *kafka.Message) {
		message := string(msg.Value)
		log.Print(message)
	}); err != nil {
		log.Fatal(err)
	}
}

func test4() {
	kf := new(queue.KafkaConnection).Init("binance-1")
	kf.Delivery("localhost")

	if err := kf.SendJSON(binance_models.WsRequest{
		Method: "",
		Params: nil,
		Id:     0,
	}, "test"); err != nil {
		log.Fatal(err)
	}

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
