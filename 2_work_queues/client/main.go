package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/streadway/amqp"
)

var (
	sleepRetry           = 5 * time.Second
	amqpAddr             = "amqp://guest:guest@rabbitmq:5672/"
	maxSleepTimeDispatch = 5
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	role := flag.String("role", "worker", "Role of the client (processor/dispatcher)")
	flag.Parse()
	log.Printf("Started client with role %s", *role)
	conn, err := amqp.Dial(amqpAddr)
	for err != nil {
		conn, err = amqp.Dial(amqpAddr)
		log.Printf("Cannot connect to %s", amqpAddr)
		time.Sleep(sleepRetry)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"task_queue", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare a queue")

	if *role == "worker" {
		worker(ch, q)
	} else if *role == "task_enqueuer" {
		enqueueTasks(ch, q)
	} else {
		log.Fatalf("Role %s is unrecognized", role)
	}
}

func worker(ch *amqp.Channel, q amqp.Queue) {
	err := ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			dot_count := bytes.Count(d.Body, []byte("."))
			t := time.Duration(dot_count)
			time.Sleep(t * time.Second)
			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func enqueueTasks(ch *amqp.Channel, q amqp.Queue) {
	hostname, err := os.Hostname()
	failOnError(err, "Cannot get hostname")
	body := fmt.Sprintf("Hello from %s", hostname)

	rand.Seed(time.Now().Unix())
	sleepTime := time.Duration(rand.Intn(maxSleepTimeDispatch)+1) * time.Second
	log.Printf("%s will enqueue task every %.f seconds", hostname, sleepTime.Seconds())

	for {
		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,
			amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  "text/plain",
				Body:         []byte(body),
			})
		failOnError(err, "Failed to publish a message")
		time.Sleep(sleepTime)
	}
}
