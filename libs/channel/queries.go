package channel

//const sqlInsertUserChannel = `
//    insert into user_channels (id, channel_id, user_id) VALUES
//    ($1, $2, $3)
//`

//const sqlInsertUserCategory = `
//    insert into categories (id, title, user_id, parent_id) VALUES
//    ($1, $2, $3, $4)
//    returning id
//`

//const sqlFetchChannelByUrl = `
//    select id from channels where channel_url = $1
//`

//const sqlFetchUserChannels = `
//select
//	c.ID as channel_id,
//	c.title as channel_title,
//	c.channel_url as channel_url,
//	cat.id as category_id,
//	cat.title as category_title
//from
//	channels c
//join user_channels uc on
//	uc.channel_id = c.id
//left join categories cat on
//	cat.id = c.category_id
//where
//    uc.user_id = $1
//`

//const sqlFetchChannelsToUpdate = `
//select
//    id,
//	channel_url
//from
//	channels c
//where
//	now() at time zone 'utc' - c.last_successful_update >= c.refresh_interval
//order by
//	last_successful_update asc
//limit 1000
//`

//const sqlFetchLastPostDate = `
//select
//    pub_at
//from
//    posts
//where
//    channel_id = $1
//order by
//    pub_at desc
//limit 1
//`

//const sqlFetchChannelUsers = `
//select
//    user_id
//from
//    user_channels
//where
//   channel_id = $1
//`

//const sqlUpdateLastSuccessfulUpdate = `
//update channels
//set last_successful_update = $1
//where id = $2
//`

//const sqlInsertPost = `
//insert into posts (
//    id,
//	pub_at,
//	url,
//	title,
//	description,
//	content,
//	author_name,
//	author_email,
//	channel_id
//)
//values (
//    $1,
//	$2,
//	$3,
//    $4,
//    $5,
//    $6,
//    $7,
//    $8,
//	$9
//)
//`

const SqlFetchChannelsWithinCategoryTree = `
with recursive subcategories as (
select
	id,
	title,
	parent_id
from
	categories c
where
	id = $1
union
select
	c.id,
	c.title,
	c.parent_id
from
	categories c
inner join subcategories s on
	s.id = c.parent_id ) select
	c.id
from
	subcategories s
join channels c on
	c.category_id = s.id
`
