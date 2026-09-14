ALTER TABLE menu_items DROP COLUMN IF EXISTS fulfillment_type;
DROP TABLE IF EXISTS recipe_items;
DROP TABLE IF EXISTS ingredient_movements;
DROP TABLE IF EXISTS store_ingredient_inventory;
DROP TABLE IF EXISTS ingredients;

