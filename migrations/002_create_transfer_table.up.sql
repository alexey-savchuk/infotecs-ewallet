CREATE TABLE transfer (
    transfer_id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    time        TIMESTAMP DEFAULT now() NOT NULL,
    from_wallet UUID REFERENCES wallet(wallet_id) NOT NULL,
    to_wallet   UUID REFERENCES wallet(wallet_id) NOT NULL,
    amount      NUMERIC(15, 2) DEFAULT 100.00 NOT NULL
);
