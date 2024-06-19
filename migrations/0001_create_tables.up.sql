CREATE TABLE "users" (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    first_name VARCHAR,
    last_name VARCHAR,
    username VARCHAR UNIQUE NOT NULL,
    password VARCHAR NOT NULL
);

CREATE TABLE "channels" (
    id integer primary key,
    tg_id integer unique,
    link varchar,
    peer_type varchar,
    username varchar,
    active_usernames jsonb,
    title varchar,
    about varchar,
    category varchar,
    country varchar,
    language varchar,
    image100 varchar,
    image640 varchar,
    participants_count integer,
    tgstat_restrictions jsonb,
    last_updated integer
);

CREATE TABLE "posts" (
    post_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    id bigint,
    date integer,
    views integer,
    link varchar,
    channel_id integer,
    is_deleted integer,
    text varchar,
    media jsonb
);