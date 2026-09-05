-- +goose Up
-- +goose StatementBegin

CREATE TYPE contract_status AS ENUM ('draft', 'active', 'suspended', 'terminated');

CREATE TABLE contracts (
    id uuid PRIMARY KEY,
    company_id uuid NOT NULL,
    status contract_status NOT NULL,
    valid_from timestamptz NOT NULL,
    valid_until timestamptz,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT contracts_valid_period_check CHECK (valid_until IS NULL OR valid_until > valid_from)
);

CREATE UNIQUE INDEX uidx_contracts_one_active_per_company
    ON contracts (company_id)
    WHERE status = 'active'::contract_status;


CREATE TYPE contract_service_code AS ENUM ('trip_creation', 'trip_participants', 'notifications');

CREATE TABLE contract_services (
    contract_id uuid NOT NULL,
    service_code contract_service_code NOT NULL,
    allowed boolean NOT NULL DEFAULT false,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT contract_services_pkey PRIMARY KEY (contract_id, service_code),
    CONSTRAINT fk_contract_services_contract_id FOREIGN KEY (contract_id) REFERENCES contracts (id) ON DELETE CASCADE
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS contract_services;
DROP TYPE IF EXISTS contract_service_code;
DROP INDEX IF EXISTS uidx_contracts_one_active_per_company;
DROP TABLE IF EXISTS contracts;
DROP TYPE IF EXISTS contract_status;

-- +goose StatementEnd
