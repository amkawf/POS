CREATE TABLE menu_category_items (
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

CREATE INDEX idx_menu_category_items_menu_item_id
    ON menu_category_items(menu_item_id);