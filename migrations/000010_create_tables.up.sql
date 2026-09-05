CREATE TABLE tables (
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

CREATE INDEX idx_tables_store_id
    ON tables(store_id);

CREATE INDEX idx_tables_store_status
    ON tables(store_id, status);

CREATE INDEX idx_tables_company_id
    ON tables(company_id);

ALTER TABLE orders
    ADD CONSTRAINT fk_orders_table
    FOREIGN KEY (table_id)
    REFERENCES tables(id)
    ON DELETE SET NULL;
