package comms

import "github.com/streadway/amqp"

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
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")

	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()

	if err != nil {
		return nil, err
	}

	return &Connection{
		Connection: conn,
		Channel:    ch,
	}, nil
}

func (connection *Connection) ConnectionClose() {
	connection.Channel.Close()
	connection.Connection.Close()
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
