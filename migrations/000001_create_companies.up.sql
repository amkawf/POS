CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE companies (
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