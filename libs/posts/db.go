package posts

import (
	"database/sql"
	"time"

	"github.com/lchsk/rss/libs/pagination"

	sq "github.com/Masterminds/squirrel"
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
	SQ      *sq.StatementBuilderType
	Queries map[string]*sql.Stmt
}

const (
	FetchPostsModeInbox = iota
	FetchPostsModeChannel
	FetchPostsModeChannels
)

type FetchPostsOptions struct {
	ChannelId      string
	ChannelIds     []string
	FetchPostsMode int
}

type PostData struct {
	Id          string `json:"id"`
	CreatedAt   string `json:"created_at"`
	PubAt       string `json:"pub_at"`
	Title       string `json:"title"`
	Url         string `json:"url"`
	Description string `json:"description"`
	Content     string `json:"content"`
	AuthorName  string `json:"author_name"`
	AuthorEmail string `json:"author_email"`
}

func (ca *PostsAccess) UpdatePostStatusForUser(postId string, userId string, status string) error {
	query := ca.SQ.Update("user_posts").Set("status", status).Where(sq.Eq{"user_id": userId, "post_id": postId})

	_, err := query.RunWith(ca.Db).Exec()

	return err
}

func (ca *PostsAccess) FetchPost(postId string) (*PostData, error) {
	post := &PostData{}

	postQuery := ca.SQ.Select("id, created_at, pub_at, title, url, description, content, author_name, author_email").From("posts").Where(sq.Eq{
		"id": postId,
	}).Limit(1)

	err := postQuery.RunWith(ca.Db).Scan(&post.Id, &post.CreatedAt, &post.PubAt, &post.Title, &post.Url, &post.Description,
		&post.Content, &post.AuthorName, &post.AuthorEmail)

	if err != nil {
		return nil, err
	}

	return post, nil
}

func (ca *PostsAccess) getPostsCount(options FetchPostsOptions, userId string) (int, error) {
	var postsCount int
	var err error

	if options.FetchPostsMode == FetchPostsModeChannel {
		err = ca.Db.QueryRow(sqlFetchUserPostsChannelCount, userId, options.ChannelId).Scan(&postsCount)
	} else if options.FetchPostsMode == FetchPostsModeInbox {
		err = ca.Db.QueryRow(sqlFetchUserPostsInboxCount, userId).Scan(&postsCount)
	} else if options.FetchPostsMode == FetchPostsModeChannels {
		users := ca.SQ.Select("count(p.id)").From("posts p").Join(
			"user_posts up on up.post_id = p.id",
		).Where(sq.Eq{
			"p.channel_id": options.ChannelIds,
			"up.user_id":   userId})
		err = users.RunWith(ca.Db).Scan(&postsCount)
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
	} else if options.FetchPostsMode == FetchPostsModeChannels {
		users := ca.SQ.Select("p.id, p.pub_at, p.title, p.channel_id, up.status").From("posts p").Join(
			"user_posts up on up.post_id = p.id",
		).Where(sq.Eq{
			"p.channel_id": options.ChannelIds,
			"up.user_id":   userId}).OrderBy("p.pub_at ASC").Limit(uint64(paginationValues.Limit)).Offset(uint64(paginationValues.Offset))

		rows, err = users.RunWith(ca.Db).Query()
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

func InitPostsAccess(db *sql.DB, psql *sq.StatementBuilderType) (*PostsAccess, error) {
	queries := map[string]*sql.Stmt{}

	queriesToPrepare := map[string]string{
	}

	for name, sql := range queriesToPrepare {
		stmt, err := db.Prepare(sql)

		if err != nil {
			return nil, err
		}

		queries[name] = stmt
	}

	ca := &PostsAccess{Db: db, SQ: psql, Queries: queries}

	return ca, nil
}
