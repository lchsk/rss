// +build integration

package main

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/lchsk/rss/libs/channel"
	"github.com/lchsk/rss/libs/comms"
	"github.com/lchsk/rss/libs/db"
	"github.com/lchsk/rss/libs/demo"
	"github.com/stretchr/testify/assert"
)

var DBA *db.DbAccess
var queue *comms.Connection

func setupIntegrationTests() {
	log.Println("Setting up integration tests")

	os.Setenv("INTEGRATION_TEST", "true")
	godotenv.Load("../.env")

	queue, _ = comms.ConnectionInit("amqp://guest:guest@localhost:5672/")

	DBA, _ = db.GetDBConnection()
	db.InstallSchema(DBA.DB, "../sql/schema.sql")
	demo.InstallDemo(DBA)
}

func TestMain(m *testing.M) {
	setupIntegrationTests()

	os.Exit(m.Run())
}

func TestUpdateChannels(t *testing.T) {
	channels, err := DBA.Channel.FetchChannelsToUpdate()

	var cnt int

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

	channel.QueueConn = queue

	updateChannelsTime := time.Now().UTC()
	err = DBA.Channel.UpdateChannels()
	assert.Nil(t, err)

	time.Sleep(100 * time.Millisecond)
	DBA.DB.QueryRow(`select count(id) from articles`).Scan(&cnt)

	assert.Equal(t, 20, cnt)

	var updateTime time.Time
	DBA.DB.QueryRow(`
select last_successful_update from channels
where channel_url = 'http://localhost:8000/api/debug/channels/538_nate'
`).Scan(&updateTime)

	assert.True(t, updateTime.After(updateChannelsTime))

	DBA.DB.QueryRow(`select count(id) from user_articles`).Scan(&cnt)
	assert.Equal(t, 20, cnt)
}
