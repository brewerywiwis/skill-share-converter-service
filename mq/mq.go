package mq

import (
	"fmt"
	"log"
	"skillshare/converter/config"

	"github.com/streadway/amqp"
)

var connection *amqp.Connection
var channel *amqp.Channel
var CurrentQueue string
var RoutingKey string

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func CreateChannel(currentQueue string, routingKey string) *amqp.Channel {
	rabbitMQ := config.GetRabbitMQConfig()
	CurrentQueue = currentQueue
	RoutingKey = routingKey

	connection, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s", rabbitMQ.USERNAME, rabbitMQ.PASSWORD, rabbitMQ.HOST, rabbitMQ.PORT))
	failOnError(err, "Failed to connect to RabbitMQ")
	// defer conn.Close()

	ch, err := connection.Channel()
	channel = ch
	failOnError(err, "Failed to open a channel")
	// defer ch.Close()

	err = channel.ExchangeDeclare(
		rabbitMQ.SensorGatewayExchange, // name
		"direct",                       // type
		true,                           // durable
		false,                          // auto-deleted
		false,                          // internal
		false,                          // no-wait
		nil,                            // arguments
	)
	failOnError(err, "Failed to declare an sensor gateway exchange")

	queue, err := channel.QueueDeclare(
		currentQueue, // name
		false,        // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare a earth queue")

	err = channel.QueueBind(
		queue.Name,                     // queue name
		routingKey,                     // routing key
		rabbitMQ.SensorGatewayExchange, // exchange
		false,
		nil,
	)
	failOnError(err, "Failed to declare an "+currentQueue+" queue")

	return channel
}

func ClearConnection() {
	if connection != nil && !connection.IsClosed() {
		connection.Close()
	}
}
func GetChannel() *amqp.Channel {
	return channel
}

func Publish(data []byte) error {
	err := channel.Publish(
		"sensor_gateway", // exchange
		RoutingKey,       // routing key
		false,            // mandatory
		false,            // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(data),
		})
	failOnError(err, "Failed to publish a message")
	return nil
}
func Consume() <-chan amqp.Delivery {
	msgs, err := channel.Consume(
		CurrentQueue, // queue
		"",           // consumer
		true,         // auto ack
		false,        // exclusive
		false,        // no local
		false,        // no wait
		nil,          // args
	)
	failOnError(err, "Failed to register a consumer")
	return msgs
}
