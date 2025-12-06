-- name: GetClient :one
SELECT * FROM clients WHERE id = ? LIMIT 1;

-- name: GetClientByEmail :one
SELECT * FROM clients WHERE email = ? LIMIT 1;

-- name: ListClients :many
SELECT * FROM clients ORDER BY created_at DESC;

-- name: SearchClients :many
SELECT * FROM clients
WHERE name LIKE ? OR company LIKE ? OR email LIKE ?
ORDER BY name;

-- name: CreateClient :one
INSERT INTO clients (id, name, email, phone, company, address_street, address_city, address_state, address_zip, address_country)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: UpdateClient :one
UPDATE clients
SET name = ?, email = ?, phone = ?, company = ?,
    address_street = ?, address_city = ?, address_state = ?,
    address_zip = ?, address_country = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: DeleteClient :exec
DELETE FROM clients WHERE id = ?;

-- name: CountClients :one
SELECT COUNT(*) FROM clients;
