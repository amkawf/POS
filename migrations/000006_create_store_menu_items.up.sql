CREATE TABLE store_menu_items (
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

CREATE INDEX idx_store_menu_items_menu_item_id
    ON store_menu_items(menu_item_id);