CREATE TABLE IF NOT EXISTS permissions (
    id bigserial PRIMARY KEY,
    code text NOT NULL
);

CREATE TABLE IF NOT EXISTS users_permissions (
    user_id bigint NOT NULL REFERENCES users ON DELETE CASCADE, -- the primary key of our users
    permission_id bigint NOT NULL REFERENCES permissions ON DELETE CASCADE, -- corresponding pk entry in the permissions table
    PRIMARY KEY (user_id, permission_id)
);

INSERT INTO permissions (code)
VALUES ('movies:read'), ('movies:write');