-- name: GetTotalRevenue :one
SELECT COALESCE(SUM(total), 0) as total
FROM invoices
WHERE status = 'paid';

-- name: GetTotalRevenueByYear :one
SELECT COALESCE(SUM(total), 0) as total
FROM invoices
WHERE status = 'paid' AND strftime('%Y', issue_date) = ?;

-- name: GetOutstandingAmount :one
SELECT COALESCE(SUM(total), 0) as total
FROM invoices
WHERE status IN ('sent', 'overdue');

-- name: GetTotalExpensesByYear :one
SELECT COALESCE(SUM(amount), 0) as total
FROM expenses
WHERE strftime('%Y', date) = ?;

-- name: GetMonthlyRevenue :many
SELECT
    strftime('%m', issue_date) as month,
    COALESCE(SUM(CASE WHEN status = 'paid' THEN total ELSE 0 END), 0) as revenue
FROM invoices
WHERE strftime('%Y', issue_date) = ?
GROUP BY strftime('%m', issue_date)
ORDER BY month;

-- name: GetMonthlyExpenses :many
SELECT
    strftime('%m', date) as month,
    COALESCE(SUM(amount), 0) as expenses
FROM expenses
WHERE strftime('%Y', date) = ?
GROUP BY strftime('%m', date)
ORDER BY month;

-- name: GetInvoiceStatusCounts :many
SELECT status, COUNT(*) as count FROM invoices GROUP BY status;

-- name: GetExpensesByCategory :many
SELECT
    ec.name as category_name,
    ec.slug as category_slug,
    COALESCE(SUM(e.amount), 0) as total
FROM expense_categories ec
LEFT JOIN expenses e ON e.category_id = ec.id
GROUP BY ec.id, ec.name, ec.slug
ORDER BY total DESC;

-- name: GetClientTotalBilled :one
SELECT COALESCE(SUM(total), 0) as total
FROM invoices
WHERE client_id = ? AND status = 'paid';

-- name: GetClientInvoiceCount :one
SELECT COUNT(*) as count FROM invoices WHERE client_id = ?;

-- name: GetUpcomingDueInvoices :many
SELECT * FROM invoices
WHERE status = 'sent' AND due_date <= date('now', '+7 days')
ORDER BY due_date;
