CREATE TABLE menu_items (
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

CREATE INDEX idx_menu_items_company_id
    ON menu_items(company_id);