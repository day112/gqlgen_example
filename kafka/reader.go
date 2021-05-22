package main

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
)

func main() {
	fmt.Println("Kafka Consumer using Golang")

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"0.0.0.0:9092"},
		Topic:   "goscreencasts",
		GroupID: "goscreencasts.io",
	})

	defer reader.Close()

	fmt.Println("Consumer is running...")

	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			fmt.Println(err)
		}
		fmt.Printf("Topic : %s Message : %s Partition :%d \n", msg.Topic, string(msg.Value), msg.Partition)

	}
}
