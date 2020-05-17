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

	Article1 string
	Article2 string
	Article3 string
	Article4 string
	Article5 string
}

var bugs UserData

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

	bugs.Article1 = uuid.New().String()
	bugs.Article2 = uuid.New().String()
	bugs.Article3 = uuid.New().String()
	bugs.Article4 = uuid.New().String()
	bugs.Article5 = uuid.New().String()

	ca.InsertArticle(bugs.Article1, &now1, "url", "Article 1", "description", "content",
		"authorName", "authorEmail", bugs.Channel538FeedId,
	)
	ca.InsertArticle(bugs.Article2, &now2, "url", "Article 2", "description", "content",
		"authorName", "authorEmail", bugs.Channel538FeedId,
	)
	ca.InsertArticle(bugs.Article3, &now3, "url", "Article 3", "description", "content",
		"authorName", "authorEmail", bugs.Channel538FeedId,
	)
	ca.InsertArticle(bugs.Article4, &now4, "url", "Article 4", "description", "content",
		"authorName", "authorEmail", bugs.Channel538FeedId,
	)
	ca.InsertArticle(bugs.Article5, &now5, "url", "Article 5", "description", "content",
		"authorName", "authorEmail", bugs.Channel538FeedId,
	)
}

func installUserPosts(dba *db.DbAccess) {
	ca := dba.Channel

	articles := []string{bugs.Article1, bugs.Article2, bugs.Article3, bugs.Article4, bugs.Article5}

	ca.InsertUserArticles(bugs.Channel538FeedId, articles)
}

func installChannels(dba *db.DbAccess) {
	ca := dba.Channel

	// https://fivethirtyeight.com/features/feed
	channel, _ := ca.InsertChannel("http://localhost:8000/api/debug/channels/538_features", &bugs.CategoryNewsId)
	bugs.Channel538FeedId = channel.ID

	// https://fivethirtyeight.com/politics/feed
	channel, _ = ca.InsertChannel("http://localhost:8000/api/debug/channels/538_politics", &bugs.CategoryPoliticsId)
	bugs.Channel538PoliticsId = channel.ID

	// https://fivethirtyeight.com/contributors/nate-silver/feed/
	channel, _ = ca.InsertChannel("http://localhost:8000/api/debug/channels/538_nate", &bugs.CategoryBlogsId)
	bugs.Channel538NateId = channel.ID

	// https://rss.nytimes.com/services/xml/rss/nyt/Politics.xml
	channel, _ = ca.InsertChannel("http://localhost:8000/api/debug/channels/nyt_politics", &bugs.CategoryPoliticsId)
	bugs.ChannelNYTUSPoliticsId = channel.ID

	// https: //rss.nytimes.com/services/xml/rss/nyt/Science.xml
	channel, _ = ca.InsertChannel("http://localhost:8000/api/debug/channels/nyt_science", nil)
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
