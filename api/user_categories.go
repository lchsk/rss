package main

import (
	"encoding/json"
	"net/http"

	"github.com/lchsk/rss/libs/channel"
)

type UserCategoriesResponse struct {
	UserCategories []channel.UserCategory `json:"user_categories"`
}

func handlerFetchCurrentUserCategories(w http.ResponseWriter, req *http.Request) {
	tokenAuth, errToken := ExtractTokenMetadata(req)
	if errToken != nil {
		w.WriteHeader(401)
		return
	}

	var err error = nil
	var userCategories []channel.UserCategory

	if err == nil {
		userCategories, err = DBA.Channel.FetchUserCategories(tokenAuth.UserId)

		if err != nil {
			w.WriteHeader(400)
			return
		}
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(UserCategoriesResponse{
		UserCategories: userCategories,
	})
}
