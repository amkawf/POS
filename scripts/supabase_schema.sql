-- =============================================================================
-- SUPABASE SCHEMA & DEV SEED DATA FOR RESTAURANT POS PLATFORM
-- =============================================================================
-- Jalankan seluruh script ini di Supabase Dashboard -> SQL Editor -> Run
-- =============================================================================

CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- 1. COMPANIES
CREATE TABLE IF NOT EXISTS companies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(50) NOT NULL,
    name VARCHAR(150) NOT NULL,
    legal_name VARCHAR(200),
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
    timezone VARCHAR(50) NOT NULL DEFAULT 'Asia/Jakarta',
    currency_code CHAR(3) NOT NULL DEFAULT 'IDR',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_companies_code UNIQUE (code),
    CONSTRAINT chk_companies_status
        CHECK (status IN ('ACTIVE', 'INACTIVE')),
    CONSTRAINT chk_companies_currency_code
        CHECK (currency_code ~ '^[A-Z]{3}$')
);

-- 2. STORES
CREATE TABLE IF NOT EXISTS stores (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(150) NOT NULL,
    address TEXT,
    timezone VARCHAR(50),
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_stores_company
        FOREIGN KEY (company_id)
        REFERENCES companies(id),

    CONSTRAINT uq_stores_company_code
        UNIQUE (company_id, code),

    CONSTRAINT chk_stores_status
        CHECK (status IN ('ACTIVE', 'INACTIVE'))
);

CREATE INDEX IF NOT EXISTS idx_stores_company_id
    ON stores(company_id);

-- 3. MENUS
CREATE TABLE IF NOT EXISTS menus (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    code VARCHAR(50) NOT NULL,
    name VARCHAR(150) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_menus_company
        FOREIGN KEY (company_id)
        REFERENCES companies(id),

    CONSTRAINT uq_menus_company_code
        UNIQUE (company_id, code),

    CONSTRAINT chk_menus_status
        CHECK (status IN ('ACTIVE', 'INACTIVE'))
);

CREATE INDEX IF NOT EXISTS idx_menus_company_id
    ON menus(company_id);

-- 4. MENU CATEGORIES
CREATE TABLE IF NOT EXISTS menu_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    menu_id UUID NOT NULL,
    name VARCHAR(100) NOT NULL,
    sort_order INTEGER NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_menu_categories_menu
        FOREIGN KEY (menu_id)
        REFERENCES menus(id)
        ON DELETE CASCADE,

    CONSTRAINT chk_menu_categories_sort_order
        CHECK (sort_order >= 0),

    CONSTRAINT chk_menu_categories_status
        CHECK (status IN ('ACTIVE', 'INACTIVE'))
);

CREATE INDEX IF NOT EXISTS idx_menu_categories_menu_id
    ON menu_categories(menu_id);

-- 5. MENU ITEMS
CREATE TABLE IF NOT EXISTS menu_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    sku VARCHAR(50) NOT NULL,
    name VARCHAR(150) NOT NULL,
    description TEXT,
    base_price NUMERIC(15, 2) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_menu_items_company
        FOREIGN KEY (company_id)
        REFERENCES companies(id),

    CONSTRAINT uq_menu_items_company_sku
        UNIQUE (company_id, sku),

    CONSTRAINT chk_menu_items_base_price
        CHECK (base_price >= 0),

    CONSTRAINT chk_menu_items_status
        CHECK (status IN ('ACTIVE', 'INACTIVE'))
);

CREATE INDEX IF NOT EXISTS idx_menu_items_company_id
    ON menu_items(company_id);

