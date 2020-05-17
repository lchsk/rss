package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/lchsk/rss/libs/posts"
)

func handlerFetchChannelPosts(w http.ResponseWriter, req *http.Request) {
	tokenAuth, errToken := ExtractTokenMetadata(req)
	if errToken != nil {
		w.WriteHeader(401)
		return
	}

	vars := mux.Vars(req)
	channelId, ok := vars["channel"]

	if !ok {
		w.WriteHeader(400)
		return
	}

	page, err := strconv.Atoi(req.URL.Query().Get("page"))

	if err != nil {
		page = 1
	}

	options := posts.FetchPostsOptions{
		FetchPostsMode: posts.FetchPostsModeChannel,
		ChannelId:      channelId,
	}

	inboxPosts, err := DBA.Posts.FetchInboxPosts(options, tokenAuth.UserId, page, perPage)

	if err != nil {
		log.Printf("Error fetching channel posts: %s", err)
		w.WriteHeader(400)
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(inboxPosts)
}
