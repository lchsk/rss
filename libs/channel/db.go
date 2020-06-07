package channel

import (
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lchsk/rss/libs/comms"
	"github.com/mmcdole/gofeed"

	_ "github.com/lib/pq"
)

type Channel struct {
	ID string `json:"id"`
}

type UserChannel struct {
	ChannelId     string  `json:"channel_id"`
	ChannelTitle  string  `json:"channel_title"`
	ChannelUrl    string  `json:"channel_url"`
	CategoryId    *string `json:"category_id"`
	CategoryTitle *string `json:"category_title"`

	DbCategoryId    sql.NullString `json:"-"`
	DbCategoryTitle sql.NullString `json:"-"`
}

type ChannelToUpdate struct {
	ChannelId  string
	ChannelUrl string
}

type ChannelAccess struct {
	Db      *sql.DB
	SQ      *sq.StatementBuilderType
	Queries map[string]*sql.Stmt
}

var QueueConn *comms.Connection

func (ca *ChannelAccess) UpdateChannelsDirectly() error {
	channels, err := ca.FetchChannelsToUpdate()

	if err != nil {
		log.Printf("Error in channel update: %s\n", err)
		return err
	}

	for _, channel := range channels {
		fp := gofeed.NewParser()
		feed, err := fp.ParseURL(channel.ChannelUrl)

		if err != nil {
			log.Printf("Error getting channel data for %s: %s\n", channel.ChannelUrl, err)
			return err
		}
		ca.UpdateChannel(channel.ChannelId, feed)
	}

	return nil
}

func (ca *ChannelAccess) UpdateChannels() error {
	channels, err := ca.FetchChannelsToUpdate()

	if err != nil {
		log.Printf("Error in channel update: %s\n", err)
		return err
	}

	for _, channel := range channels {
		refreshMsg := comms.RefreshChannel{Id: channel.ChannelId, Url: channel.ChannelUrl}

		message, err := comms.BuildMessage(refreshMsg)

		if err == nil {
			QueueConn.Publish("", "hello", message)
			log.Printf("Published channel update message for channel id=%s\n", channel.ChannelId)
			return nil
		} else if err != nil {
			log.Printf("Error building channel update message: %s\n", err)
			return err
		}
	}

	return nil
}

func (ca *ChannelAccess) InsertArticle(id string,
	pubAt *time.Time,
	url string,
	title string,
	description string,
	content string,
	authorName string,
	authorEmail string,
	channelId string) error {
	stmt := ca.Queries["insertArticle"]

	_, err := stmt.Exec(id, pubAt, url, title, description, content, authorName, authorEmail, channelId)

	if err != nil {
		log.Printf("Error on InsertArticle: %s", err)
	}

	// TODO: Log postgres error

	return err
}

