package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/lchsk/rss/libs/posts"
)

var perPage int = 20

func handlerFetchInboxPosts(w http.ResponseWriter, req *http.Request) {
	tokenAuth, errToken := ExtractTokenMetadata(req)
	if errToken != nil {
		w.WriteHeader(401)
		return
	}

	page, err := strconv.Atoi(req.URL.Query().Get("page"))

	if err != nil {
		page = 1
	}

	options := posts.FetchPostsOptions{
		FetchPostsMode: posts.FetchPostsModeInbox,
	}

	inboxPosts, err := DBA.Posts.FetchInboxPosts(options, tokenAuth.UserId, page, perPage)

	if err != nil {
		log.Printf("Error fetching inbox posts: %s", err)
		w.WriteHeader(400)
		return
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(inboxPosts)
}
