-- +goose Up
CREATE TABLE invoice_sequences (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    client_prefix TEXT NOT NULL,
    year INTEGER NOT NULL,
    last_number INTEGER NOT NULL DEFAULT 0,
    UNIQUE(client_prefix, year)
);

-- +goose Down
DROP TABLE invoice_sequences;
