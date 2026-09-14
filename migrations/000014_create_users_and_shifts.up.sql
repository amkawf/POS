-- Aktifkan ekstensi pgcrypto untuk hashing bcrypt langsung di PostgreSQL jika belum aktif
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- 1. Tabel Master Akun Pengguna & Staf Restoran
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(150),
    password_hash VARCHAR(255),
    pin_hash VARCHAR(255) NOT NULL,
    role VARCHAR(30) NOT NULL CHECK (role IN ('OWNER', 'MANAGER', 'CASHIER', 'KITCHEN')),
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE', 'INACTIVE')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_company_email UNIQUE (company_id, email)
);

CREATE INDEX IF NOT EXISTS idx_users_company ON users(company_id);
CREATE INDEX IF NOT EXISTS idx_users_company_status ON users(company_id, status);

-- 2. Tabel Relasi Staf ke Cabang Toko (User-Store Mapping)
CREATE TABLE IF NOT EXISTS user_stores (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    store_id UUID NOT NULL REFERENCES stores(id) ON DELETE CASCADE,
    assigned_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, store_id)
);

CREATE INDEX IF NOT EXISTS idx_user_stores_store ON user_stores(store_id);

-- 3. Tabel Sesi Buka/Tutup Kasir (Cashier Register Shift)
CREATE TABLE IF NOT EXISTS cashier_shifts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    store_id UUID NOT NULL REFERENCES stores(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    opened_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    starting_cash NUMERIC(15, 2) NOT NULL DEFAULT 0,
    closed_at TIMESTAMPTZ,
    actual_ending_cash NUMERIC(15, 2),
    expected_ending_cash NUMERIC(15, 2),
    notes TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'OPEN' CHECK (status IN ('OPEN', 'CLOSED')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_cashier_shifts_lookup ON cashier_shifts(store_id, user_id, status);

-- 4. Seed Data Staf Awal untuk Demo / Testing Cepat
-- Memasukkan 3 akun staf di DEV_COMPANY_ID & DEV_STORE_ID dengan PIN bcrypt terenkripsi
DO $$
DECLARE
    dev_comp_id UUID := '11111111-1111-1111-1111-111111111111';
    dev_str_id  UUID := '22222222-2222-2222-2222-222222222222';
    uid_cashier UUID;
    uid_kitchen UUID;
    uid_manager UUID;
BEGIN
    -- Akun 1: Kasir (PIN: 1234)
    INSERT INTO users (company_id, name, email, pin_hash, role)
    VALUES (
        dev_comp_id, 
        'Budi Santoso', 
        'kasir@pos.local', 
        crypt('1234', gen_salt('bf', 10)), 
        'CASHIER'
    )
    ON CONFLICT (company_id, email) DO UPDATE 
        SET pin_hash = EXCLUDED.pin_hash
    RETURNING id INTO uid_cashier;

    -- Akun 2: Koki Dapur (PIN: 5678)
    INSERT INTO users (company_id, name, email, pin_hash, role)
    VALUES (
        dev_comp_id, 
        'Chef Juna', 
        'kitchen@pos.local', 
        crypt('5678', gen_salt('bf', 10)), 
        'KITCHEN'
    )
    ON CONFLICT (company_id, email) DO UPDATE 
        SET pin_hash = EXCLUDED.pin_hash
    RETURNING id INTO uid_kitchen;

    -- Akun 3: Manager Toko (PIN: 9999, Password: password123)
    INSERT INTO users (company_id, name, email, password_hash, pin_hash, role)
    VALUES (
        dev_comp_id, 
        'Rian Manager', 
        'manager@pos.local', 
        crypt('password123', gen_salt('bf', 10)), 
        crypt('9999', gen_salt('bf', 10)), 
        'MANAGER'
    )
    ON CONFLICT (company_id, email) DO UPDATE 
        SET pin_hash = EXCLUDED.pin_hash, password_hash = EXCLUDED.password_hash
    RETURNING id INTO uid_manager;

    -- Petakan ketiga staf ke Toko Utama (DEV_STORE_ID)
    IF uid_cashier IS NOT NULL THEN
        INSERT INTO user_stores (user_id, store_id) VALUES (uid_cashier, dev_str_id) ON CONFLICT DO NOTHING;
    END IF;
    IF uid_kitchen IS NOT NULL THEN
        INSERT INTO user_stores (user_id, store_id) VALUES (uid_kitchen, dev_str_id) ON CONFLICT DO NOTHING;
    END IF;
    IF uid_manager IS NOT NULL THEN
        INSERT INTO user_stores (user_id, store_id) VALUES (uid_manager, dev_str_id) ON CONFLICT DO NOTHING;
    END IF;
END $$;

