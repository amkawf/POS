-- Development-only seed data. Not a migration.
-- Matches the temporary DEV_COMPANY_ID / DEV_STORE_ID constants hardcoded
-- in frontend/apps/pos/src/config.ts pending real auth/store context.

INSERT INTO companies (id, code, name)
VALUES ('11111111-1111-1111-1111-111111111111', 'C1', 'Dev Company')
ON CONFLICT (id) DO NOTHING;

INSERT INTO stores (id, company_id, code, name)
VALUES ('22222222-2222-2222-2222-222222222222', '11111111-1111-1111-1111-111111111111', 'S1', 'Dev Store')
ON CONFLICT (id) DO NOTHING;

INSERT INTO menu_items (id, company_id, sku, name, base_price)
VALUES
    ('33333333-3333-3333-3333-333333333333', '11111111-1111-1111-1111-111111111111', 'FD-001', 'Nasi Goreng', 28000),
    ('44444444-4444-4444-4444-444444444444', '11111111-1111-1111-1111-111111111111', 'DR-001', 'Es Teh', 8000),
    ('55555555-5555-5555-5555-555555555555', '11111111-1111-1111-1111-111111111111', 'FD-002', 'Ayam Bakar', 32000)
ON CONFLICT (id) DO NOTHING;
