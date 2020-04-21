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
	tokenAuth, errToken := ExtractTokenMetadata(req)
	if errToken != nil {
		w.WriteHeader(401)
		return
	}

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
		var insertUserChannel bool = false

		var dbErr error
		channel, dbErr = DBA.Channel.FetchChannelByUrl(input.ChannelUrl)

		if dbErr == nil {
			// Channel exists - we can insert user channel
			insertUserChannel = true
		} else {
			// Channel doesn't exist - need to add it
			channel, dbErr = DBA.Channel.InsertChannel(input.ChannelUrl)

			if dbErr == nil {
				log.Printf("Added new channel URL: %s\n", input.ChannelUrl)
				insertUserChannel = true
			} else {
				log.Printf("Error insertint new channel %s: %s\n", input.ChannelUrl, dbErr)
				err = errors.New(errDbError)
			}
		}

		if insertUserChannel {
			userChannelErr := DBA.Channel.InsertUserChannel(channel.ID, tokenAuth.UserId)

			if userChannelErr != nil {
				log.Printf("Error inserting user channel user_id=%s channel_id=%s : %s\n", channel.ID, tokenAuth.UserId, userChannelErr)
				err = errors.New(errDbError)
			}
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
