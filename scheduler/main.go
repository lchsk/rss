package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/lchsk/rss/libs/comms"
	"github.com/lchsk/rss/libs/db"
	"github.com/lchsk/rss/libs/tasktimer"
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
		log.Fatalf("Error connecting to rabbit: %s", err)
	}

	queueConn = conn

	if err := comms.DeclareQueues(queueConn.Channel); err != nil {
		log.Fatalf("Error declaring queues: %s\n", err)
	}
}

func main() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		// TODO: Add db
		queueConn.ConnectionClose()

		os.Exit(0)
	}()
	mgr := tasktimer.TaskManager{}

	// TODO: Improve the interface for adding a task
	mgr.Interval = 1 * time.Second
	mgr.Tasks = append(mgr.Tasks, &tasktimer.Task{
		Every:         400 * time.Millisecond,
		LastExecution: time.Now().UTC(),
		Func: func() {
			err := DBA.Channel.UpdateChannels()

			if err == nil {
				log.Println("UpdateChannels finished successfully")
			} else {
				log.Printf("Error updating channels: %s", err)
			}
		},
	})

	mgr.Wait()
}
