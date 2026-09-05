-- name: ListActiveMenuCategoriesByCompany :many
SELECT mc.id, mc.menu_id, mc.name, mc.sort_order, mc.status
FROM menu_categories mc
JOIN menus m ON mc.menu_id = m.id
WHERE m.company_id = $1
  AND mc.status = 'ACTIVE'
  AND m.status = 'ACTIVE'
ORDER BY mc.sort_order ASC, mc.name ASC;

-- name: ListMenuItemCategoryIDsByCompany :many
SELECT mci.menu_item_id, mci.menu_category_id
FROM menu_category_items mci
JOIN menu_categories mc ON mci.menu_category_id = mc.id
JOIN menus m ON mc.menu_id = m.id
WHERE m.company_id = $1
  AND mc.status = 'ACTIVE'
  AND m.status = 'ACTIVE'
ORDER BY mci.sort_order ASC;
