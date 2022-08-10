package main

import (
	"binance_collector/binance"
	"binance_collector/binance/binance_models"
	"binance_collector/database"
	"binance_collector/queue"
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
)

var collect []interface{}

func cassandra() {
	c := new(database.Cassandra).Init()
	for i := 0; i < 1000; i++ {
		c.WriteData()
	}

	log.Print(c)
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
	if market, err := new(binance.Market).Init().GetDepth("BTCUSDT", 5000); err != nil {
		log.Fatal(err)
	} else {
		data, _ := json.Marshal(market)
		log.Print(string(data))
	}
}

func test2() {
	if market, err := new(binance.Market).Init().GetTicker("24hr"); err != nil {
		log.Fatal(err)
	} else {

		data, _ := json.Marshal(market[0])
		log.Print(string(data))
	}
}
