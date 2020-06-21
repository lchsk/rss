package main

import (
	"encoding/json"
	"github.com/lchsk/rss/libs/user"
	"net/http"

	"github.com/google/uuid"
)

func handlerFetchCurrentUser(w http.ResponseWriter, req *http.Request) {
	tokenAuth, errToken := ExtractTokenMetadata(req)
	if errToken != nil {
		w.WriteHeader(401)
		return
	}

	userId, errFetch := FetchAuth(tokenAuth)

	if errFetch != nil {
		w.WriteHeader(401)
		return
	}

	if _, err := uuid.Parse(userId); err != nil {
		w.WriteHeader(400)
		return
	}

	var err error = nil
	var u *user.User

	if err == nil {
		var dbErr error
		u, dbErr = DBA.User.FindUserById(userId)

		if dbErr != nil {
			w.WriteHeader(404)
			return
		}
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(u)
}
