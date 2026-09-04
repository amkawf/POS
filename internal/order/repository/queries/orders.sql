-- name: CreateOrder :exec
INSERT INTO orders (
    id,
    company_id,
    store_id,
    order_number,
    order_type,
    order_source,
    status,
    customer_name,
    subtotal,
    discount_amount,
    tax_amount,
    service_amount,
    total_amount,
    notes,
    opened_at,
    completed_at,
    cancelled_at,
    created_by,
    created_at,
    updated_at
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9,
    $10,
    $11,
    $12,
    $13,
    $14,
    $15,
    $16,
    $17,
    $18,
    $19,
    $20
);

-- name: CreateOrderItem :exec
INSERT INTO order_items (
    id,
    order_id,
    menu_item_id,
    item_name,
    sku,
    quantity,
    unit_price,
    modifier_amount,
    discount_amount,
    tax_amount,
    total_amount,
    notes,
    status,
    created_at,
    updated_at
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9,
    $10,
    $11,
    $12,
    $13,
    $14,
    $15
);

-- name: GetOrderByID :one
SELECT *
FROM orders
WHERE company_id = $1
  AND store_id = $2
  AND id = $3;

-- name: ListOrdersByStore :many
SELECT *
FROM orders
WHERE company_id = $1
  AND store_id = $2
ORDER BY opened_at DESC
LIMIT $3;

-- name: ListOrdersByStatus :many
SELECT *
FROM orders
WHERE company_id = $1
  AND store_id = $2
  AND status = $3
ORDER BY opened_at DESC
LIMIT $4;

-- name: ListOrderItemsByOrderID :many
SELECT *
FROM order_items
WHERE order_id = $1
ORDER BY created_at ASC;

-- name: UpdateOrder :exec
UPDATE orders
SET status = $4,
    completed_at = $5,
    cancelled_at = $6,
    notes = COALESCE($7, notes),
    updated_at = $8
WHERE company_id = $1
  AND store_id = $2
  AND id = $3;

-- name: DeleteOrderItemsByOrderID :exec
DELETE FROM order_items
WHERE order_id = $1;

-- name: DeleteOrder :exec
DELETE FROM orders
WHERE company_id = $1
  AND store_id = $2
  AND id = $3;
