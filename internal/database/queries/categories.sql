-- name: GetCategory :one
SELECT * FROM expense_categories WHERE id = ? LIMIT 1;

-- name: GetCategoryBySlug :one
SELECT * FROM expense_categories WHERE slug = ? LIMIT 1;

-- name: ListCategories :many
SELECT * FROM expense_categories ORDER BY name;

-- name: CreateCategory :one
INSERT INTO expense_categories (id, name, slug)
VALUES (?, ?, ?)
RETURNING *;

-- name: UpdateCategory :one
UPDATE expense_categories SET name = ?, slug = ? WHERE id = ? RETURNING *;

-- name: DeleteCategory :exec
DELETE FROM expense_categories WHERE id = ?;

-- name: CountCategories :one
SELECT COUNT(*) FROM expense_categories;
