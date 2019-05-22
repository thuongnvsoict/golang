package main

import (
	"github.com/labstack/gommon/log"
	"github.com/streadway/amqp"
	"os"
	"strings"
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
		"task_queue", // name
		true,   // durable
		false,   // delete when usused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	// Send
	body := bodyForm(os.Args)
	err = ch.Publish(
		"",
		q.Name,		// routing key
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType: ("text/plain"),
			Body:[]byte(body),
		})

	log.Printf("[x] Sent %s" , body)
	failOnError(err, "Failed to publish a message")



}

func bodyForm(args []string) string{
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "hello"
	}else{
		s = strings.Join(args[1:], " ")
	}
	return s
}


