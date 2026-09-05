CREATE TABLE kitchen_tickets (
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

CREATE TABLE kitchen_ticket_items (
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

CREATE INDEX idx_kitchen_tickets_store_status
    ON kitchen_tickets(store_id, status);

CREATE INDEX idx_kitchen_tickets_order_id
    ON kitchen_tickets(order_id);

CREATE INDEX idx_kitchen_tickets_created_at
    ON kitchen_tickets(store_id, created_at);

CREATE INDEX idx_kitchen_ticket_items_ticket_id
    ON kitchen_ticket_items(ticket_id);
