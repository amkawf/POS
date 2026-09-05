-- name: ListActiveTablesByStore :many
SELECT *
FROM tables
WHERE company_id = $1
  AND store_id = $2
  AND status != 'INACTIVE'
ORDER BY table_number ASC;

-- name: GetTableByID :one
SELECT *
FROM tables
WHERE company_id = $1
  AND store_id = $2
  AND id = $3;

-- name: UpdateTableStatus :exec
UPDATE tables
SET status = $4, updated_at = NOW()
WHERE company_id = $1
  AND store_id = $2
  AND id = $3;
