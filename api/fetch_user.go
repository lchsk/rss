package main

import (
	"net/http"

	"github.com/google/uuid"
)

func handlerFetchUser(w http.ResponseWriter, req *http.Request) {
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

	if err == nil {
		var dbErr error
		_, dbErr = DBA.User.FindUserById(userId)

		if dbErr != nil {
			w.WriteHeader(404)
			return
		}
	}

	w.WriteHeader(200)
}
