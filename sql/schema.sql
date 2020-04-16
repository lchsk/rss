drop extension if exists pgcrypto;
create extension pgcrypto;

-- users

drop table if exists users cascade;
CREATE TABLE users (
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

-- channels

drop table if exists channels cascade;
create table channels (
  id uuid not null primary key,
  created_at timestamp without time zone default (now() at time zone 'utc') not null,
  updated_at timestamp without time zone default (now() at time zone 'utc') not null,

  title text not null,
  description text not null,
  website_url text not null,
  channel_url text not null,
  managing_editor text not null,
  pub_date_str text not null,
  pub_date timestamp without time zone,
  category_id uuid null,

  constraint fk_channels_category_id
     foreign key (category_id)
     references categories (id)
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

  constraint fk_user_channels_channel_id
     foreign key (channel_id)
     references channels (id),
  constraint fk_user_channels_user_id
     foreign key (user_id)
     references users (id)
);

drop index if exists idx_user_channels_channel_id_user_id;
create unique index idx_user_channels_channel_id_user_id on user_channels(channel_id, user_id);
