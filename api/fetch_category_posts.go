package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/lchsk/rss/libs/posts"
)

func handlerFetchCategoryPosts(w http.ResponseWriter, req *http.Request) {
	tokenAuth, errToken := ExtractTokenMetadata(req)
	if errToken != nil {
		w.WriteHeader(401)
		return
	}

	vars := mux.Vars(req)
	categoryId, ok := vars["category"]

	if !ok {
		w.WriteHeader(400)
		return
	}

	page, err := strconv.Atoi(req.URL.Query().Get("page"))

	if err != nil {
		page = 1
	}

	rows, err := DBA.DB.Query(SqlFetchChannelsWithinCategoryTree, categoryId)

	if err != nil {
		w.WriteHeader(400)
		return
	}

	var channels []string

	for rows.Next() {
		var channelId string

		if err := rows.Scan(&channelId); err != nil {
		}

		channels = append(channels, channelId)
	}

	options := posts.FetchPostsOptions{
		FetchPostsMode: posts.FetchPostsModeChannels,
		ChannelIds:     channels,
	}

	inboxPosts, err := DBA.Posts.FetchInboxPosts(options, tokenAuth.UserId, page, perPage)

	if err != nil {
		log.Printf("Error fetching category posts: %s", err)
		w.WriteHeader(400)
		return
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(inboxPosts)
}

const SqlFetchChannelsWithinCategoryTree = `
with recursive subcategories as (
select
	id,
	title,
	parent_id
from
	categories c
where
	id = $1
union
select
	c.id,
	c.title,
	c.parent_id
from
	categories c
inner join subcategories s on
	s.id = c.parent_id ) select
	c.id
from
	subcategories s
join channels c on
	c.category_id = s.id
`
