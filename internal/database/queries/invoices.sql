-- name: GetInvoice :one
SELECT * FROM invoices WHERE id = ? LIMIT 1;

-- name: GetInvoiceByNumber :one
SELECT * FROM invoices WHERE number = ? LIMIT 1;

-- name: GetInvoiceByToken :one
SELECT * FROM invoices WHERE public_token = ? LIMIT 1;

-- name: ListInvoices :many
SELECT * FROM invoices ORDER BY created_at DESC;

-- name: ListInvoicesByStatus :many
SELECT * FROM invoices WHERE status = ? ORDER BY created_at DESC;

-- name: ListInvoicesByClient :many
SELECT * FROM invoices WHERE client_id = ? ORDER BY created_at DESC;

-- name: CreateInvoice :one
INSERT INTO invoices (id, number, client_id, status, issue_date, due_date, subtotal, tax_rate, tax_amount, total, notes, public_token)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateInvoice :one
UPDATE invoices
SET client_id = ?, status = ?, issue_date = ?, due_date = ?,
    subtotal = ?, tax_rate = ?, tax_amount = ?, total = ?,
    notes = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: UpdateInvoiceStatus :one
UPDATE invoices SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ? RETURNING *;

-- name: UpdateInvoiceTotals :one
UPDATE invoices
SET subtotal = ?, tax_amount = ?, total = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: DeleteInvoice :exec
DELETE FROM invoices WHERE id = ?;

-- name: GetRecentInvoices :many
SELECT * FROM invoices ORDER BY created_at DESC LIMIT ?;

-- name: CountInvoices :one
SELECT COUNT(*) FROM invoices;

-- name: CountInvoicesByStatus :one
SELECT COUNT(*) FROM invoices WHERE status = ?;

-- name: GetNextInvoiceSequence :one
INSERT INTO invoice_sequences (client_prefix, year, last_number)
VALUES (?, ?, 1)
ON CONFLICT(client_prefix, year) DO UPDATE SET last_number = last_number + 1
RETURNING last_number;

-- name: MarkOverdueInvoices :exec
UPDATE invoices
SET status = 'overdue', updated_at = CURRENT_TIMESTAMP
WHERE status = 'sent' AND due_date < date('now');
