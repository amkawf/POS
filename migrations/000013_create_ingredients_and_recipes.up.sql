-- 1. Tabel master bahan baku
CREATE TABLE IF NOT EXISTS ingredients (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    code VARCHAR(50),
    name VARCHAR(100) NOT NULL,
    unit VARCHAR(20) NOT NULL,
    min_stock_alert NUMERIC(12, 4) NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_ingredients_company ON ingredients(company_id);

-- 2. Tabel stok fisik bahan baku per toko (Toleran minus untuk Opsi B)
CREATE TABLE IF NOT EXISTS store_ingredient_inventory (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    store_id UUID NOT NULL REFERENCES stores(id) ON DELETE CASCADE,
    ingredient_id UUID NOT NULL REFERENCES ingredients(id) ON DELETE CASCADE,
    stock NUMERIC(12, 4) NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_store_ingredient UNIQUE (store_id, ingredient_id)
);

CREATE INDEX IF NOT EXISTS idx_store_ingredient_inventory_lookup ON store_ingredient_inventory(store_id, ingredient_id);

-- 3. Buku besar mutasi bahan baku (Audit Trail)
CREATE TABLE IF NOT EXISTS ingredient_movements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    store_id UUID NOT NULL REFERENCES stores(id) ON DELETE CASCADE,
    ingredient_id UUID NOT NULL REFERENCES ingredients(id) ON DELETE CASCADE,
    quantity NUMERIC(12, 4) NOT NULL,
    type VARCHAR(30) NOT NULL,
    reference_id UUID,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_ingredient_movements_lookup ON ingredient_movements(store_id, ingredient_id);

-- 4. Tabel resep menu (Bill of Materials / BOM)
CREATE TABLE IF NOT EXISTS recipe_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    menu_item_id UUID NOT NULL REFERENCES menu_items(id) ON DELETE CASCADE,
    ingredient_id UUID NOT NULL REFERENCES ingredients(id) ON DELETE CASCADE,
    quantity_per_portion NUMERIC(12, 4) NOT NULL CHECK (quantity_per_portion > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_recipe_menu_ingredient UNIQUE (menu_item_id, ingredient_id)
);

CREATE INDEX IF NOT EXISTS idx_recipe_items_menu_item ON recipe_items(menu_item_id);

-- 5. Tambahan kolom penanda tipe pemenuhan di menu_items
ALTER TABLE menu_items
ADD COLUMN IF NOT EXISTS fulfillment_type VARCHAR(20) NOT NULL DEFAULT 'BATCH_COOKING';

