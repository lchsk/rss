package posts

const sqlFetchUserPostsInbox = `
select
	a.id,
	a.pub_at,
	a.title,
	a.channel_id,
	ua.status
from
	articles a
join user_articles ua on
	ua.article_id = a.id
where
	ua.user_id = $1
order by
	a.pub_at asc
limit $2 offset $3
`

const sqlFetchUserPostsInboxCount = `
select
	count(a.id)
from
	articles a
join user_articles ua on
	ua.article_id = a.id
where
	ua.user_id = $1
`
