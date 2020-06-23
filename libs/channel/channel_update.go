package channel

import (
	sq "github.com/Masterminds/squirrel"
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
