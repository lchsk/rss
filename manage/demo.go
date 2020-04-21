package main

import (
	"github.com/lchsk/rss/db"
)

type UserData struct {
	UserId string

	CategorySportId    string
	CategoryBaseballId string
	CategoryFootballId string
	CategoryBlogsId    string
	CategoryNewsId     string
	CategoryPoliticsId string

	Channel538FeedId       string
	Channel538PoliticsId   string
	Channel538NateId       string
	ChannelNYTUSPoliticsId string
	ChannelNYTScienceId    string
}

var bugs UserData

func installDemo(dba *db.DbAccess) {
	installUsers(dba)
	installCategories(dba)
	installChannels(dba)
	installUserChannels(dba)
}

func installUsers(dba *db.DbAccess) {
	ua := dba.User

	u, _ := ua.InsertUser("bugs", "bugs@bunny.com", "bunny")

	bugs.UserId = u.ID
}

func installChannels(dba *db.DbAccess) {
	ca := dba.Channel

	channel, _ := ca.InsertChannel("https://fivethirtyeight.com/features/feed/", &bugs.CategoryNewsId)
	bugs.Channel538FeedId = channel.ID

	channel, _ = ca.InsertChannel("https://fivethirtyeight.com/politics/feed/", &bugs.CategoryPoliticsId)
	bugs.Channel538PoliticsId = channel.ID

	channel, _ = ca.InsertChannel("https://fivethirtyeight.com/contributors/nate-silver/feed/", &bugs.CategoryBlogsId)
	bugs.Channel538NateId = channel.ID

	channel, _ = ca.InsertChannel("https://rss.nytimes.com/services/xml/rss/nyt/Politics.xml", &bugs.CategoryPoliticsId)
	bugs.ChannelNYTUSPoliticsId = channel.ID

	channel, _ = ca.InsertChannel("https: //rss.nytimes.com/services/xml/rss/nyt/Science.xml", nil)
	bugs.ChannelNYTScienceId = channel.ID
}

func installUserChannels(dba *db.DbAccess) {
	ca := dba.Channel

	ca.InsertUserChannel(bugs.Channel538FeedId, bugs.UserId)
	ca.InsertUserChannel(bugs.Channel538NateId, bugs.UserId)
	ca.InsertUserChannel(bugs.ChannelNYTUSPoliticsId, bugs.UserId)
	ca.InsertUserChannel(bugs.ChannelNYTScienceId, bugs.UserId)
}

func installCategories(dba *db.DbAccess) {
	ca := dba.Channel

	id, err := ca.InsertUserCategory("Sport", bugs.UserId, nil)
	panicIfErr(err)
	bugs.CategorySportId = id.String()

	_, err = ca.InsertUserCategory("Baseball", bugs.UserId, &bugs.CategorySportId)
	panicIfErr(err)

	_, err = ca.InsertUserCategory("Football", bugs.UserId, &bugs.CategorySportId)
	panicIfErr(err)

	id, err = ca.InsertUserCategory("Blogs", bugs.UserId, nil)
	panicIfErr(err)
	bugs.CategoryBlogsId = id.String()

	id, err = ca.InsertUserCategory("News", bugs.UserId, nil)
	panicIfErr(err)
	bugs.CategoryNewsId = id.String()

	id, err = ca.InsertUserCategory("Politics", bugs.UserId, &bugs.CategoryNewsId)
	panicIfErr(err)
	bugs.CategoryPoliticsId = id.String()
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}
