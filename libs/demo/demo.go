package demo

import (
	"time"

	"github.com/google/uuid"
	"github.com/lchsk/rss/libs/db"
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

	Post1 string
	Post2 string
	Post3 string
	Post4 string
	Post5 string
	Post6 string
	Post7 string
}

var bugs UserData

var Bugs *UserData = &bugs

func InstallDemo(dba *db.DbAccess) {
	installUsers(dba)
	installCategories(dba)
	installChannels(dba)
	installUserChannels(dba)
	installPosts(dba)
	installUserPosts(dba)
}

func installUsers(dba *db.DbAccess) {
	ua := dba.User

	u, _ := ua.InsertUser("bugs", "bugs@bunny.com", "bunny")

	bugs.UserId = u.ID
}

func installPosts(dba *db.DbAccess) {
	ca := dba.Channel

	now1 := time.Now().UTC()
	now2 := time.Now().UTC().Add(time.Minute)
	now3 := time.Now().UTC().Add(time.Minute * 2)
	now4 := time.Now().UTC().Add(time.Minute * 3)
	now5 := time.Now().UTC().Add(time.Minute * 4)

	bugs.Post1 = uuid.New().String()
	bugs.Post2 = uuid.New().String()
	bugs.Post3 = uuid.New().String()
	bugs.Post4 = uuid.New().String()
	bugs.Post5 = uuid.New().String()
	bugs.Post6 = uuid.New().String()
	bugs.Post7 = uuid.New().String()

	ca.InsertPost(bugs.Post1, &now1, "url", "Post 1", "description", "content",
		"authorName", "authorEmail", bugs.Channel538FeedId,
	)
	ca.InsertPost(bugs.Post2, &now2, "url", "Post 2", "description", "content",
		"authorName", "authorEmail", bugs.Channel538FeedId,
	)
	ca.InsertPost(bugs.Post3, &now3, "url", "Post 3", "description", "content",
		"authorName", "authorEmail", bugs.Channel538FeedId,
	)
	ca.InsertPost(bugs.Post4, &now4, "url", "Post 4", "description", "content",
		"authorName", "authorEmail", bugs.Channel538FeedId,
	)
	ca.InsertPost(bugs.Post5, &now5, "url", "Post 5", "description", "content",
		"authorName", "authorEmail", bugs.Channel538FeedId,
	)

	ca.InsertPost(bugs.Post6, &now4, "url", "Post 6", "description", "content",
		"authorName", "authorEmail", bugs.Channel538NateId,
	)
	ca.InsertPost(bugs.Post7, &now5, "url", "Post 7", "description", "content",
		"authorName", "authorEmail", bugs.Channel538NateId,
	)
}

func installUserPosts(dba *db.DbAccess) {
	ca := dba.Channel

	posts := []string{bugs.Post1, bugs.Post2, bugs.Post3, bugs.Post4, bugs.Post5}
	ca.InsertUserPosts(bugs.Channel538FeedId, posts)

	postsNate := []string{bugs.Post6, bugs.Post7}
	ca.InsertUserPosts(bugs.Channel538NateId, postsNate)

}

func installChannels(dba *db.DbAccess) {
	ca := dba.Channel

	// https://fivethirtyeight.com/features/feed
	channel, _ := ca.InsertChannel("http://localhost:8000/api/debug/channels/538_features")
	bugs.Channel538FeedId = channel.ID

	// https://fivethirtyeight.com/politics/feed
	channel, _ = ca.InsertChannel("http://localhost:8000/api/debug/channels/538_politics")
	bugs.Channel538PoliticsId = channel.ID

	// https://fivethirtyeight.com/contributors/nate-silver/feed/
	channel, _ = ca.InsertChannel("http://localhost:8000/api/debug/channels/538_nate")
	bugs.Channel538NateId = channel.ID

	// https://rss.nytimes.com/services/xml/rss/nyt/Politics.xml
	channel, _ = ca.InsertChannel("http://localhost:8000/api/debug/channels/nyt_politics")
	bugs.ChannelNYTUSPoliticsId = channel.ID

	// https: //rss.nytimes.com/services/xml/rss/nyt/Science.xml
	channel, _ = ca.InsertChannel("http://localhost:8000/api/debug/channels/nyt_science")
	bugs.ChannelNYTScienceId = channel.ID
}

func installUserChannels(dba *db.DbAccess) {
	ca := dba.Channel

	ca.InsertUserChannel(bugs.Channel538FeedId, bugs.UserId, bugs.CategorySportId)
	ca.InsertUserChannel(bugs.Channel538NateId, bugs.UserId, bugs.CategoryBlogsId)
	ca.InsertUserChannel(bugs.ChannelNYTUSPoliticsId, bugs.UserId, bugs.CategoryPoliticsId)
	ca.InsertUserChannel(bugs.ChannelNYTScienceId, bugs.UserId, bugs.CategoryNewsId)
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
