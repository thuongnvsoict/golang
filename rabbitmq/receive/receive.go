package main

import (
	"bytes"
	"github.com/streadway/amqp"
	"log"
	"time"
)

// Check return value for each amqp call
func failOnError(err error, msg string){
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main(){
	/*----------------------------
		CONNECT TO RABBITMQ SERVER
	 -----------------------------*/
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	/*-----------------------
		CREATE A CHANNEL
	 ------------------------*/
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	/*-----------------------
		DECLARE A QUEUE
	 ------------------------*/
	q, err := ch.QueueDeclare(
		"task_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	/*-----------------------
		SETTING QOS
	 ------------------------*/
	err = ch.Qos(
		1,
		0,
		false,
	)
	failOnError(err, "Failed to set QoS")

	/*-----------------------
		DEFINE CONSUME
	 ------------------------*/
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	/*----------------------------
		RECEIVE MESSAGE AND SHOW
	 -----------------------------*/
	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			dotCount := bytes.Count(d.Body, []byte("."))
			t := time.Duration(dotCount)
			time.Sleep(t * time.Second)
			log.Printf("Done")
			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}