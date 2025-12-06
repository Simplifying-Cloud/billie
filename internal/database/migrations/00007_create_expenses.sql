-- +goose Up
CREATE TABLE expenses (
    id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
    description TEXT NOT NULL,
    category_id TEXT REFERENCES expense_categories(id) ON DELETE SET NULL,
    merchant_id TEXT REFERENCES merchants(id) ON DELETE SET NULL,
    amount REAL NOT NULL,
    date DATE NOT NULL,
    receipt_path TEXT,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE expenses;
