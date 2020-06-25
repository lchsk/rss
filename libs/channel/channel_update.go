package channel

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/mmcdole/gofeed"
	"log"
)

func (ca *ChannelAccess) RefreshChannel(channelId string) {
	query := ca.SQ.Select("channel_url").From("channels").Where(sq.Eq{
		"id": channelId,
	}).Limit(1)

	var url string
	err := query.RunWith(ca.DbCache).Scan(&url)

	if err != nil {
		log.Printf("Error finding url for channel id: %s", channelId)
		return
	}

	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)

	if err != nil {
		log.Printf("Error getting channel data for %s: %s\n", url, err)
		return
	}
	ca.UpdateChannel(channelId, feed)
}

func (ca *ChannelAccess) AddPostsToUser(channelId string, userId string) {
	query := ca.SQ.Select("id").From("posts").Where(sq.Eq{
		"channel_id": channelId,
	}).OrderBy("created_at DESC")

	rows, err := query.RunWith(ca.DbCache).Query()

	if err != nil {
		log.Printf("Error adding posts from channel_id=%s to user=%s", channelId, userId)
		return
	}

	// TODO: Rewrite adding those posts

	for rows.Next() {
		var postId string
		if err := rows.Scan(&postId); err != nil {
			continue
		}

		query:= ca.SQ.
			Insert("user_posts").Columns("id, user_id, post_id, status").
			Values(uuid.New(), userId, postId, "unread")

		_, err := query.RunWith(ca.DbCache).Exec()

		if err != nil {
			log.Printf("Error adding entry to user_posts, user_id=%s post_id=%s", userId, postId)
		}
	}
}
