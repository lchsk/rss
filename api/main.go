package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lchsk/rss/db"
	"github.com/lchsk/rss/libs/api"
	"github.com/lchsk/rss/user"
)

const (
	// URLs
	registerUserUrl = "/api/users"

	// API Errors
	errInvalidUsernameLen = "invalid_username_len"
	errInvalidPasswordLen = "invalid_password_len"
	errInvalidEmailLen    = "invalid_email_len"
	errInvalidEmail       = "invalid_email"
	errInvalidInputFormat = "invalid_input_format"
	// Generic database error, such as violeted index
	errDbError = "db_error"
)

var DBA *db.DbAccess

type UserRegistrationInput struct {
	Username string
	Password string
	Email    string
}

type UserRegistrationResponse struct {
	User *user.User `json:"user"`
}

func validateRegisterUser(input *UserRegistrationInput) error {
	if len(input.Username) < 6 || len(input.Username) > 20 {
		return errors.New(errInvalidUsernameLen)
	}

	if len(input.Password) < 8 || len(input.Password) > 30 {
		return errors.New(errInvalidPasswordLen)
	}

	if len(input.Email) < 6 || len(input.Email) > 20 {
		return errors.New(errInvalidEmailLen)
	}

	if !isEmailValid(input.Email) {
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

func getRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc(registerUserUrl, handlerRegisterUser).Methods(http.MethodPost)

	return router
}

func runAPI() {
	router := getRouter()

	log.Fatal(http.ListenAndServe(":8000", api.CommonMiddleware(router)))
}

func main() {
	conn, _ := db.GetDBConn("rss", "rss", "rss_db", "5432")
	dba, _ := db.InitDbAccess(conn)

	DBA = dba

	runAPI()
}