-- 6. STORE MENU ITEMS
CREATE TABLE IF NOT EXISTS store_menu_items (
    store_id UUID NOT NULL,
    menu_item_id UUID NOT NULL,
    price_override NUMERIC(15, 2),
    is_available BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT pk_store_menu_items
        PRIMARY KEY (store_id, menu_item_id),

    CONSTRAINT fk_store_menu_items_store
        FOREIGN KEY (store_id)
        REFERENCES stores(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_store_menu_items_menu_item
        FOREIGN KEY (menu_item_id)
        REFERENCES menu_items(id)
        ON DELETE CASCADE,

    CONSTRAINT chk_store_menu_items_price_override
        CHECK (price_override IS NULL OR price_override >= 0)
);

CREATE INDEX IF NOT EXISTS idx_store_menu_items_menu_item_id
    ON store_menu_items(menu_item_id);

-- 7. MENU CATEGORY ITEMS
CREATE TABLE IF NOT EXISTS menu_category_items (
    menu_category_id UUID NOT NULL,
    menu_item_id UUID NOT NULL,
    sort_order INTEGER NOT NULL DEFAULT 0,

    CONSTRAINT pk_menu_category_items
        PRIMARY KEY (menu_category_id, menu_item_id),

    CONSTRAINT fk_menu_category_items_category
        FOREIGN KEY (menu_category_id)
        REFERENCES menu_categories(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_menu_category_items_menu_item
        FOREIGN KEY (menu_item_id)
        REFERENCES menu_items(id)
        ON DELETE CASCADE,

    CONSTRAINT chk_menu_category_items_sort_order
        CHECK (sort_order >= 0)
);

CREATE INDEX IF NOT EXISTS idx_menu_category_items_menu_item_id
    ON menu_category_items(menu_item_id);

-- 8. ORDERS
CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    store_id UUID NOT NULL,
    table_id UUID,
    customer_session_id UUID,
    order_number VARCHAR(50) NOT NULL,
    order_type VARCHAR(20) NOT NULL,
    order_source VARCHAR(20) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'DRAFT',
    customer_name VARCHAR(150),
    subtotal NUMERIC(15, 2) NOT NULL DEFAULT 0,
    discount_amount NUMERIC(15, 2) NOT NULL DEFAULT 0,
    tax_amount NUMERIC(15, 2) NOT NULL DEFAULT 0,
    service_amount NUMERIC(15, 2) NOT NULL DEFAULT 0,
    total_amount NUMERIC(15, 2) NOT NULL DEFAULT 0,
    notes TEXT,
    opened_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMPTZ,
    cancelled_at TIMESTAMPTZ,
    created_by UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_orders_company
        FOREIGN KEY (company_id)
        REFERENCES companies(id),

    CONSTRAINT fk_orders_store
        FOREIGN KEY (store_id)
        REFERENCES stores(id),

    CONSTRAINT uq_orders_store_order_number
        UNIQUE (store_id, order_number),

    CONSTRAINT chk_orders_order_type
        CHECK (order_type IN ('DINE_IN', 'TAKEAWAY', 'DELIVERY')),

    CONSTRAINT chk_orders_order_source
        CHECK (order_source IN ('POS', 'WAITER', 'QR')),

    CONSTRAINT chk_orders_status
        CHECK (
            status IN (
                'DRAFT',
                'OPEN',
                'COMPLETED',
                'CANCELLED',
                'VOID',
                'REFUNDED'
            )
        ),

    CONSTRAINT chk_orders_amounts
        CHECK (
            subtotal >= 0
            AND discount_amount >= 0
            AND tax_amount >= 0
            AND service_amount >= 0
            AND total_amount >= 0
        )
);

CREATE INDEX IF NOT EXISTS idx_orders_company_id ON orders(company_id);
CREATE INDEX IF NOT EXISTS idx_orders_store_id ON orders(store_id);
CREATE INDEX IF NOT EXISTS idx_orders_store_status ON orders(store_id, status);
CREATE INDEX IF NOT EXISTS idx_orders_store_opened_at ON orders(store_id, opened_at);

-- 9. ORDER ITEMS
CREATE TABLE IF NOT EXISTS order_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL,
    menu_item_id UUID NOT NULL,
    item_name VARCHAR(150) NOT NULL,
    sku VARCHAR(50) NOT NULL,
    quantity NUMERIC(15, 3) NOT NULL,
    unit_price NUMERIC(15, 2) NOT NULL,
    modifier_amount NUMERIC(15, 2) NOT NULL DEFAULT 0,
    discount_amount NUMERIC(15, 2) NOT NULL DEFAULT 0,
    tax_amount NUMERIC(15, 2) NOT NULL DEFAULT 0,
    total_amount NUMERIC(15, 2) NOT NULL DEFAULT 0,
    notes TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_order_items_order
        FOREIGN KEY (order_id)
        REFERENCES orders(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_order_items_menu_item
        FOREIGN KEY (menu_item_id)
        REFERENCES menu_items(id),

    CONSTRAINT chk_order_items_quantity
        CHECK (quantity > 0),

    CONSTRAINT chk_order_items_amounts
        CHECK (
            unit_price >= 0
            AND modifier_amount >= 0
            AND discount_amount >= 0
            AND tax_amount >= 0
            AND total_amount >= 0
        ),

    CONSTRAINT chk_order_items_status
        CHECK (status IN ('ACTIVE', 'VOID'))
);

CREATE INDEX IF NOT EXISTS idx_order_items_order_id ON order_items(order_id);
CREATE INDEX IF NOT EXISTS idx_order_items_menu_item_id ON order_items(menu_item_id);

-- 10. TABLES
CREATE TABLE IF NOT EXISTS tables (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    store_id UUID NOT NULL,
    table_number VARCHAR(50) NOT NULL,
    capacity INTEGER NOT NULL DEFAULT 4,
    status VARCHAR(20) NOT NULL DEFAULT 'AVAILABLE',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_tables_company
        FOREIGN KEY (company_id)
        REFERENCES companies(id),

    CONSTRAINT fk_tables_store
        FOREIGN KEY (store_id)
        REFERENCES stores(id),

    CONSTRAINT uq_tables_store_number
        UNIQUE (store_id, table_number),

    CONSTRAINT chk_tables_capacity
        CHECK (capacity > 0),

    CONSTRAINT chk_tables_status
        CHECK (status IN ('AVAILABLE', 'OCCUPIED', 'RESERVED', 'INACTIVE'))
);

CREATE INDEX IF NOT EXISTS idx_tables_store_id ON tables(store_id);
CREATE INDEX IF NOT EXISTS idx_tables_store_status ON tables(store_id, status);
CREATE INDEX IF NOT EXISTS idx_tables_company_id ON tables(company_id);

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_constraint WHERE conname = 'fk_orders_table'
    ) THEN
        ALTER TABLE orders
            ADD CONSTRAINT fk_orders_table
            FOREIGN KEY (table_id)
            REFERENCES tables(id)
            ON DELETE SET NULL;
    END IF;
