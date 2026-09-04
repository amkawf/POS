-- name: ListActiveMenuItemsByCompany :many
SELECT *
FROM menu_items
WHERE company_id = $1
  AND status = 'ACTIVE'
ORDER BY name;
