CREATE TYPE loan_status AS ENUM (
'proposed','approved','invested','disbursed'
);

CREATE TABLE loan (
    loan_id uuid PRIMARY KEY,
    principal_amount bigint NOT NULL,
    interest_rate bigint NOT NULL,
    agreement_letter text,
    status loan_status NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone,
    disburse_at timestamp with time zone,
    borrower_id uuid NOT NULL
);

CREATE TABLE loan_investment (
    loan_investment_id uuid PRIMARY KEY,
    loan_id uuid NOT NULL,
    investor_id uuid NOT NULL,
    amount bigint NOT NULL,
    invested_at timestamp with time zone NOT NULL
);

CREATE TABLE loan_status_history (
    loan_status_history_id BIGSERIAL PRIMARY KEY,
    loan_id uuid NOT NULL,
    before jsonb,
    after jsonb,
    updated_by uuid,
    updated_at timestamp with time zone
);