END $$;

-- 11. KITCHEN TICKETS & ITEMS
CREATE TABLE IF NOT EXISTS kitchen_tickets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL,
    store_id UUID NOT NULL,
    order_id UUID NOT NULL,
    order_number VARCHAR(50) NOT NULL,
    order_type VARCHAR(20) NOT NULL,
    table_id UUID,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    priority VARCHAR(20) NOT NULL DEFAULT 'NORMAL',
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    started_at TIMESTAMPTZ,
    ready_at TIMESTAMPTZ,
    served_at TIMESTAMPTZ,

    CONSTRAINT fk_kitchen_tickets_company
        FOREIGN KEY (company_id)
        REFERENCES companies(id),

    CONSTRAINT fk_kitchen_tickets_store
        FOREIGN KEY (store_id)
        REFERENCES stores(id),

    CONSTRAINT fk_kitchen_tickets_order
        FOREIGN KEY (order_id)
        REFERENCES orders(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_kitchen_tickets_table
        FOREIGN KEY (table_id)
        REFERENCES tables(id)
        ON DELETE SET NULL,

    CONSTRAINT chk_kitchen_tickets_status
        CHECK (status IN ('PENDING', 'PREPARING', 'READY', 'SERVED', 'CANCELLED')),

    CONSTRAINT chk_kitchen_tickets_priority
        CHECK (priority IN ('NORMAL', 'RUSH', 'VIP'))
);

CREATE TABLE IF NOT EXISTS kitchen_ticket_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ticket_id UUID NOT NULL,
    order_item_id UUID,
    menu_item_id UUID NOT NULL,
    item_name VARCHAR(150) NOT NULL,
    sku VARCHAR(50) NOT NULL,
    quantity NUMERIC(15, 3) NOT NULL,
    notes TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_kitchen_ticket_items_ticket
        FOREIGN KEY (ticket_id)
        REFERENCES kitchen_tickets(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_kitchen_ticket_items_menu_item
        FOREIGN KEY (menu_item_id)
        REFERENCES menu_items(id),

    CONSTRAINT chk_kitchen_ticket_items_quantity
        CHECK (quantity > 0),

    CONSTRAINT chk_kitchen_ticket_items_status
        CHECK (status IN ('PENDING', 'PREPARING', 'READY', 'SERVED', 'CANCELLED'))
);

