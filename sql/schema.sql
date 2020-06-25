drop extension if exists pgcrypto;
create extension pgcrypto;

-- migrations

drop table if exists migrations cascade;
create table migrations (
  id uuid not null primary key,
  filename text not null,
  created_at timestamp without time zone default (now() at time zone 'utc') not null
);

-- users

drop table if exists users cascade;
create table users (
  id uuid not null primary key,
  username text not null,
  email text not null,
  password text not null,
  created_at timestamp without time zone default (now() at time zone 'utc') not null,
  updated_at timestamp without time zone default (now() at time zone 'utc') not null
);

drop index if exists idx_users_email;
create unique index idx_users_email on users(lower(email));

drop index if exists idx_users_username;
create unique index idx_users_username on users(lower(username));

-- channel_categories

drop table if exists categories cascade;
create table categories (
  id uuid not null primary key,
  created_at timestamp without time zone default (now() at time zone 'utc') not null,
  updated_at timestamp without time zone default (now() at time zone 'utc') not null,

  title text not null,
  user_id uuid not null,
  parent_id uuid null,

  constraint fk_categories_user_id
     foreign key (user_id)
     references users (id),

  constraint fk_categories_parent_id
     foreign key (parent_id)
     references categories (id)
);
create index idx_categories_user_id on categories (user_id);
create index idx_categories_parent_id on categories (parent_id);

-- channels

drop table if exists channels cascade;
create table channels (
  id uuid not null primary key,
  created_at timestamp without time zone default (now() at time zone 'utc') not null,
  updated_at timestamp without time zone default (now() at time zone 'utc') not null,

  channel_url text not null,
  title text not null default '',
  description text not null default '',
  website_url text not null default '',
  managing_editor text not null default '',
  pub_date_str text not null default '',
  pub_date timestamp without time zone null,

  -- In seconds
  refresh_interval interval not null default '30 minutes',
  last_successful_update timestamp without time zone default (now() at time zone 'utc') not null
);

drop index if exists idx_channels_channel_url;
create unique index idx_channels_channel_url on channels(channel_url);

-- user_channels

drop table if exists user_channels cascade;
create table user_channels (
  id uuid not null primary key,
  created_at timestamp without time zone default (now() at time zone 'utc') not null,

  channel_id uuid not null,
  user_id uuid not null,
  category_id uuid not null,

  constraint fk_user_channels_channel_id
     foreign key (channel_id)
     references channels (id),
  constraint fk_user_channels_user_id
     foreign key (user_id)
     references users (id),
  constraint fk_channels_category_id
      foreign key (category_id)
          references categories (id)
);

drop index if exists idx_user_channels_channel_id_user_id;
create unique index idx_user_channels_channel_id_user_id on user_channels(channel_id, user_id);
create unique index idx_user_channels_channel_id_category_id on user_channels(channel_id, category_id);

-- posts

drop table if exists posts cascade;
create table posts (
  id uuid not null primary key,
  created_at timestamp without time zone default (now() at time zone 'utc') not null,

  pub_at timestamp without time zone default (now() at time zone 'utc') not null,

  url text not null,
  title text not null,
  description text not null,
  content text not null,
  author_name text not null,
  author_email text not null,

  channel_id uuid not null,

  constraint fk_posts_channel_id
     foreign key (channel_id)
     references channels (id)
);
create index idx_posts_channel_id on posts(channel_id);

-- user_posts

drop table if exists user_posts cascade;

drop type if exists user_posts_status;
create type user_posts_status as enum ('unread', 'read');

create table user_posts (
  id uuid not null primary key,
  created_at timestamp without time zone default (now() at time zone 'utc') not null,

  user_id uuid not null,
  post_id uuid not null,

  status user_posts_status not null default 'unread',

  constraint fk_user_posts_user_id
     foreign key (user_id)
     references users (id),
  constraint fk_user_posts_post_id
     foreign key (post_id)
     references posts (id)
);
create index idx_user_posts_user_id on user_posts(user_id);
create index ids_user_posts_post_id on user_posts(post_id);
create unique index ids_user_posts_unique_user_post on user_posts(user_id, post_id);