func (ca *ChannelAccess) UpdateLastSuccessfulUpdateToNow(channelId string, title string, description string, link string, editor string) error {
	values := map[string]interface{}{
		"last_successful_update": time.Now().UTC(),
		"title":                  title,
		"description":            description,
		"website_url":            link,
		"managing_editor":        editor,
	}
	query := ca.SQ.Update("channels").SetMap(
		values,
	).Where(sq.Eq{"id": channelId})

	_, err := query.RunWith(ca.Db).Exec()

	return err
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

func (ca *ChannelAccess) FetchChannelsToUpdate() ([]*ChannelToUpdate, error) {
	stmt := ca.Queries["fetchChannelsToUpdate"]

	rows, err := stmt.Query()
	defer rows.Close()

	if err != nil {
		return nil, err
	}

	var channels []*ChannelToUpdate

	for rows.Next() {
		channel := &ChannelToUpdate{}

		if err := rows.Scan(&channel.ChannelId, &channel.ChannelUrl); err != nil {
			return nil, err
		}

		channels = append(channels, channel)
	}

	return channels, nil
}

func (ca *ChannelAccess) InsertUserArticles(channelId string, articleIds []string) {
	stmt := ca.Queries["fetchChannelUsers"]

	rows, err := stmt.Query(channelId)

	if err != nil {
		log.Printf("Error getting user articles: %s", err)
		return
	}

	defer rows.Close()

	var userId string
	for rows.Next() {
		if err := rows.Scan(&userId); err != nil {
			log.Printf("Could not read userId: %s", err)
			continue
		}

		valueStrings := make([]string, 0, len(articleIds))
		valueArgs := make([]interface{}, 0, len(articleIds)*3)

		for i, articleId := range articleIds {
			valueStrings = append(valueStrings,
				fmt.Sprintf("($%d, $%d, $%d)", i*3+1, i*3+2, i*3+3))
			valueArgs = append(valueArgs, uuid.New().String())
			valueArgs = append(valueArgs, userId)
			valueArgs = append(valueArgs, articleId)
		}

		stmt := fmt.Sprintf(`
		insert into user_articles (id, user_id, article_id) values %s
		`, strings.Join(valueStrings, ","))
		_, err := ca.Db.Exec(stmt, valueArgs...)

		if err != nil {
			log.Printf("Error inserting user articles user_id=%s channel_id=%s: %s", userId, channelId, err)
		}
	}

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
			&uc.ChannelTitle,
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

func (ca *ChannelAccess) UpdateChannel(channelId string, feed *gofeed.Feed) error {
	log.Printf("Updating channel_id=%s", channelId)

	// Load last article published for this channel
	stmt := ca.Queries["fetchLastArticleDate"]

	var date time.Time
	err := stmt.QueryRow(channelId).Scan(&date)

	var minPubTime time.Time

	if err == sql.ErrNoRows {
		minPubTime = time.Unix(0, 0)
	} else {
		minPubTime = date
	}

	var articleIds []string

	for _, item := range feed.Items {
		var pubAt *time.Time

		if item.PublishedParsed == nil {
			current := time.Now().UTC()
			pubAt = &current
		} else {
			pubAt = item.PublishedParsed
		}

		if pubAt.Before(minPubTime) || pubAt.Equal(minPubTime) {
			continue
		}

		authorName := ""
		authorEmail := ""

		if item.Author != nil {
			authorName = item.Author.Name
			authorEmail = item.Author.Email
		}

		// TODO: Escape title, description, content, and other strings
		articleId := uuid.New().String()

		err := ca.InsertArticle(articleId, pubAt,
			item.Link, item.Title, item.Description, item.Content,
			authorName, authorEmail, channelId,
		)

		if err == nil {
			articleIds = append(articleIds, articleId)
		} else {
			log.Printf("Could not insert article to channel id=%s url=%s", channelId, item.Link)
		}
	}

	editor := ""

	if feed.Author != nil {
		editor = feed.Author.Name
	}

	err = ca.UpdateLastSuccessfulUpdateToNow(channelId, feed.Title, feed.Description, feed.Link, editor)

	if err != nil {
		log.Printf("Error updating channel=%s data: %s", channelId, err)
		return err
	}

	if len(articleIds) > 0 {
		ca.InsertUserArticles(channelId, articleIds)
	}

	log.Printf("Channel channel_id=%s updated", channelId)
	return nil
}

func InitChannelAccess(db *sql.DB, psql *sq.StatementBuilderType) (*ChannelAccess, error) {
	queries := map[string]*sql.Stmt{}

	queriesToPrepare := map[string]string{
		"insertChannel":              sqlInsertChannel,
		"insertUserChannel":          sqlInsertUserChannel,
		"insertUserCategory":         sqlInsertUserCategory,
		"insertArticle":              sqlInsertArticle,
		"fetchChannelByUrl":          sqlFetchChannelByUrl,
		"fetchUserChannels":          sqlFetchUserChannels,
		"fetchChannelsToUpdate":      sqlFetchChannelsToUpdate,
		"fetchLastArticleDate":       sqlFetchLastArticleDate,
		"fetchChannelUsers":          sqlFetchChannelUsers,
		"updateLastSuccessfulUpdate": sqlUpdateLastSuccessfulUpdate,
	}

	for name, sql := range queriesToPrepare {
		stmt, err := db.Prepare(sql)

		if err != nil {
			return nil, err
		}

		queries[name] = stmt
	}

	ca := &ChannelAccess{Db: db, SQ: psql, Queries: queries}

	return ca, nil
}
