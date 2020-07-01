package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type ChangeUserPasswordInput struct {
	Current string `json:"current"`
	New     string `json:"new"`
}

func validatePasswordChange(userId string, data ChangeUserPasswordInput) error {
	if len(data.New) < 8 || len(data.New) > 32 {
		return errors.New(errInvalidPasswordLen)
	}

	user, err := DBA.User.FindUserById(userId)

	if err != nil {
		log.Printf("Cannot find user ID=%s", userId)
		return err
	}

	_, err = DBA.User.FindUserByCredentials(user.Email, data.Current)

	if err != nil {
		log.Printf("Cannot find user by credentials")
		return err
	}

	return nil
}

func handlerChangeUserPassword(w http.ResponseWriter, req *http.Request) {
	tokenAuth, errToken := ExtractTokenMetadata(req)
	if errToken != nil {
		w.WriteHeader(401)
		return
	}

	decoder := json.NewDecoder(req.Body)

	var data ChangeUserPasswordInput
	err := decoder.Decode(&data)

	if err != nil {
		w.WriteHeader(400)
		return
	}

	if err := validatePasswordChange(tokenAuth.UserId, data); err != nil {
		log.Printf("Error in validatePasswordChange: %s", err)
		w.WriteHeader(400)
		return
	}

	err = DBA.User.UpdateUserPassword(tokenAuth.UserId, data.New)

	if err != nil {
		log.Printf("Error in UpdateUserPassword: %s", err)
		w.WriteHeader(400)
		return
	}

	w.WriteHeader(200)
}
