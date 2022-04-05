CREATE TABLE IF NOT EXISTS words(
    id serial primary key,
    word text not null default '',
    frequency numeric not null default 0,
    is_primary boolean not null default false
);

CREATE TABLE IF NOT EXISTS auth(
    id serial primary key,
    username text not null,
    pw text not null,
    email text not null,
    created_at TIMESTAMP not null default NOW(),
    updated_at TIMESTAMP not null default NOW()
);