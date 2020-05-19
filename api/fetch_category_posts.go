package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/lchsk/rss/libs/channel"
	"github.com/lchsk/rss/libs/posts"
)

func handlerFetchCategoryPosts(w http.ResponseWriter, req *http.Request) {
	tokenAuth, errToken := ExtractTokenMetadata(req)
	if errToken != nil {
		w.WriteHeader(401)
		return
	}

	vars := mux.Vars(req)
	categoryId, ok := vars["category"]

	if !ok {
		w.WriteHeader(400)
		return
	}

	page, err := strconv.Atoi(req.URL.Query().Get("page"))

	if err != nil {
		page = 1
	}

	rows, err := DBA.DB.Query(channel.SqlFetchChannelsWithinCategoryTree, categoryId)

	if err != nil {
		w.WriteHeader(400)
		return
	}

	var channels []string

	for rows.Next() {
		var channelId string

		if err := rows.Scan(&channelId); err != nil {
		}

		channels = append(channels, channelId)
	}

	options := posts.FetchPostsOptions{
		FetchPostsMode: posts.FetchPostsModeChannels,
		ChannelIds:     channels,
	}

	inboxPosts, err := DBA.Posts.FetchInboxPosts(options, tokenAuth.UserId, page, perPage)

	if err != nil {
		log.Printf("Error fetching channel posts: %s", err)
		w.WriteHeader(400)
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(inboxPosts)
}
