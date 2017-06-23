-- +goose Up
CREATE schema fwatcher;
CREATE TABLE IF NOT EXISTS fwatcher.status (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    filename text NOT NULL,
    current_status text NOT NULL,
    error_string text NOT NULL,
    time_of_processing TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE fwatcher.status;
DROP SCHEMA fwatcher CASCADE;
