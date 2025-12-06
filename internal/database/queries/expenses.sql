-- name: GetExpense :one
SELECT * FROM expenses WHERE id = ? LIMIT 1;

-- name: ListExpenses :many
SELECT * FROM expenses ORDER BY date DESC;

-- name: ListExpensesByCategory :many
SELECT * FROM expenses WHERE category_id = ? ORDER BY date DESC;

-- name: ListExpensesByMerchant :many
SELECT * FROM expenses WHERE merchant_id = ? ORDER BY date DESC;

-- name: ListExpensesByDateRange :many
SELECT * FROM expenses
WHERE date >= ? AND date <= ?
ORDER BY date DESC;

-- name: CreateExpense :one
INSERT INTO expenses (id, description, category_id, merchant_id, amount, date, receipt_path, notes)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateExpense :one
UPDATE expenses
SET description = ?, category_id = ?, merchant_id = ?, amount = ?,
    date = ?, receipt_path = ?, notes = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: UpdateExpenseReceipt :one
UPDATE expenses SET receipt_path = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ? RETURNING *;

-- name: DeleteExpense :exec
DELETE FROM expenses WHERE id = ?;

-- name: GetRecentExpenses :many
SELECT * FROM expenses ORDER BY date DESC LIMIT ?;

-- name: GetTotalExpenses :one
SELECT COALESCE(SUM(amount), 0) as total FROM expenses;

-- name: GetTotalExpensesByCategory :one
SELECT COALESCE(SUM(amount), 0) as total FROM expenses WHERE category_id = ?;

-- name: GetTotalExpensesByDateRange :one
SELECT COALESCE(SUM(amount), 0) as total FROM expenses WHERE date >= ? AND date <= ?;

-- name: CountExpenses :one
SELECT COUNT(*) FROM expenses;
