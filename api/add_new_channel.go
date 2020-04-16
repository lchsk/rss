package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/lchsk/rss/channel"
	"github.com/lchsk/rss/libs/api"
)

type AddNewChannelInput struct {
	ChannelUrl string `json:"channel_url"`
}

type AddNewChannelResponse struct {
	Channel *channel.Channel `json:"channel"`
}

func isUrlValid(channelUrl string) bool {
	_, err := url.ParseRequestURI(channelUrl)

	return err == nil
}

func handlerAddNewChannelUrl(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)

	var input AddNewChannelInput
	err := decoder.Decode(&input)

	if err != nil {
		log.Println(err)
		err = errors.New(errInvalidInputFormat)
	}

	if err == nil {
		if !isUrlValid(input.ChannelUrl) {
			log.Printf("Invalid channel URL: %s\n", input.ChannelUrl)
			err = errors.New(errInvalidInputFormat)
		}
	}

	var channel *channel.Channel

	if err == nil {
		var dbErr error
		channel, dbErr = DBA.Channel.InsertChannel(input.ChannelUrl)

		if dbErr != nil {
			log.Println(dbErr)
			err = errors.New(errDbError)
		} else {
			log.Printf("Added new channel URL: %s\n", input.ChannelUrl)
		}
	}

	if err != nil {
		w.WriteHeader(400)

		json.NewEncoder(w).Encode(api.ErrorResponse{
			Errors: []api.Error{{ErrorCode: fmt.Sprintf("%s", err)}},
		})
		return
	}

	w.WriteHeader(201)
	json.NewEncoder(w).Encode(AddNewChannelResponse{
		Channel: channel,
	})
}
