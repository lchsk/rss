// +build integration

package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogout__success(t *testing.T) {
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
	req, err := http.NewRequest("POST", "/authentication", bytes.NewBuffer(inputJson))

	assert.Nil(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlerAuthentication)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, 200, rr.Code)

	type Response struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	var resp Response
	json.Unmarshal(rr.Body.Bytes(), &resp)

	assert.Equal(t, true, len(resp.AccessToken) != 0)
	assert.Equal(t, true, len(resp.RefreshToken) != 0)

	// Logout
	req, err = http.NewRequest("GET", "/logout", nil)
	req.AddCookie(getCookie("token", resp.AccessToken, AccessCookieDuration))

	assert.Nil(t, err)

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(handlerLogout)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, 200, rr.Code)

	// Try to fetch user data
	req, err = http.NewRequest("GET", "/users/1", nil)

	assert.Nil(t, err)

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(handlerFetchUser)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, 401, rr.Code)

	refreshInput := RefreshTokenInput{
		RefreshToken: resp.RefreshToken,
	}

	refreshInputJson, _ := json.Marshal(refreshInput)
	req, err = http.NewRequest("POST", "/authentication", bytes.NewBuffer(refreshInputJson))

	req, err = http.NewRequest("GET", "/authentication/refresh", bytes.NewBuffer(refreshInputJson))

	assert.Nil(t, err)

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Refresh)

	handler.ServeHTTP(rr, req)
	assert.Equal(t, 200, rr.Code)

	json.Unmarshal(rr.Body.Bytes(), &resp)

	assert.Equal(t, true, len(resp.AccessToken) != 0)
	assert.Equal(t, true, len(resp.RefreshToken) != 0)

	// Try to get authed data using new access token
	req, err = http.NewRequest("GET", "/users/1", nil)
	req.AddCookie(getCookie("token", resp.AccessToken, AccessCookieDuration))

	assert.Nil(t, err)

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(handlerFetchUser)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, 200, rr.Code)
}
