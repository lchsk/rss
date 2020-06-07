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
  category_id uuid null,

  -- In seconds
  refresh_interval interval not null default '30 minutes',
  last_successful_update timestamp without time zone default (now() at time zone 'utc') not null,

  constraint fk_channels_category_id
     foreign key (category_id)
     references categories (id)
);

drop index if exists idx_channels_channel_url;
create unique index idx_channels_channel_url on channels(channel_url);
create index channels_category_id on channels(category_id);

-- user_channels

drop table if exists user_channels cascade;
create table user_channels (
  id uuid not null primary key,
  created_at timestamp without time zone default (now() at time zone 'utc') not null,

  channel_id uuid not null,
  user_id uuid not null,

  constraint fk_user_channels_channel_id
     foreign key (channel_id)
     references channels (id),
  constraint fk_user_channels_user_id
     foreign key (user_id)
     references users (id)
);

drop index if exists idx_user_channels_channel_id_user_id;
create unique index idx_user_channels_channel_id_user_id on user_channels(channel_id, user_id);

-- articles

drop table if exists articles cascade;
create table articles (
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

  constraint fk_articles_channel_id
     foreign key (channel_id)
     references channels (id)
);
create index idx_articles_channel_id on articles(channel_id);

-- user_articles

drop table if exists user_articles cascade;

drop type if exists user_articles_status;
create type user_articles_status as enum ('unread', 'read');

create table user_articles (
  id uuid not null primary key,
  created_at timestamp without time zone default (now() at time zone 'utc') not null,

  user_id uuid not null,
  article_id uuid not null,

  status user_articles_status not null default 'unread',

  constraint fk_user_articles_user_id
     foreign key (user_id)
     references users (id),
  constraint fk_user_articles_article_id
     foreign key (article_id)
     references articles (id)
);
create index idx_user_articles_user_id on user_articles(user_id);
create index ids_user_articles_article_id on user_articles(article_id);
