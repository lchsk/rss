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

func TestRegisterUser(t *testing.T) {
	setupSchema(DBA.DB)

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
	req, err := http.NewRequest("POST", "/api/users", bytes.NewBuffer(inputJson))

	assert.Nil(t, err)

	rr := httptest.NewRecorder()
	router := getRouter()
	router.ServeHTTP(rr, req)

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
