// +build integration

package main

import (
	"testing"
	"time"

	"github.com/lchsk/rss/channel"
	"github.com/lchsk/rss/demo"
	"github.com/lchsk/rss/libs/comms"
	"github.com/stretchr/testify/assert"
)

func TestUpdateChannels__success(t *testing.T) {
	demo.InstallDemo(DBA)
	var cnt int

	channels, err := DBA.Channel.FetchChannelsToUpdate()

	DBA.DB.QueryRow(`select count(id) from articles`).Scan(&cnt)
	assert.Equal(t, cnt, 0)

	assert.Nil(t, err)
	assert.Nil(t, channels)

	DBA.DB.Exec(`
update channels
set last_successful_update = now() at time zone 'utc' - interval '45 minutes'
where channel_url = 'http://localhost:8000/api/debug/channels/538_nate'
`)

	channels, err = DBA.Channel.FetchChannelsToUpdate()

	assert.Nil(t, err)
	assert.Equal(t, len(channels), 1)

	conn, _ := comms.ConnectionInit("amqp://guest:guest@localhost:5672/")
	channel.QueueConn = conn

	updateChannelsTime := time.Now().UTC()
	err = DBA.Channel.UpdateChannels()
	assert.Nil(t, err)

	time.Sleep(100 * time.Millisecond)
	DBA.DB.QueryRow(`select count(id) from articles`).Scan(&cnt)

	assert.Equal(t, cnt, 20)

	var updateTime time.Time
	DBA.DB.QueryRow(`
select last_successful_update from channels
where channel_url = 'http://localhost:8000/api/debug/channels/538_nate'
`).Scan(&updateTime)

	assert.True(t, updateTime.After(updateChannelsTime))

	DBA.DB.QueryRow(`select count(id) from user_articles`).Scan(&cnt)
	assert.Equal(t, 20, cnt)
}
