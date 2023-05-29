package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
)

func main() {
	// Define RabbitMQ server URL.
	amqpServerURL := os.Getenv("AMQP_SERVER_URL")

	// Create a new RabbitMQ connection.
	connectRabbitMQ, err := amqp.Dial(amqpServerURL)
	if err != nil {
		panic(err)
	}
	defer connectRabbitMQ.Close()

	// Let's start by opening a channel to our RabbitMQ
	// instance over the connection we have already
	// established.
	channelRabbitMQ, err := connectRabbitMQ.Channel()
	if err != nil {
		panic(err)
	}
	defer channelRabbitMQ.Close()

	// With the instance and declare Queues that we can
	// publish and subscribe to.
	_, err = channelRabbitMQ.QueueDeclare(
		"QueueService1", // queue name
		true,            // durable
		false,           // auto delete
		false,           // exclusive
		false,           // no wait
		nil,             // arguments
	)
	if err != nil {
		panic(err)
	}

	// Create a new Gin instance.
	app := gin.Default()

	// Add route for sending/publishing a message to Service 1.
	app.GET("/send", func(c *gin.Context) {
		// Create a message to publish.
		message := amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(c.Query("msg")),
		}

		// Attempt to publish a message to the queue.
		err := channelRabbitMQ.Publish(
			"",              // exchange
			"QueueService1", // queue name
			false,           // mandatory
			false,           // immediate
			message,         // message to publish
		)
		if err != nil {
			c.String(500, err.Error())
			return
		}

		c.String(200, "Message published successfully")
	})

	// Start the Gin server.
	log.Fatal(app.Run(":3000"))
}
