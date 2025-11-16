-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS teams (
    team_name VARCHAR(255) PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS users (
    user_id VARCHAR(50) PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    team_name VARCHAR(255) NOT NULL REFERENCES teams(team_name),
    is_active BOOLEAN NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_users_team_name ON users(team_name);

CREATE TYPE pr_status AS ENUM ('OPEN', 'MERGED');
CREATE TABLE IF NOT EXISTS pull_requests (
    pull_request_id VARCHAR(50) PRIMARY KEY,
    pull_request_name VARCHAR(255) NOT NULL,
    author_id VARCHAR(50) NOT NULL REFERENCES users(user_id),
    status pr_status NOT NULL,
    created_at  TIMESTAMP DEFAULT now(),
    merged_at  TIMESTAMP
);

CREATE TABLE IF NOT EXISTS reviewers (
    user_id VARCHAR(50) NOT NULL REFERENCES users(user_id),
    pull_request_id VARCHAR(50) NOT NULL REFERENCES pull_requests(pull_request_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS reviewers;
DROP TABLE IF EXISTS pull_requests;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS teams;
DROP INDEX IF EXISTS idx_users_team_name;
DROP TYPE IF EXISTS pr_status;
-- +goose StatementEnd
