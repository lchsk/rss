// +build database

package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lchsk/rss/libs/demo"
	"github.com/lchsk/rss/libs/posts"
	"github.com/stretchr/testify/assert"
)

func TestFetchChannelPosts(t *testing.T) {
	setupSchema(DBA.DB)
	demo.InstallDemo(DBA)

	perPage = 2

	userTokens := authUser("bugs@bunny.com", "bunny")

	channelId := demo.Bugs.Channel538NateId
	req, err := http.NewRequest("GET", "/api/posts/channels/"+channelId, nil)
	req.AddCookie(getCookie("token", userTokens.AccessToken, AccessCookieDuration))

	assert.Nil(t, err)

	rr := httptest.NewRecorder()
	router := getRouter()
	router.ServeHTTP(rr, req)

	assert.Equal(t, 200, rr.Code)

	var resp posts.InboxPosts
	json.Unmarshal(rr.Body.Bytes(), &resp)

	assert.Equal(t, 2, len(resp.Posts))

	post1 := resp.Posts[0]
	post2 := resp.Posts[1]
	pagination := resp.Pagination

	assert.Equal(t, "Article 6", post1.Title)
	assert.Equal(t, "Article 7", post2.Title)
	assert.Equal(t, posts.Pagination{
		CurrentPage: 1,
		LastPage:    1,
		Next:        -1,
		Prev:        -1,
	}, pagination)
}