-- name: GetMerchant :one
SELECT * FROM merchants WHERE id = ? LIMIT 1;

-- name: GetMerchantByName :one
SELECT * FROM merchants WHERE name = ? LIMIT 1;

-- name: ListMerchants :many
SELECT * FROM merchants ORDER BY name;

-- name: SearchMerchants :many
SELECT * FROM merchants WHERE name LIKE ? ORDER BY name LIMIT 20;

-- name: CreateMerchant :one
INSERT INTO merchants (id, name)
VALUES (?, ?)
RETURNING *;

-- name: UpdateMerchant :one
UPDATE merchants SET name = ? WHERE id = ? RETURNING *;

-- name: DeleteMerchant :exec
DELETE FROM merchants WHERE id = ?;

-- name: CountMerchants :one
SELECT COUNT(*) FROM merchants;
