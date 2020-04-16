package channel

import (
	"database/sql"

	"github.com/google/uuid"

	_ "github.com/lib/pq"
)

const sqlInsertChannel = `
    INSERT INTO channels (id, channel_url) VALUES
    ($1, $2)
    RETURNING id
`

type Channel struct {
	ID string `json:"id"`
}

type ChannelAccess struct {
	Queries map[string]*sql.Stmt
}

func (ca *ChannelAccess) InsertChannel(channelUrl string) (*Channel, error) {
	c := &Channel{}

	stmt := ca.Queries["insertChannel"]

	id := uuid.New()

	err := stmt.QueryRow(id, channelUrl).Scan(&c.ID)

	// TODO: Log postgres error

	return c, err
}

func InitChannelAccess(db *sql.DB) (*ChannelAccess, error) {
	queries := map[string]*sql.Stmt{}

	queriesToPrepare := map[string]string{
		"insertChannel": sqlInsertChannel,
	}

	for name, sql := range queriesToPrepare {
		stmt, err := db.Prepare(sql)

		if err != nil {
			return nil, err
		}

		queries[name] = stmt
	}

	ca := &ChannelAccess{Queries: queries}

	return ca, nil
}
