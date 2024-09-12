package data

const dbCreateTables = `
CREATE TABLE IF NOT EXISTS hivemapper.cursor (
	name TEXT PRIMARY KEY,
	cursor TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS hivemapper.blocks (
	id SERIAL PRIMARY KEY,
	number INTEGER NOT NULL,
	hash TEXT NOT NULL,
	hash TEXT NOT NULL,
	timestamp TIMESTAMP NOT NULL
);
CREATE TABLE IF NOT EXISTS hivemapper.prices (
	timestamp TIMESTAMP PRIMARY KEY,
	price DECIMAL NOT NULL
);

CREATE TABLE IF NOT EXISTS hivemapper.transactions (
	id SERIAL PRIMARY KEY,
	block_id INTEGER NOT NULL,
	hash TEXT NOT NULL UNIQUE,
	CONSTRAINT fk_block FOREIGN KEY (block_id) REFERENCES hivemapper.blocks(id)
);

CREATE TABLE IF NOT EXISTS hivemapper.derived_addresses (
	id SERIAL PRIMARY KEY,
	transaction_id INTEGER NOT NULL,
	address TEXT NOT NULL,
	derivedAddress TEXT NOT NULL UNIQUE,
	CONSTRAINT fk_transaction FOREIGN KEY (transaction_id) REFERENCES hivemapper.transactions(id)
);


CREATE TABLE IF NOT EXISTS hivemapper.fleets (
	id SERIAL PRIMARY KEY,
	transaction_id INTEGER NOT NULL,
	address TEXT NOT NULL UNIQUE,
	CONSTRAINT fk_transaction FOREIGN KEY (transaction_id) REFERENCES hivemapper.transactions(id)
);

CREATE TABLE IF NOT EXISTS hivemapper.drivers (
	id SERIAL PRIMARY KEY,
	transaction_id INTEGER NOT NULL,
	address TEXT NOT NULL UNIQUE,
	CONSTRAINT fk_transaction FOREIGN KEY (transaction_id) REFERENCES hivemapper.transactions(id)
);

CREATE TABLE IF NOT EXISTS hivemapper.fleet_drivers (
	id SERIAL PRIMARY KEY,
	transaction_id INTEGER NOT NULL,
	fleet_id INTEGER NOT NULL,
	driver_id INTEGER NOT NULL,
	CONSTRAINT fk_fleet FOREIGN KEY (fleet_id) REFERENCES hivemapper.fleets(id),
	CONSTRAINT fk_driver FOREIGN KEY (driver_id) REFERENCES hivemapper.drivers(id),
	CONSTRAINT fk_transaction FOREIGN KEY (transaction_id) REFERENCES hivemapper.transactions(id)
);

CREATE TABLE IF NOT EXISTS hivemapper.mints (
	id SERIAL PRIMARY KEY,
	transaction_id INTEGER NOT NULL,
	to_address TEXT NOT NULL,
	amount DECIMAL NOT NULL,
	CONSTRAINT fk_transaction FOREIGN KEY (transaction_id) REFERENCES hivemapper.transactions(id)
);

CREATE TABLE IF NOT EXISTS hivemapper.burns (
	id SERIAL PRIMARY KEY,
	transaction_id INTEGER NOT NULL,
	from_address TEXT NOT NULL,
	amount DECIMAL NOT NULL,
	CONSTRAINT fk_transaction FOREIGN KEY (transaction_id) REFERENCES hivemapper.transactions(id)
);

CREATE TABLE IF NOT EXISTS hivemapper.payments (
	id SERIAL PRIMARY KEY,
	mint_id INTEGER NOT NULL,
	CONSTRAINT fk_mint FOREIGN KEY (mint_id) REFERENCES hivemapper.mints(id)
);

CREATE TABLE IF NOT EXISTS hivemapper.map_create (
	id SERIAL PRIMARY KEY,
	burn_id INTEGER NOT NULL,
	CONSTRAINT fk_mint FOREIGN KEY (burn_id) REFERENCES hivemapper.burns(id)
);

CREATE TABLE IF NOT EXISTS hivemapper.map_consumption_reward (
	id SERIAL PRIMARY KEY,
	mint_id INTEGER NOT NULL,
	CONSTRAINT fk_mint FOREIGN KEY (mint_id) REFERENCES hivemapper.mints(id)
);

CREATE TABLE IF NOT EXISTS hivemapper.ai_payments (
	id SERIAL PRIMARY KEY,
	mint_id INTEGER NOT NULL,
	CONSTRAINT fk_mint FOREIGN KEY (mint_id) REFERENCES hivemapper.mints(id)
);

CREATE TABLE IF NOT EXISTS hivemapper.operational_payments (
	id SERIAL PRIMARY KEY,
	mint_id INTEGER NOT NULL,
	CONSTRAINT fk_mint FOREIGN KEY (mint_id) REFERENCES hivemapper.mints(id)
);
CREATE TABLE IF NOT EXISTS hivemapper.reward_payments (
	id SERIAL PRIMARY KEY,
	mint_id INTEGER NOT NULL,
	CONSTRAINT fk_mint FOREIGN KEY (mint_id) REFERENCES hivemapper.mints(id)
);

CREATE TABLE IF NOT EXISTS hivemapper.split_payments (
	id SERIAL PRIMARY KEY,
	transaction_id INTEGER NOT NULL,
	driver_mint_id INTEGER NOT NULL,
	fleet_mint_id INTEGER NOT NULL,
	CONSTRAINT fk_transaction FOREIGN KEY (transaction_id) REFERENCES hivemapper.transactions(id),
	CONSTRAINT fk_driver_mint FOREIGN KEY (driver_mint_id) REFERENCES hivemapper.mints(id),
	CONSTRAINT fk_fleet_mint FOREIGN KEY (fleet_mint_id) REFERENCES hivemapper.mints(id)
);

CREATE TABLE IF NOT EXISTS hivemapper.no_split_payments (
	id SERIAL PRIMARY KEY,
	mint_id INTEGER NOT NULL,
	CONSTRAINT fk_mint FOREIGN KEY (mint_id) REFERENCES hivemapper.mints(id)
);

CREATE TABLE IF NOT EXISTS hivemapper.transfers (
	id SERIAL PRIMARY KEY,
	transaction_id INTEGER NOT NULL,
	from_address TEXT NOT NULL,
	to_address TEXT NOT NULL,
	amount DECIMAL NOT NULL,
	CONSTRAINT fk_transaction FOREIGN KEY (transaction_id) REFERENCES hivemapper.transactions(id)
);

`
