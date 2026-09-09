-- name: CreateKitchenTicket :one
INSERT INTO kitchen_tickets (
    id,
    company_id,
    store_id,
    order_id,
    order_number,
    order_type,
    table_id,
    status,
    priority,
    notes,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
)
RETURNING *;

-- name: CreateKitchenTicketItem :one
INSERT INTO kitchen_ticket_items (
    id,
    ticket_id,
    order_item_id,
    menu_item_id,
    item_name,
    sku,
    quantity,
    notes,
    status,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
)
RETURNING *;

-- name: ListKitchenTicketsByStore :many
SELECT
    kt.id,
    kt.company_id,
    kt.store_id,
    kt.order_id,
    kt.order_number,
    kt.order_type,
    kt.table_id,
    kt.status,
    kt.priority,
    kt.notes,
    kt.created_at,
    kt.updated_at,
    kt.started_at,
    kt.ready_at,
    kt.served_at,
    t.table_number
FROM kitchen_tickets kt
LEFT JOIN tables t ON t.id = kt.table_id
WHERE kt.store_id = $1
    AND (sqlc.narg('status')::varchar IS NULL OR kt.status = sqlc.narg('status'))
ORDER BY
    CASE kt.priority
        WHEN 'VIP' THEN 1
        WHEN 'RUSH' THEN 2
        ELSE 3
    END ASC,
    kt.created_at ASC;

-- name: ListKitchenTicketItemsByTicketIDs :many
SELECT
    id,
    ticket_id,
    order_item_id,
    menu_item_id,
    item_name,
    sku,
    quantity,
    notes,
    status,
    created_at,
    updated_at
FROM kitchen_ticket_items
WHERE ticket_id = ANY($1::uuid[])
ORDER BY created_at ASC;

-- name: GetKitchenTicketByID :one
SELECT
    kt.id,
    kt.company_id,
    kt.store_id,
    kt.order_id,
    kt.order_number,
    kt.order_type,
    kt.table_id,
    kt.status,
    kt.priority,
    kt.notes,
    kt.created_at,
    kt.updated_at,
    kt.started_at,
    kt.ready_at,
    kt.served_at,
    t.table_number
FROM kitchen_tickets kt
LEFT JOIN tables t ON t.id = kt.table_id
WHERE kt.id = $1;

-- name: ListKitchenTicketItemsByTicketID :many
SELECT
    id,
    ticket_id,
    order_item_id,
    menu_item_id,
    item_name,
    sku,
    quantity,
    notes,
    status,
    created_at,
    updated_at
FROM kitchen_ticket_items
WHERE ticket_id = $1
ORDER BY created_at ASC;

-- name: UpdateKitchenTicketStatus :one
UPDATE kitchen_tickets
SET
    status = $2::varchar,
    started_at = CASE WHEN $2::varchar = 'PREPARING' AND started_at IS NULL THEN NOW() ELSE started_at END,
    ready_at = CASE WHEN $2::varchar = 'READY' AND ready_at IS NULL THEN NOW() ELSE ready_at END,
    served_at = CASE WHEN $2::varchar = 'SERVED' AND served_at IS NULL THEN NOW() ELSE served_at END,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateKitchenTicketItemsStatusByTicketID :exec
UPDATE kitchen_ticket_items
SET
    status = $2,
    updated_at = NOW()
WHERE ticket_id = $1;

-- name: GetKitchenTicketByOrderID :one
SELECT
    id,
    company_id,
    store_id,
    order_id,
    order_number,
    order_type,
    table_id,
    status,
    priority,
    notes,
    created_at,
    updated_at,
    started_at,
    ready_at,
    served_at
FROM kitchen_tickets
WHERE order_id = $1
LIMIT 1;
