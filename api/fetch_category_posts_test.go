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

func TestFetchCategoryPosts(t *testing.T) {
	setupSchema(DBA.DB)
	demo.InstallDemo(DBA)

	perPage = 2

	userTokens := authUser("bugs@bunny.com", "bunny")

	categoryId := demo.Bugs.CategoryNewsId
	req, err := http.NewRequest("GET", "/api/posts/categories/"+categoryId, nil)
	req.AddCookie(getCookie("token", userTokens.AccessToken, AccessCookieDuration))

	assert.Nil(t, err)

	rr := httptest.NewRecorder()
	router := getRouter()
	router.ServeHTTP(rr, req)

	assert.Equal(t, 200, rr.Code)

	var resp posts.InboxPosts
	json.Unmarshal(rr.Body.Bytes(), &resp)

	postsData := resp.Posts
	pagination := resp.Pagination
	assert.Equal(t, 2, len(postsData))

	assert.Equal(t, "Post 1", postsData[0].Title)
	assert.Equal(t, "Post 2", postsData[1].Title)
	assert.Equal(t, posts.Pagination{
		CurrentPage: 1,
		LastPage:    3,
		Next:        2,
		Prev:        -1,
	}, pagination)
}
