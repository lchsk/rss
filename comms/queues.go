package comms

import (
	"log"

	"github.com/streadway/amqp"
)

type Queue struct {
	Name             string
	Durable          bool
	DeleteWhenUnused bool
	Exclusive        bool
	NoWait           bool
}

type Connection struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

func ConnectionInit(url string) (*Connection, error) {
	conn, err := amqp.Dial(url)

	if err != nil {
		return nil, err
	}

	log.Println("Rabbit connection established")

	ch, err := conn.Channel()

	if err != nil {
		return nil, err
	}

	log.Println("Rabbit channel opened")

	return &Connection{
		Connection: conn,
		Channel:    ch,
	}, nil
}

func (connection *Connection) ConnectionClose() {
	log.Printf("Cleaning up connection\n")
	connection.Channel.Close()
	connection.Connection.Close()
	log.Printf("Connection cleaned up\n")
}

func getQueues() []Queue {
	return []Queue{
		Queue{
			Name: "hello",
		},
	}
}

func DeclareQueues(ch *amqp.Channel) error {
	for _, q := range getQueues() {
		_, err := ch.QueueDeclare(
			q.Name,
			q.Durable,
			q.DeleteWhenUnused,
			q.Exclusive,
			q.NoWait,
			nil, // arguments
		)

		if err != nil {
			return err
		}
	}

	return nil
}

func (conn *Connection) Publish(exchange string, routing string, body []byte) error {
	return conn.Channel.Publish(
		exchange,
		routing,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
}
