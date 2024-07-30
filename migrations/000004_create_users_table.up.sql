CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone not null DEFAULT NOW(),
    name text not null,
    email citext UNIQUE not null,
    password_hash bytea not null,
    activated bool not null,
    version integer not null DEFAULT 1
);