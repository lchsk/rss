package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/lchsk/rss/libs/api"
	"github.com/lchsk/rss/libs/email"
	"github.com/lchsk/rss/libs/user"
)

type UserRegistrationInput struct {
	Username string
	Password string
	Email    string
}

type UserRegistrationResponse struct {
	User *user.User `json:"user"`
}

func validateRegisterUser(input *UserRegistrationInput) error {
	if len(input.Username) < 6 || len(input.Username) > 32 {
		return errors.New(errInvalidUsernameLen)
	}

	if len(input.Password) < 8 || len(input.Password) > 32 {
		return errors.New(errInvalidPasswordLen)
	}

	if len(input.Email) < 6 || len(input.Email) > 32 {
		return errors.New(errInvalidEmailLen)
	}

	if !email.IsEmailValid(input.Email) {
		return errors.New(errInvalidEmail)
	}

	return nil
}

func handlerRegisterUser(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)

	var data UserRegistrationInput
	err := decoder.Decode(&data)

	if err == nil {
		err = validateRegisterUser(&data)
	} else {
		err = errors.New(errInvalidInputFormat)
	}

	var newUser *user.User

	if err == nil {
		var dbErr error
		newUser, dbErr = DBA.User.InsertUser(data.Username, data.Email, data.Password)

		if dbErr != nil {
			err = errors.New(errDbError)
		}
	}

	if err != nil {
		w.WriteHeader(400)

		json.NewEncoder(w).Encode(api.ErrorResponse{
			Errors: []api.Error{{ErrorCode: fmt.Sprintf("%s", err)}},
		})
		return
	}

	w.WriteHeader(201)
	json.NewEncoder(w).Encode(UserRegistrationResponse{
		User: newUser,
	})
}
