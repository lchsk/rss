// +build database

package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetchUser(t *testing.T) {
	setupSchema(DBA.DB)

	type Input struct {
		Password string
		Email    string
	}

	DBA.User.InsertUser("donaldduck", "donald@duck.com", "donaldduck")

	input := Input{
		Password: "donaldduck",
		Email:    "donald@duck.com",
	}

	inputJson, _ := json.Marshal(input)
	req, err := http.NewRequest("POST", "/api/authentication", bytes.NewBuffer(inputJson))

	assert.Nil(t, err)

	rr := httptest.NewRecorder()
	router := getRouter()
	router.ServeHTTP(rr, req)

	assert.Equal(t, 200, rr.Code)

	type Response struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	var resp Response
	json.Unmarshal(rr.Body.Bytes(), &resp)

	assert.Equal(t, true, len(resp.AccessToken) != 0)
	assert.Equal(t, true, len(resp.RefreshToken) != 0)

	// Fetch user without the access token
	req, err = http.NewRequest("GET", "/api/users/current", nil)
	assert.Nil(t, err)

	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	assert.Equal(t, 401, rr.Code)

	// Fetch user with the access token
	req, err = http.NewRequest("GET", "/api/users/current", nil)
	req.AddCookie(getCookie("token", resp.AccessToken, AccessCookieDuration))
	assert.Nil(t, err)

	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, 200, rr.Code)
}
