// +build database

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lchsk/rss/libs/demo"
	"github.com/lchsk/rss/libs/posts"
	"github.com/stretchr/testify/assert"
)

func TestFetchPost(t *testing.T) {
	setupSchema(DBA.DB)
	demo.InstallDemo(DBA)

	perPage = 2

	userTokens := authUser("bugs@bunny.com", "bunny")

	req, err := http.NewRequest("GET", fmt.Sprintf("/api/posts/%s", demo.Bugs.Article1), nil)
	req.AddCookie(getCookie("token", userTokens.AccessToken, AccessCookieDuration))

	assert.Nil(t, err)

	rr := httptest.NewRecorder()
	router := getRouter()
	router.ServeHTTP(rr, req)

	assert.Equal(t, 200, rr.Code)

	var resp posts.PostData
	json.Unmarshal(rr.Body.Bytes(), &resp)

	assert.Equal(t, "Article 1", resp.Title)
	assert.Equal(t, "url", resp.Url)
	assert.Equal(t, "description", resp.Description)
	assert.Equal(t, "content", resp.Content)
	assert.Equal(t, "authorName", resp.AuthorName)
	assert.Equal(t, "authorEmail", resp.AuthorEmail)
}
