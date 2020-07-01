// +build database

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/lchsk/rss/libs/demo"
	"github.com/lchsk/rss/libs/posts"
	"github.com/stretchr/testify/assert"
)

func TestMarkPostAsRead(t *testing.T) {
	setupSchema(DBA.DB)
	demo.InstallDemo(DBA)

	userTokens := authUser("bugs@bunny.com", "bunny")

	type Input struct {
		Status string `json:"status"`
	}
	input := Input{
		Status: "read",
	}

	inputJson, _ := json.Marshal(input)
	req, err := http.NewRequest("PATCH", fmt.Sprintf("/api/posts/%s", demo.Bugs.Post1), bytes.NewBuffer(inputJson))
	req.AddCookie(getCookie("token", userTokens.AccessToken, AccessCookieDuration))

	assert.Nil(t, err)

	rr := httptest.NewRecorder()
	router := getRouter()
	router.ServeHTTP(rr, req)

	assert.Equal(t, 200, rr.Code)

	query := DBA.SQ.Select("status").From("user_posts").Where(sq.Eq{
		"post_id": demo.Bugs.Post1,
		"user_id": demo.Bugs.UserId,
	}).Limit(1)

	var status string
	err = query.RunWith(DBA.DB).Scan(&status)

	assert.Nil(t, err)
	assert.Equal(t, "read", status)
}

func TestFetchPost(t *testing.T) {
	setupSchema(DBA.DB)
	demo.InstallDemo(DBA)

	perPage = 2

	userTokens := authUser("bugs@bunny.com", "bunny")

	req, err := http.NewRequest("GET", fmt.Sprintf("/api/posts/%s", demo.Bugs.Post1), nil)
	req.AddCookie(getCookie("token", userTokens.AccessToken, AccessCookieDuration))

	assert.Nil(t, err)

	rr := httptest.NewRecorder()
	router := getRouter()
	router.ServeHTTP(rr, req)

	assert.Equal(t, 200, rr.Code)

	var resp posts.PostData
	json.Unmarshal(rr.Body.Bytes(), &resp)

	assert.Equal(t, "Post 1", resp.Title)
	assert.Equal(t, "url1", resp.Url)
	assert.Equal(t, "description", resp.Description)
	assert.Equal(t, "content", resp.Content)
	assert.Equal(t, "authorName", resp.AuthorName)
	assert.Equal(t, "authorEmail", resp.AuthorEmail)
}
