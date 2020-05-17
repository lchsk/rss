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

func (ca *PostsAccess) FetchInboxPosts(userId string, page int, perPage int) (*InboxPosts, error) {
	var postsCount int
	err := ca.Db.QueryRow(sqlFetchUserPostsInboxCount, userId).Scan(&postsCount)

	if err != nil {
		return nil, err
	}

	paginationValues, err := pagination.GetPaginationValues(page, postsCount, perPage)

	if err != nil {
		return nil, err
	}

	rows, err := ca.Db.Query(sqlFetchUserPostsInbox, userId, paginationValues.Limit, paginationValues.Offset)

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
