package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/lchsk/rss/db"
	"github.com/lchsk/rss/libs/comms"
)

var queueConn *comms.Connection
var DBA *db.DbAccess

func init() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal(err)
	}

	dba, err := db.GetDBConnection()

	if err != nil {
		log.Fatal(err)
	}

	DBA = dba

	conn, err := comms.ConnectionInit("amqp://guest:guest@localhost:5672/")

	if err != nil {
		log.Fatalf("Cannot connect to rabbit: %s", err)
	}

	queueConn = conn

	if err := comms.DeclareQueues(queueConn.Channel); err != nil {
		log.Fatalf("Could not declare queues: %s\n", err)
	}
}

func waitForMessages() {
	// TODO: Configure logging

	messages, _ := queueConn.Channel.Consume(
		"hello", // queue
		"",      // consumer
		true,    // auto-ack
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // args
	)

	forever := make(chan bool)

	go func() {
		for rawMsg := range messages {
			var msg comms.Message
			err := json.Unmarshal(rawMsg.Body, &msg)
			fmt.Println(err)
			data, _ := json.Marshal(msg.Body)

			if msg.Type == comms.RefreshChannelType {
				go refreshChannelHandler(&msg, data)
			}

			fmt.Println("err", err, msg.Type, msg.Time)
		}
	}()

	log.Printf("Waiting for messages...\n")
	<-forever
}

func main() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		// TODO Add db
		queueConn.ConnectionClose()

		os.Exit(0)
	}()

	waitForMessages()
}
