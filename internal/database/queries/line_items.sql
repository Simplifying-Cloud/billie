-- name: GetLineItem :one
SELECT * FROM line_items WHERE id = ? LIMIT 1;

-- name: ListLineItemsByInvoice :many
SELECT * FROM line_items WHERE invoice_id = ? ORDER BY sort_order;

-- name: CreateLineItem :one
INSERT INTO line_items (id, invoice_id, description, quantity, rate, amount, sort_order)
VALUES (?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateLineItem :one
UPDATE line_items
SET description = ?, quantity = ?, rate = ?, amount = ?, sort_order = ?
WHERE id = ?
RETURNING *;

-- name: DeleteLineItem :exec
DELETE FROM line_items WHERE id = ?;

-- name: DeleteLineItemsByInvoice :exec
DELETE FROM line_items WHERE invoice_id = ?;

-- name: CountLineItemsByInvoice :one
SELECT COUNT(*) FROM line_items WHERE invoice_id = ?;
