CREATE TABLE wallet (
    wallet_id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    balance   NUMERIC(15, 2) DEFAULT 100.00 NOT NULL
);
