package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type PatchUserPostInput struct {
	Status string
}

func handlerPosts(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	postId, ok := vars["id"]

	if !ok {
		w.WriteHeader(400)
		return
	}

	if req.Method == http.MethodGet {
		post, err := DBA.Posts.FetchPost(postId)

		if err != nil {
			log.Printf("Error fetching post %s: %s", postId, err)
			w.WriteHeader(400)
			return
		}

		w.WriteHeader(200)
		json.NewEncoder(w).Encode(post)
	} else if req.Method == http.MethodPatch {
		tokenAuth, errToken := ExtractTokenMetadata(req)
		if errToken != nil {
			w.WriteHeader(401)
			return
		}

		decoder := json.NewDecoder(req.Body)

		var data PatchUserPostInput
		err := decoder.Decode(&data)

		if err != nil {
			w.WriteHeader(400)
			return
		}

		err = DBA.Posts.UpdatePostStatusForUser(postId, tokenAuth.UserId, data.Status)

		if err != nil {
			w.WriteHeader(400)
			return
		}

		w.WriteHeader(200)
	}
}
