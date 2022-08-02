package queue

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
	"log"
	"time"
)

type KafkaConnection struct {
	Server string
	Group  string

	Producer *kafka.Producer
	Consumer *kafka.Consumer
}

func (c *KafkaConnection) Init(group string) *KafkaConnection {
	c.Group = group

	c.Producer = nil
	c.Consumer = nil

	return c
}

func (c *KafkaConnection) Delivery(server string) *KafkaConnection {
	c.Server = server

	if producer, err := kafka.NewProducer(
		&kafka.ConfigMap{"bootstrap.servers": c.Server},
	); err != nil {
		log.Fatal(err)
	} else {
		c.Producer = producer

		go func() {
			for e := range c.Producer.Events() {
				switch ev := e.(type) {
				case *kafka.Message:
					if ev.TopicPartition.Error != nil {
						fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
					} else {
						//fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
					}
				}
			}
		}()

	}

	return c
}

func (c *KafkaConnection) Receiver(server string) *KafkaConnection {
	c.Server = server

	if consumer, err := kafka.NewConsumer(
		&kafka.ConfigMap{
			"bootstrap.servers":  c.Server,
			"group.id":           c.Group,
			"auto.offset.reset":  "earliest", //"earliest", "latest",
			"enable.auto.commit": false,
		},
	); err != nil {
		log.Fatal(err)
	} else {
		c.Consumer = consumer
	}

	return c
}

func (c *KafkaConnection) Send(content []byte, topic string) error {
	if c.Producer == nil {
		return errors.New("producer not initialized")
	}

	uid := uuid.New().String()

	if pErr := c.Producer.Produce(
		&kafka.Message{
			TopicPartition: kafka.TopicPartition{
				Topic:     &topic,
				Partition: kafka.PartitionAny,
			},
			Value:     content,
			Key:       []byte(uid),
			Timestamp: time.Now().UTC(),
		}, nil); pErr != nil {
		return pErr
	}

	c.Producer.Flush(15 * 1000)

	return nil
}

func (c *KafkaConnection) SendJSON(content interface{}, topic string) error {
	if c.Producer == nil {
		return errors.New("producer not initialized")
	}

	uid := uuid.New().String()

	if contentBytes, jErr := json.Marshal(content); jErr != nil {
		return jErr
	} else {
		if pErr := c.Producer.Produce(
			&kafka.Message{
				TopicPartition: kafka.TopicPartition{
					Topic:     &topic,
					Partition: kafka.PartitionAny,
				},
				Value:     contentBytes,
				Key:       []byte(uid),
				Timestamp: time.Now().UTC(),
			}, nil); pErr != nil {
			return pErr
		}
	}

	c.Producer.Flush(15 * 1000)

	return nil
}

func (c *KafkaConnection) Receive(topic string, callback func(msg *kafka.Message)) error {
	if c.Consumer == nil {
		return errors.New("consumer not initialized")
	}

	if sErr := c.Consumer.Subscribe(topic, nil); sErr != nil {
		return sErr
	}

	for {
		if msg, rErr := c.Consumer.ReadMessage(-1); rErr != nil {
			return rErr
		} else {
			callback(msg)
		}
	}
}

func (c *KafkaConnection) CloseDelivery() error {
	return c.Consumer.Close()
}

func (c *KafkaConnection) CloseReceiver() {
	c.Producer.Close()
}
