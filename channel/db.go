package channel

import (
	"database/sql"
	"log"

	"github.com/google/uuid"
	"github.com/mmcdole/gofeed"

	_ "github.com/lib/pq"
)

const sqlInsertChannel = `
    insert into channels (id, channel_url, category_id) values
    ($1, $2, $3)
    returning id
`

const sqlInsertUserChannel = `
    insert into user_channels (id, channel_id, user_id) VALUES
    ($1, $2, $3)
`

const sqlInsertUserCategory = `
    insert into categories (id, title, user_id, parent_id) VALUES
    ($1, $2, $3, $4)
    returning id
`

const sqlFetchChannelByUrl = `
    select id from channels where channel_url = $1
`

const sqlFetchUserChannels = `
select
	c.ID as channel_id,
	c.channel_url as channel_url,
	cat.id as category_id,
	cat.title as category_title
from
	channels c
join user_channels uc on
	uc.channel_id = c.id
left join categories cat on
	cat.id = c.category_id
where
    uc.user_id = $1
`

const sqlFetchChannelsToUpdate = `
select
	channel_url
from
	channels c
where
	now() at time zone 'utc' - c.last_successful_update >= c.refresh_interval
order by
	last_successful_update asc
limit 1000
`

type Channel struct {
	ID string `json:"id"`
}

type UserChannel struct {
	ChannelId     string  `json:"channel_id"`
	ChannelUrl    string  `json:"channel_url"`
	CategoryId    *string `json:"category_id"`
	CategoryTitle *string `json:"category_title"`

	DbCategoryId    sql.NullString `json:"-"`
	DbCategoryTitle sql.NullString `json:"-"`
}

type ChannelAccess struct {
	Queries map[string]*sql.Stmt
}

func (ca *ChannelAccess) InsertChannel(channelUrl string, categoryId *string) (*Channel, error) {
	c := &Channel{}

	stmt := ca.Queries["insertChannel"]

	id := uuid.New()

	err := stmt.QueryRow(id, channelUrl, categoryId).Scan(&c.ID)

	// TODO: Log postgres error

	return c, err
}

func (ca *ChannelAccess) InsertUserCategory(title string, userId string, parentId *string) (uuid.UUID, error) {
	stmt := ca.Queries["insertUserCategory"]

	id := uuid.New()

	_, err := stmt.Exec(id, title, userId, parentId)

	return id, err
}

func (ca *ChannelAccess) InsertUserChannel(channelId string, userId string) error {
	stmt := ca.Queries["insertUserChannel"]

	_, err := stmt.Exec(uuid.New(), channelId, userId)

	if err == nil {
		log.Printf("Inserted user channel channel_id=%s user_id=%s\n", channelId, userId)
	} else {
		log.Printf("Error inserting user channel channel_id=%s user_id=%s : %s\n", channelId, userId, err)
	}

	return err
}

func (ca *ChannelAccess) FetchChannelByUrl(channelUrl string) (*Channel, error) {
	c := &Channel{}

	stmt := ca.Queries["fetchChannelByUrl"]

	err := stmt.QueryRow(channelUrl).Scan(&c.ID)

	// TODO: Log postgres error

	return c, err
}

func (ca *ChannelAccess) FetchChannelsToUpdate() ([]string, error) {
	stmt := ca.Queries["fetchChannelsToUpdate"]

	rows, err := stmt.Query()
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	var urls []string

	for rows.Next() {
		var url string

		if err := rows.Scan(&url); err != nil {
			return nil, err
		}

		urls = append(urls, url)
	}

	return urls, nil
}

func (ca *ChannelAccess) FetchUserChannels(userId string) ([]UserChannel, error) {
	userChannels := []UserChannel{}

	stmt := ca.Queries["fetchUserChannels"]

	rows, err := stmt.Query(userId)
	defer rows.Close()

	if err != nil {
		log.Printf("Error fetching user channels for user_id=%s: %s\n", userId, err)
		return nil, err
	}

	for rows.Next() {
		uc := UserChannel{}

		if err := rows.Scan(
			&uc.ChannelId,
			&uc.ChannelUrl,
			&uc.DbCategoryId,
			&uc.DbCategoryTitle,
		); err != nil {
			log.Printf("Error reading user channels for user_id=%s: %s\n", userId, err)
			return nil, err
		}

		if uc.DbCategoryId.Valid {
			uc.CategoryId = &uc.DbCategoryId.String
		}

		if uc.DbCategoryTitle.Valid {
			uc.CategoryTitle = &uc.DbCategoryTitle.String
		}

		userChannels = append(userChannels, uc)
	}

	return userChannels, err
}

func (ca *ChannelAccess) UpdateChannel(feed *gofeed.Feed) error {
	// Load last article published for this channel
	// Filter which articles we need to insert

	// Insert new articles articles

	// Update channel properties in channels
	// Update last successful sync

	return nil
}

func InitChannelAccess(db *sql.DB) (*ChannelAccess, error) {
	queries := map[string]*sql.Stmt{}

	queriesToPrepare := map[string]string{
		"insertChannel":         sqlInsertChannel,
		"insertUserChannel":     sqlInsertUserChannel,
		"insertUserCategory":    sqlInsertUserCategory,
		"fetchChannelByUrl":     sqlFetchChannelByUrl,
		"fetchUserChannels":     sqlFetchUserChannels,
		"fetchChannelsToUpdate": sqlFetchChannelsToUpdate,
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
