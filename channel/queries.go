package channel

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
    id,
	channel_url
from
	channels c
where
	now() at time zone 'utc' - c.last_successful_update >= c.refresh_interval
order by
	last_successful_update asc
limit 1000
`
