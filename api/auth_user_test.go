// +build database

package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

type UserTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func authUser(email string, password string) UserTokens {
	type Input struct {
		Password string
		Email    string
	}

	input := Input{
		Password: password,
		Email:    email,
	}

	inputJson, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", "/authentication", bytes.NewBuffer(inputJson))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlerAuthentication)

	handler.ServeHTTP(rr, req)

	var resp UserTokens
	json.Unmarshal(rr.Body.Bytes(), &resp)

	return resp
}
