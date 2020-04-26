package comms

import (
	"encoding/json"
	"reflect"
	"time"
)

type Message struct {
	Type string
	Time time.Time
	Body interface{}
}

type RefreshChannel struct {
	Url string
}

const RefreshChannelType = "RefreshChannel"

func BuildMessage(body interface{}) ([]byte, error) {
	msgType := reflect.TypeOf(body).Name()

	msg := Message{
		Type: msgType,
		Body: body,
		Time: time.Now().UTC(),
	}

	marshalled, err := json.Marshal(msg)

	if err == nil {
		return marshalled, nil
	}

	return nil, err
}
