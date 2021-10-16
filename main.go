package main

import (
	"encoding/json"
	"log"
	"os"
	"receiver/config"
	"receiver/database"
	"receiver/mq"
)

type ConverterMessage struct {
	VideoLink string `json:"video_link"`
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	config.Init()
	database.Init()
	defer database.Disconnect()

	var currentQueue string
	if len(os.Args) < 2 {
		log.Printf("Usage: %s [info] [warning] [error]", os.Args[0])
		os.Exit(0)
	} else if len(os.Args) == 2 {
		currentQueue = os.Args[1]
	}

	rabbitMQ := config.GetRabbitMQConfig()

	routingKey := currentQueue + rabbitMQ.RoutingKeySuffix

	mq.CreateChannel(currentQueue, routingKey)
	defer mq.ClearConnection()

	msgs := mq.Consume()

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var data ConverterMessage
			json.Unmarshal(d.Body, &data)
			log.Printf("Received a message: %s", data)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
