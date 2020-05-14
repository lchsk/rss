package main

import (
	"encoding/json"
	"net/http"

	"github.com/lchsk/rss/libs/channel"
)

type UserChannelsResponse struct {
	UserChannels []channel.UserChannel `json:"user_channels"`
}

func handlerFetchCurrentUserChannels(w http.ResponseWriter, req *http.Request) {
	tokenAuth, errToken := ExtractTokenMetadata(req)
	if errToken != nil {
		w.WriteHeader(401)
		return
	}

	var err error = nil
	var userChannels []channel.UserChannel

	if err == nil {
		userChannels, err = DBA.Channel.FetchUserChannels(tokenAuth.UserId)

		if err != nil {
			w.WriteHeader(400)
			return
		}
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(UserChannelsResponse{
		UserChannels: userChannels,
	})
}