CREATE INDEX IF NOT EXISTS idx_kitchen_tickets_store_status ON kitchen_tickets(store_id, status);
CREATE INDEX IF NOT EXISTS idx_kitchen_tickets_order_id ON kitchen_tickets(order_id);
CREATE INDEX IF NOT EXISTS idx_kitchen_tickets_created_at ON kitchen_tickets(store_id, created_at);
CREATE INDEX IF NOT EXISTS idx_kitchen_ticket_items_ticket_id ON kitchen_ticket_items(ticket_id);

-- =============================================================================
-- SEED DATA (DEV COMPANY, STORE, MENU, ITEMS & TABLES)
-- =============================================================================
INSERT INTO companies (id, code, name)
VALUES ('11111111-1111-1111-1111-111111111111', 'C1', 'Dev Company')
ON CONFLICT (id) DO NOTHING;

INSERT INTO stores (id, company_id, code, name)
VALUES ('22222222-2222-2222-2222-222222222222', '11111111-1111-1111-1111-111111111111', 'S1', 'Dev Store')
ON CONFLICT (id) DO NOTHING;

INSERT INTO menus (id, company_id, code, name)
VALUES ('66666666-6666-6666-6666-666666666666', '11111111-1111-1111-1111-111111111111', 'MAIN', 'Main Menu')
ON CONFLICT (id) DO NOTHING;

INSERT INTO menu_categories (id, menu_id, name, sort_order)
VALUES
    ('77777777-7777-7777-7777-777777777771', '66666666-6666-6666-6666-666666666666', 'Makanan', 1),
    ('77777777-7777-7777-7777-777777777772', '66666666-6666-6666-6666-666666666666', 'Minuman', 2),
    ('77777777-7777-7777-7777-777777777773', '66666666-6666-6666-6666-666666666666', 'Camilan', 3)
ON CONFLICT (id) DO NOTHING;

INSERT INTO menu_items (id, company_id, sku, name, base_price)
VALUES
    ('33333333-3333-3333-3333-333333333333', '11111111-1111-1111-1111-111111111111', 'FD-001', 'Nasi Goreng', 28000),
    ('44444444-4444-4444-4444-444444444444', '11111111-1111-1111-1111-111111111111', 'DR-001', 'Es Teh', 8000),
    ('55555555-5555-5555-5555-555555555555', '11111111-1111-1111-1111-111111111111', 'FD-002', 'Ayam Bakar', 32000),
    ('88888888-8888-8888-8888-888888888888', '11111111-1111-1111-1111-111111111111', 'SN-001', 'Kentang Goreng', 18000)
ON CONFLICT (id) DO NOTHING;

INSERT INTO menu_category_items (menu_category_id, menu_item_id, sort_order)
VALUES
    ('77777777-7777-7777-7777-777777777771', '33333333-3333-3333-3333-333333333333', 1),
    ('77777777-7777-7777-7777-777777777772', '44444444-4444-4444-4444-444444444444', 1),
    ('77777777-7777-7777-7777-777777777771', '55555555-5555-5555-5555-555555555555', 2),
    ('77777777-7777-7777-7777-777777777773', '88888888-8888-8888-8888-888888888888', 1)
ON CONFLICT (menu_category_id, menu_item_id) DO NOTHING;

INSERT INTO tables (id, company_id, store_id, table_number, capacity, status)
VALUES
    ('99999999-9999-9999-9999-999999999901', '11111111-1111-1111-1111-111111111111', '22222222-2222-2222-2222-222222222222', 'Meja 01', 2, 'AVAILABLE'),
    ('99999999-9999-9999-9999-999999999902', '11111111-1111-1111-1111-111111111111', '22222222-2222-2222-2222-222222222222', 'Meja 02', 4, 'AVAILABLE'),
    ('99999999-9999-9999-9999-999999999903', '11111111-1111-1111-1111-111111111111', '22222222-2222-2222-2222-222222222222', 'Meja 03', 4, 'AVAILABLE'),
    ('99999999-9999-9999-9999-999999999904', '11111111-1111-1111-1111-111111111111', '22222222-2222-2222-2222-222222222222', 'Meja 04', 6, 'AVAILABLE'),
    ('99999999-9999-9999-9999-999999999905', '11111111-1111-1111-1111-111111111111', '22222222-2222-2222-2222-222222222222', 'Meja 05', 8, 'AVAILABLE')
ON CONFLICT (store_id, table_number) DO NOTHING;

