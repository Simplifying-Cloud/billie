-- +goose Up

-- Clients indexes
CREATE INDEX idx_clients_email ON clients(email);
CREATE INDEX idx_clients_name ON clients(name);
CREATE INDEX idx_clients_company ON clients(company);
CREATE INDEX idx_clients_created_at ON clients(created_at DESC);

-- Invoices indexes
CREATE INDEX idx_invoices_client_id ON invoices(client_id);
CREATE INDEX idx_invoices_status ON invoices(status);
CREATE INDEX idx_invoices_due_date ON invoices(due_date);
CREATE INDEX idx_invoices_public_token ON invoices(public_token);
CREATE INDEX idx_invoices_created_at ON invoices(created_at DESC);
CREATE INDEX idx_invoices_number ON invoices(number);

-- Line items indexes
CREATE INDEX idx_line_items_invoice_id ON line_items(invoice_id);

-- Expenses indexes
CREATE INDEX idx_expenses_category_id ON expenses(category_id);
CREATE INDEX idx_expenses_merchant_id ON expenses(merchant_id);
CREATE INDEX idx_expenses_date ON expenses(date DESC);
CREATE INDEX idx_expenses_created_at ON expenses(created_at DESC);

-- +goose Down
DROP INDEX IF EXISTS idx_clients_email;
DROP INDEX IF EXISTS idx_clients_name;
DROP INDEX IF EXISTS idx_clients_company;
DROP INDEX IF EXISTS idx_clients_created_at;
DROP INDEX IF EXISTS idx_invoices_client_id;
DROP INDEX IF EXISTS idx_invoices_status;
DROP INDEX IF EXISTS idx_invoices_due_date;
DROP INDEX IF EXISTS idx_invoices_public_token;
DROP INDEX IF EXISTS idx_invoices_created_at;
DROP INDEX IF EXISTS idx_invoices_number;
DROP INDEX IF EXISTS idx_line_items_invoice_id;
DROP INDEX IF EXISTS idx_expenses_category_id;
DROP INDEX IF EXISTS idx_expenses_merchant_id;
DROP INDEX IF EXISTS idx_expenses_date;
DROP INDEX IF EXISTS idx_expenses_created_at;
