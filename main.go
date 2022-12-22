package main

import (
	"context"
	"encoding/json"
	"go_demo/db"
	"go_demo/rabbit"
	"go_demo/server"
	"log"
	"sync"
	"time"
)

const (
	// RabbitMQURL is the URL of the RabbitMQ server
	RabbitMQURL = "amqp://guest:guest@localhost:5672/"
	// numWorkers is the number of workers to consume messages
	numWorkers = 24
)

func main() {
	// run db.Migrate() to create the messages table.
	db.Migrate()
	// start the http server in a goroutine
	go server.StartServer()

	// create a new RabbitMQ connection
	conn, err := rabbit.NewConnection(RabbitMQURL)
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	// create a new queue
	queue := "messages"
	if err := conn.CreateQueue(queue); err != nil {
		panic(err)
	}

	counter := 0
	lastCounter := 0

	go consumeMessages(&counter, conn, queue, &sync.WaitGroup{})

	// infinite loop with a sleep of 1 second to display the number of messages consumed
	go func() {
		for {
			consumed := counter - lastCounter
			lastCounter = counter

			log.Println("Consumed", counter, "messages speed: ", consumed, "msg/s")
			time.Sleep(1 * time.Second)
		}
	}()

	// block until context is cancelled
	<-context.Background().Done()
}

func consumeMessages(counter *int, conn *rabbit.Connection, queue string, wg *sync.WaitGroup) error {
	// Create a channel to receive messages
	msgs, err := conn.Consume(queue)
	if err != nil {
		return err
	}

	// consume messages from the channel
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for msg := range msgs {
				// Decode the message body into an Message struct
				var message db.Message
				if err := json.Unmarshal(msg.Body, &message); err != nil {
					log.Println("Error decoding message body: ", err)
					continue
				}

				// if no error ack the message else reject it
				if err != nil {
					msg.Nack(false, true)
					log.Println("Error consuming message: ", err)
					continue
				} else {
					msg.Ack(false)
				}

				*counter++
			}
		}()
	}

	// Wait for all workers to finish

	wg.Wait()

	return nil
}
