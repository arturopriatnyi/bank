CREATE TABLE transactions (
    id VARCHAR(36) PRIMARY KEY,
    sender_id INTEGER NOT NULL,
    receiver_id INTEGER NOT NULL,
    money_value INTEGER NOT NULL,
    money_currency VARCHAR(3) NOT NULL
);
