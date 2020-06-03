package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func handlerFetchPost(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	postId, ok := vars["id"]

	if !ok {
		w.WriteHeader(400)
		return
	}

	post, err := DBA.Posts.FetchPost(postId)

	if err != nil {
		log.Printf("Error fetching post %s: %s", postId, err)
		w.WriteHeader(400)
		return
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(post)
}
