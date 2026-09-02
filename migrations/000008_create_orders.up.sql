CREATE TABLE orders (
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

CREATE INDEX idx_orders_company_id
    ON orders(company_id);

CREATE INDEX idx_orders_store_id
    ON orders(store_id);

CREATE INDEX idx_orders_store_status
    ON orders(store_id, status);

CREATE INDEX idx_orders_store_opened_at
    ON orders(store_id, opened_at);