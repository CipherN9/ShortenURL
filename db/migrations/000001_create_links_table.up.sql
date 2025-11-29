CREATE TABLE IF NOT EXISTS links
(
    id           UUID PRIMARY KEY,
    initial_link TEXT        NOT NULL,
    shorten_link TEXT        NOT NULL UNIQUE,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    modified_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
