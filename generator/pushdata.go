package main

import (
	"fmt"
	"go_demo/db"
	"go_demo/rabbit"
	"sync"
)

// create main function that will publish a specific number of random messages to the queue
func main() {
	// create a new RabbitMQ connection
	conn, err := rabbit.NewConnection("")
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	// add new work group
	var wg sync.WaitGroup

	counter := 0
	// for 2 workers consume messages from the queue
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// publish 1000 messages to the queue
			for i := 0; i < 100000; i++ {
				message := &db.Message{}
				randomMessage := message.RandomMessage()
				if err := conn.Publish("messages", *randomMessage); err != nil {
					fmt.Println(err)
				}
				counter++
			}
		}()

		// wait for all goroutines to finish
		wg.Wait()

	}

	fmt.Println("Published", counter, "messages to the queue")

}
