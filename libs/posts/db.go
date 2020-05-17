package posts

import (
	"database/sql"
	"time"

	"github.com/lchsk/rss/libs/pagination"

	_ "github.com/lib/pq"
)

type PostBrief struct {
	Id        string    `json:"id"`
	PubAt     time.Time `json:"pub_at"`
	Title     string    `json:"title"`
	ChannelId string    `json:"channel_id"`
	Status    string    `json:"status"`
}

type Pagination struct {
	CurrentPage int `json:"current_page"`
	LastPage    int `json:"last_page"`
	Next        int `json:"next"`
	Prev        int `json:"prev"`
}

type InboxPosts struct {
	Posts      []*PostBrief `json:"posts"`
	Pagination Pagination   `json:"pagination"`
}

type PostsAccess struct {
	Db      *sql.DB
	Queries map[string]*sql.Stmt
}

const (
	FetchPostsModeInbox = iota
	FetchPostsModeChannel
)

type FetchPostsOptions struct {
	ChannelId      string
	FetchPostsMode int
}

func (ca *PostsAccess) getPostsCount(options FetchPostsOptions, userId string) (int, error) {
	var postsCount int
	var err error

	if options.FetchPostsMode == FetchPostsModeChannel {
		err = ca.Db.QueryRow(sqlFetchUserPostsChannelCount, userId, options.ChannelId).Scan(&postsCount)
	} else if options.FetchPostsMode == FetchPostsModeInbox {
		err = ca.Db.QueryRow(sqlFetchUserPostsInboxCount, userId).Scan(&postsCount)
	}

	return postsCount, err
}

func (ca *PostsAccess) getPosts(options FetchPostsOptions, userId string,
	paginationValues *pagination.PaginationValues,
) (*sql.Rows, error) {
	var rows *sql.Rows
	var err error

	if options.FetchPostsMode == FetchPostsModeChannel {
		rows, err = ca.Db.Query(sqlFetchUserPostsChannel, userId, options.ChannelId, paginationValues.Limit, paginationValues.Offset)
	} else if options.FetchPostsMode == FetchPostsModeInbox {
		rows, err = ca.Db.Query(sqlFetchUserPostsInbox, userId, paginationValues.Limit, paginationValues.Offset)
	}

	return rows, err
}

func (ca *PostsAccess) FetchInboxPosts(options FetchPostsOptions, userId string, page int, perPage int) (*InboxPosts, error) {
	postsCount, err := ca.getPostsCount(options, userId)

	if err != nil {
		return nil, err
	}

	paginationValues, err := pagination.GetPaginationValues(page, postsCount, perPage)

	if err != nil {
		return nil, err
	}

	rows, err := ca.getPosts(options, userId, paginationValues)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var posts []*PostBrief

	for rows.Next() {
		post := &PostBrief{}

		if err := rows.Scan(&post.Id, &post.PubAt, &post.Title, &post.ChannelId, &post.Status); err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	lastPage := pagination.GetLastPage(postsCount, perPage)
	prevPage, nextPage := pagination.GetPages(page, lastPage)

	inboxPosts := &InboxPosts{
		Posts: posts,
		Pagination: Pagination{
			CurrentPage: page,
			LastPage:    lastPage,
			Prev:        prevPage,
			Next:        nextPage,
		},
	}

	return inboxPosts, nil
}

func InitPostsAccess(db *sql.DB) (*PostsAccess, error) {
	queries := map[string]*sql.Stmt{}

	queriesToPrepare := map[string]string{
		// "fetchUserPostsInbox":      sqlFetchUserPostsInbox,
		// "fetchUserPostsInboxCount": sqlFetchUserPostsInboxCount,
	}

	for name, sql := range queriesToPrepare {
		stmt, err := db.Prepare(sql)

		if err != nil {
			return nil, err
		}

		queries[name] = stmt
	}

	ca := &PostsAccess{Db: db, Queries: queries}

	return ca, nil
}
