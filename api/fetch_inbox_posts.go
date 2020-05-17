package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
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

	inboxPosts, err := DBA.Posts.FetchInboxPosts(tokenAuth.UserId, page, perPage)

	if err != nil {
		log.Printf("Error fetching inbox posts: %s", err)
		w.WriteHeader(400)
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(inboxPosts)
}
