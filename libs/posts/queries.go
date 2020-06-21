package posts

const sqlFetchUserPostsInbox = `
select
	a.id,
	a.pub_at,
	a.title,
	a.channel_id,
	ua.status
from
	posts a
join user_posts ua on
	ua.post_id = a.id
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
	posts a
join user_posts ua on
	ua.post_id = a.id
where
	ua.user_id = $1
`

const sqlFetchUserPostsChannel = `
select
	a.id,
	a.pub_at,
	a.title,
	a.channel_id,
	ua.status
from
	posts a
join user_posts ua on
	ua.post_id = a.id
where
	ua.user_id = $1
	and a.channel_id = $2
order by
	a.pub_at asc
limit $3 offset $4
`

const sqlFetchUserPostsChannelCount = `
select
	count(a.id)
from
	posts a
join user_posts ua on
	ua.post_id = a.id
where
	ua.user_id = $1
	and a.channel_id = $2
`

const sqlFetchUserPostsChannelsCount = `
select
	count(a.id)
from
	posts a
join user_posts ua on
	ua.post_id = a.id
where
	ua.user_id = $1
	and a.channel_id in ($2)
`
