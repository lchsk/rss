package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/lchsk/rss/libs/api"
	"github.com/lchsk/rss/libs/user"
)

type AuthenticationInput struct {
	Email    string
	Password string
}
type AuthenticationResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func handlerAuthentication(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)

	var data AuthenticationInput
	err := decoder.Decode(&data)

	if err != nil {
		err = errors.New(errInvalidInputFormat)
	}

	var u *user.User

	if err == nil {
		var dbErr error
		u, dbErr = DBA.User.FindUserByCredentials(data.Email, data.Password)

		if dbErr != nil {
			err = errors.New(errDbError)
		}
	}

	if err != nil {
		w.WriteHeader(401)

		json.NewEncoder(w).Encode(api.ErrorResponse{
			Errors: []api.Error{{ErrorCode: fmt.Sprintf("%s", err)}},
		})
		return
	}

	tokenData, err := CreateToken(u.ID)

	if err != nil {
		w.WriteHeader(422)
		json.NewEncoder(w).Encode(api.ErrorResponse{
			Errors: []api.Error{{ErrorCode: ""}},
		})
		return
	}

	authErr := CreateAuth(u.ID, tokenData)
	if authErr != nil {
		w.WriteHeader(422)
		json.NewEncoder(w).Encode(api.ErrorResponse{
			Errors: []api.Error{{ErrorCode: ""}},
		})
		return
	}

	http.SetCookie(w, getCookie("token", tokenData.AccessToken, AccessCookieDuration))
	http.SetCookie(w, getCookie("refresh", tokenData.RefreshToken, RefreshCookieDuration))

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(AuthenticationResponse{
		AccessToken:  tokenData.AccessToken,
		RefreshToken: tokenData.RefreshToken,
	})
}
