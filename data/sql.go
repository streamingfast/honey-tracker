package data

const fleetsCreateTable = `
CREATE TABLE IF NOT EXISTS fleets (
	id SERIAL PRIMARY KEY,
	address TEXT NOT NULL,
);
`

const driversCreateTable = `
CREATE TABLE IF NOT EXISTS drivers (
	id SERIAL PRIMARY KEY,
	address TEXT NOT NULL,
);
`

const fleetDriversCreateTable = `
CREATE TABLE IF NOT EXISTS fleet_drivers (
	id SERIAL PRIMARY KEY,
	fleet_id INTEGER NOT NULL,
	driver_id INTEGER NOT NULL,
	CONSTRAINT fk_fleet FOREIGN KEY (fleet_id) REFERENCES fleets(id),
	CONSTRAINT fk_driver FOREIGN KEY (driver_id) REFERENCES drivers(id),
);
`

const paymentsCreateTable = `
CREATE TABLE IF NOT EXISTS payments (
	id SERIAL PRIMARY KEY,
	trx_id TEXT NOT NULL,
	date TIMESTAMP NOT NULL,
	to_address TEXT NOT NULL,
	amount FLOAT NOT NULL,
);
`

const splitPaymentsCreateTable = `
CREATE TABLE IF NOT EXISTS split_payments (
	id SERIAL PRIMARY KEY,
	trx_id TEXT NOT NULL,
	date TIMESTAMP NOT NULL,
	fleet_payment_id INTEGER NOT NULL,
	driver_payment_id INTEGER NOT NULL,
	CONSTRAINT fk_fleet_payment FOREIGN KEY (fleet_payment_id) REFERENCES payments(id),
	CONSTRAINT fk_driver_payment FOREIGN KEY (driver_payment_id) REFERENCES payments(id),
);
`

const transfersCreateTable = `
CREATE TABLE IF NOT EXISTS transfers (
	id SERIAL PRIMARY KEY,
	trx_id TEXT NOT NULL,
	date TIMESTAMP NOT NULL,
	from_owner_address TEXT NOT NULL,
	from_address TEXT NOT NULL,
	to_owner_address TEXT NOT NULL,
	to_address TEXT NOT NULL,
	amount FLOAT NOT NULL,
);
`
