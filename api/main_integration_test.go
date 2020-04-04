// +build integration

package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	setupIntegrationTests()
	os.Exit(m.Run())
}

func TestRegisterUser__success(t *testing.T) {
	type Input struct {
		Username string
		Password string
		Email    string
	}

	input := Input{
		Username: "donaldduck",
		Password: "donaldduck",
		Email:    "donald@duck.com",
	}

	inputJson, _ := json.Marshal(input)
	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(inputJson))

	assert.Nil(t, err)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlerRegisterUser)
	handler.ServeHTTP(rr, req)

	assert.Equal(t, 201, rr.Code)

	type User struct {
		ID interface{}
	}

	type Response struct {
		User User
	}

	var resp Response
	json.Unmarshal(rr.Body.Bytes(), &resp)

	id, ok := resp.User.ID.(string)
	assert.True(t, ok)
	assert.Equal(t, true, len(id) != 0)
}
