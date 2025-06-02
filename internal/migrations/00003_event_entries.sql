-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS event_types (
    id          BIGSERIAL PRIMARY KEY,
    name        TEXT NOT NULL UNIQUE,
    description TEXT,
    created_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_event_types_name ON event_types(name);

CREATE TABLE IF NOT EXISTS event_entries (
    id            BIGSERIAL PRIMARY KEY,
    event_id      BIGINT NOT NULL REFERENCES events(id) ON DELETE CASCADE,
    event_type_id BIGINT NOT NULL REFERENCES event_types(id) ON DELETE RESTRICT,
    title         TEXT NOT NULL,
    description   TEXT,
    start_time    TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time      TIMESTAMP WITH TIME ZONE NOT NULL,
    location      TEXT,
    created_at    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_event_entries_event_id      ON event_entries(event_id);
CREATE INDEX IF NOT EXISTS idx_event_entries_event_type_id ON event_entries(event_type_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS event_entries;
DROP TABLE IF EXISTS event_types;
-- +goose StatementEnd
