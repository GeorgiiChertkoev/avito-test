CREATE TYPE pr_status AS ENUM ('OPEN', 'MERGED');

CREATE TABLE team (
    team_name TEXT PRIMARY KEY
);

CREATE TABLE "user" (
    user_id TEXT PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    team_name TEXT NOT NULL REFERENCES team(team_name) ON DELETE CASCADE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE pull_request (
    pull_request_id TEXT PRIMARY KEY,
    pull_request_name TEXT NOT NULL,
    author_id TEXT NOT NULL REFERENCES "user"(user_id) ON DELETE CASCADE,
    status pr_status NOT NULL DEFAULT 'OPEN',
    created_at TIMESTAMP,
    merged_at TIMESTAMP
);

CREATE TABLE pull_request_reviewer (
    pull_request_id TEXT REFERENCES pull_request(pull_request_id) ON DELETE CASCADE,
    reviewer_id TEXT REFERENCES "user"(user_id) ON DELETE CASCADE,
PRIMARY KEY (pull_request_id, reviewer_id)
);

CREATE INDEX idx_user_team_name ON "user"(team_name);
CREATE INDEX idx_pr_author_id ON pull_request(author_id);
CREATE INDEX idx_pr_status ON pull_request(status);
CREATE INDEX idx_pr_reviewer_reviewer_id ON pull_request_reviewer(reviewer_id);
