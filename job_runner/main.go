package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/lchsk/rss/libs/comms"
	"github.com/lchsk/rss/libs/db"
)

var queueConn *comms.Connection
var DBA *db.DbAccess

type Args struct {
	Task   *string
	Listen *bool

	ChannelIds []string
}

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

func readArgs() *Args {
	args := &Args{}
	args.Task = flag.String("task", "", "Task to perform")
	args.Listen = flag.Bool("listen", false, "Listen to messages")
	var channelIdsStr *string
	channelIdsStr = flag.String("channel_ids", "", "List of channel ids to refresh")
	flag.Parse()

	args.ChannelIds = strings.Split(*channelIdsStr, " ")

	return args
}

func listen() {
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

func runTask(task string, args *Args) {
	if task == "force_refresh_channels" {
		log.Printf("Force refreshing channels: %s", args.ChannelIds)

		for _, channelId := range args.ChannelIds {
			refreshChannel(channelId)
		}
	} else if task == "refresh_channels" {
		DBA.Channel.UpdateChannelsDirectly()
	}
}

func main() {
	args := readArgs()

	if *args.Listen {
		listen()
	}

	if *args.Task != "" {
		runTask(*args.Task, args)
	}
}
