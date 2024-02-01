package data

import (
	"database/sql"
	"errors"
	"fmt"

	"go.uber.org/zap"

	_ "github.com/lib/pq"
	pb "github.com/streamingfast/honey-tracker/data/pb/hivemapper/v1"
	sink "github.com/streamingfast/substreams-sink"
	pbsubstreams "github.com/streamingfast/substreams/pb/sf/substreams/v1"
)

type Psql struct {
	db     *sql.DB
	tx     *sql.Tx
	logger *zap.Logger
}

type PsqlInfo struct {
	Host     string
	Port     int
	User     string
	Password string
	Dbname   string
}

func (i *PsqlInfo) GetPsqlInfo() string {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		i.Host, i.Port, i.User, i.Password, i.Dbname,
	)
	return psqlInfo
}

func NewPostgreSQL(psqlInfo *PsqlInfo, logger *zap.Logger) *Psql {
	db, err := sql.Open("postgres", psqlInfo.GetPsqlInfo())
	if err != nil {
		panic(err)
	}
	return &Psql{
		db:     db,
		logger: logger,
	}
}

func (p *Psql) Init() error {
	_, err := p.db.Exec(dbCreateTables)
	if err != nil {
		return fmt.Errorf("creating fleets table: %w", err)
	}

	return nil
}
func (p *Psql) HandleClock(clock *pbsubstreams.Clock) (dbBlockID int64, err error) {
	row := p.tx.QueryRow("INSERT INTO hivemapper.blocks (number, hash, timestamp) VALUES ($1, $2, $3) RETURNING id", clock.Number, clock.Id, clock.Timestamp.AsTime())
	err = row.Err()
	if err != nil {
		return 0, fmt.Errorf("inserting clock: %w", err)
	}

	err = row.Scan(&dbBlockID)
	return
}

func (p *Psql) handleTransaction(dbBlockID int64, transactionHash string) (dbTransactionID int64, err error) {
	//todo: create a transaction cache
	rows, err := p.tx.Query("SELECT id FROM hivemapper.transactions WHERE hash = $1", transactionHash)
	p.logger.Debug("handling transaction", zap.String("trx_hash", transactionHash))
	if err != nil {
		return 0, fmt.Errorf("selecting transaction: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&dbTransactionID)
		return
	}

	row := p.tx.QueryRow("INSERT INTO hivemapper.transactions (hash, block_id) VALUES ($1, $2) RETURNING id", transactionHash, dbBlockID)
	err = row.Err()
	if err != nil {
		return 0, fmt.Errorf("inserting transaction: %w", err)
	}

	err = row.Scan(&dbTransactionID)
	return
}

func (p *Psql) HandleBlockUndo(blockId string) error {
	_, err := p.tx.Exec("DELETE CASCADE FROM solana_tokens.blocks WHERE hash = $1", blockId)
	if err != nil {
		return fmt.Errorf("deleting block: %w", err)
	}
	return nil
}

func (p *Psql) HandleInitializedAccount(dbBlockID int64, initializedAccounts []*pb.InitializedAccount) (err error) {
	for _, initializedAccount := range initializedAccounts {
		dbTransactionID, err := p.handleTransaction(dbBlockID, initializedAccount.TrxHash)
		if err != nil {
			return fmt.Errorf("handling transaction: %w", err)
		}
		_, err = p.tx.Exec("INSERT INTO hivemapper.derived_addresses (transaction_id, address, derivedAddress) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING", dbTransactionID, initializedAccount.Owner, initializedAccount.Account)
		if err != nil {
			return fmt.Errorf("trx_hash: %d inserting derived_addresses: %w", dbBlockID, err)
		}
	}
	return nil
}

var NotFound = errors.New("Not found")

func (p *Psql) resolveAddress(derivedAddress string) (string, error) {
	resolvedAddress := ""
	rows, err := p.tx.Query("SELECT address FROM hivemapper.derived_addresses WHERE derivedAddress = $1", derivedAddress)
	if err != nil {
		return "", fmt.Errorf("selecting derived_addresses: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&resolvedAddress)
		return resolvedAddress, nil
	}

	return "", NotFound
}

func (p *Psql) handleDriver(dbTransactionID int64, driverAddress string) (dbDriverID int64, err error) {
	rows, err := p.tx.Query("SELECT id FROM hivemapper.drivers WHERE address = $1", driverAddress)
	if err != nil {
		return 0, fmt.Errorf("selecting drivers %q : %w", driverAddress, err)
	}

	if rows.Next() {
		err = rows.Scan(&dbDriverID)
		rows.Close()
		return
	}
	rows.Close()

	row := p.tx.QueryRow("INSERT INTO hivemapper.drivers (address, transaction_id) VALUES ($1, $2) RETURNING id", driverAddress, dbTransactionID)
	err = row.Err()
	if err != nil {
		return 0, fmt.Errorf("inserting driver: %w", err)
	}

	err = row.Scan(&dbDriverID)
	return
}

func (p *Psql) handleFleet(dbTransactionID int64, fleetAddress string) (dbDriverID int64, err error) {
	rows, err := p.tx.Query("SELECT id FROM hivemapper.fleets WHERE address = $1", fleetAddress)
	if err != nil {
		return 0, fmt.Errorf("selecting fleets %q: %w", fleetAddress, err)
	}
	if rows.Next() {
		err = rows.Scan(&dbDriverID)
		rows.Close()
		return
	}

	rows.Close()

	row := p.tx.QueryRow("INSERT INTO hivemapper.fleets (address, transaction_id) VALUES ($1, $2) RETURNING id", fleetAddress, dbTransactionID)
	err = row.Err()
	if err != nil {
		return 0, fmt.Errorf("inserting driver: %w", err)
	}

	err = row.Scan(&dbDriverID)
	return
}

func (p *Psql) handleFleetDriver(dbTransactionID int64, dbFleetID int64, dbDriverID int64) (dbFleetDriverID int64, err error) {
	rows, err := p.tx.Query("SELECT id FROM hivemapper.fleet_drivers WHERE fleet_id = $1 and driver_id = $2", dbFleetID, dbDriverID)
	if err != nil {
		return 0, fmt.Errorf("selecting fleet_drivers: %w", err)
	}

	if rows.Next() {
		err = rows.Scan(&dbDriverID)
		rows.Close()
		return
	}
	rows.Close()

	row := p.tx.QueryRow("INSERT INTO hivemapper.fleet_drivers (transaction_id, fleet_id, driver_id) VALUES ($1, $2, $3) RETURNING id", dbTransactionID, dbFleetID, dbDriverID)
	err = row.Err()
	if err != nil {
		return 0, fmt.Errorf("inserting driver: %w", err)
	}

	err = row.Scan(&dbDriverID)
	return

}

func (p *Psql) HandleRegularDriverPayments(dbBlockID int64, payments []*pb.RegularDriverPayment) error {
	for _, payment := range payments {
		dbTransactionID, err := p.handleTransaction(dbBlockID, payment.Mint.TrxHash)
		if err != nil {
			return fmt.Errorf("inserting transaction: %w", err)
		}

		mintDbID, err := p.insertMint(dbTransactionID, payment.Mint)
		if err != nil {
			return fmt.Errorf("inserting mint: %w", err)
		}

		_, err = p.tx.Exec("INSERT INTO hivemapper.payments (mint_id) VALUES ($1) RETURNING id", mintDbID)
		if err != nil {
			return fmt.Errorf("inserting payment: %w", err)
		}
	}
	return nil
}

func (p *Psql) HandleNoneSplitPayments(dbBlockID int64, payments []*pb.NoSplitPayment) error {
	for _, payment := range payments {
		dbTransactionID, err := p.handleTransaction(dbBlockID, payment.Mint.TrxHash)
		if err != nil {
			return fmt.Errorf("inserting transaction: %w", err)
		}

		mintDbID, err := p.insertMint(dbTransactionID, payment.Mint)
		if err != nil {
			return fmt.Errorf("inserting mint: %w", err)
		}

		//todo: detect drive vs fleet from backend api
		_, err = p.tx.Exec("INSERT INTO hivemapper.no_split_payments (mint_id) VALUES ($1) RETURNING id", mintDbID)
		if err != nil {
			return fmt.Errorf("inserting NoneSplitPayments with mint_id %d mint_to %q tx %q: %w", mintDbID, payment.Mint.To, payment.Mint.TrxHash, err)
		}
	}
	return nil
}

func (p *Psql) HandleSplitPayments(dbBlockID int64, splitPayments []*pb.TokenSplittingPayment) error {
	for _, payment := range splitPayments {
		dbTransactionID, err := p.handleTransaction(dbBlockID, payment.ManagerMint.TrxHash)
		if err != nil {
			return fmt.Errorf("inserting transaction: %w", err)
		}

		driverAddress, err := p.resolveAddress(payment.DriverMint.To)
		if err != nil {
			return fmt.Errorf("resolving driver address: %w", err)
		}

		fleetAddress, err := p.resolveAddress(payment.ManagerMint.To)
		if err != nil {
			return fmt.Errorf("resolving fleet address: %w", err)
		}

		fleetID, err := p.handleFleet(dbTransactionID, fleetAddress)
		if err != nil {
			return fmt.Errorf("handling fleet id: %w", err)
		}

		driverID, err := p.handleDriver(dbTransactionID, driverAddress)
		if err != nil {
			return fmt.Errorf("handling driver: %w", err)
		}

		_, err = p.handleFleetDriver(dbTransactionID, fleetID, driverID)
		if err != nil {
			return fmt.Errorf("handling fleet driver: %w", err)
		}

		fleetMintID, err := p.insertMint(dbTransactionID, payment.ManagerMint)
		if err != nil {
			return fmt.Errorf("inserting fleet mint: %w", err)
		}

		driverMintID, err := p.insertMint(dbTransactionID, payment.DriverMint)
		if err != nil {
			return fmt.Errorf("inserting driver mint: %w", err)
		}

		_, err = p.tx.Exec("INSERT INTO hivemapper.split_payments (transaction_id, fleet_mint_id, driver_mint_id) VALUES ($1, $2, $3)", dbTransactionID, fleetMintID, driverMintID)
		if err != nil {
			return fmt.Errorf("inserting split payment: %w", err)
		}
	}

	return nil
}

func (p *Psql) HandleTransfers(dbBlockID int64, transfers []*pb.Transfer) error {
	for _, transfer := range transfers {
		dbTransactionID, err := p.handleTransaction(dbBlockID, transfer.TrxHash)
		if err != nil {
			return fmt.Errorf("inserting transaction: %w", err)
		}
		//{"error": "handle BlockScopedData message: rollback transaction: rolling back transaction: driver: bad connection: while handling err handle transfers: inserting transfer: pq: unexpected Parse response 'C'"}
		_, err = p.tx.Exec("INSERT INTO hivemapper.transfers (transaction_id, from_address, to_address, amount) VALUES ($1, $2, $3, $4)", dbTransactionID, transfer.From, transfer.To, transfer.Amount)
		if err != nil {
			fmt.Println("processing transfer: ", transfer.From, transfer.To, transfer.Amount, transfer.TrxHash)
			return fmt.Errorf("inserting transfer: %w", err)
		}
	}
	return nil
}

func (p *Psql) insertMint(dbTransactionID int64, mint *pb.Mint) (dbMintID int64, err error) {
	row := p.tx.QueryRow("INSERT INTO hivemapper.mints (transaction_id, to_address, amount) VALUES ($1, $2, $3) RETURNING id", dbTransactionID, mint.To, mint.Amount)
	err = row.Err()
	if err != nil {
		return 0, fmt.Errorf("inserting mint: %w", err)
	}

	err = row.Scan(&dbMintID)
	return
}

func (p *Psql) HandleMints(dbBlockID int64, mints []*pb.Mint) error {
	for _, mint := range mints {
		dbTransactionID, err := p.handleTransaction(dbBlockID, mint.TrxHash)
		if err != nil {
			return fmt.Errorf("inserting transaction: %w", err)
		}

		_, err = p.insertMint(dbTransactionID, mint)
		if err != nil {
			return fmt.Errorf("inserting mint: %w", err)
		}
	}
	return nil
}
func (p *Psql) HandleAITrainerPayments(dbBlockID int64, payments []*pb.AiTrainerPayment) error {
	for _, payment := range payments {
		dbTransactionID, err := p.handleTransaction(dbBlockID, payment.Mint.TrxHash)
		if err != nil {
			return fmt.Errorf("inserting transaction: %w", err)
		}

		mintDbID, err := p.insertMint(dbTransactionID, payment.Mint)
		if err != nil {
			return fmt.Errorf("inserting mint: %w", err)
		}

		_, err = p.tx.Exec("INSERT INTO hivemapper.ai_payments (mint_id) VALUES ($1) RETURNING id", mintDbID)
		if err != nil {
			return fmt.Errorf("inserting payment: %w", err)
		}
	}

	return nil
}

func (p *Psql) insertBurns(dbTransactionID int64, burn *pb.Burn) (dbMintID int64, err error) {
	row := p.tx.QueryRow("INSERT INTO hivemapper.burns (transaction_id, from_address, amount) VALUES ($1, $2, $3) RETURNING id", dbTransactionID, burn.From, burn.Amount)
	err = row.Err()
	if err != nil {
		return 0, fmt.Errorf("inserting burn: %w", err)
	}

	err = row.Scan(&dbMintID)
	return
}

func (p *Psql) HandleBurns(dbBlockID int64, burns []*pb.Burn) error {
	for _, burn := range burns {
		dbTransactionID, err := p.handleTransaction(dbBlockID, burn.TrxHash)
		if err != nil {
			return fmt.Errorf("inserting transaction: %w", err)
		}

		_, err = p.insertBurns(dbTransactionID, burn)
		if err != nil {
			return fmt.Errorf("inserting burn: %w", err)
		}
	}
	return nil
}

func (p *Psql) HandleAiPayments(dbBlockID int64, payments []*pb.AiTrainerPayment) error {
	for _, payment := range payments {
		dbTransactionID, err := p.handleTransaction(dbBlockID, payment.Mint.TrxHash)
		if err != nil {
			return fmt.Errorf("inserting transaction: %w", err)
		}

		dbMintID, err := p.insertMint(dbTransactionID, payment.Mint)
		if err != nil {
			return fmt.Errorf("inserting mint: %w", err)
		}
		_, err = p.tx.Exec("INSERT INTO hivemapper.ai_payments (transaction_id, mint_id) VALUES ($1, $2)", dbTransactionID, dbMintID)
		if err != nil {
			return fmt.Errorf("inserting ai payment: %w", err)
		}
	}

	return nil
}

func (p *Psql) StoreCursor(cursor *sink.Cursor) error {
	_, err := p.tx.Exec("INSERT INTO hivemapper.cursor (name, cursor) VALUES ($1, $2) ON CONFLICT (name) DO UPDATE SET cursor = $2", "hivemapper", cursor.String())
	if err != nil {
		return fmt.Errorf("inserting cursor: %w", err)
	}
	return nil
}

func (p *Psql) FetchCursor() (*sink.Cursor, error) {
	rows, err := p.db.Query("SELECT cursor FROM hivemapper.cursor WHERE name = $1", "hivemapper")
	if err != nil {
		return nil, fmt.Errorf("selecting cursor: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		var cursor string
		err = rows.Scan(&cursor)

		return sink.NewCursor(cursor)
	}
	return nil, nil
}

func (p *Psql) BeginTransaction() error {
	tx, err := p.db.Begin()
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}
	p.tx = tx
	return nil
}

func (p *Psql) CommitTransaction() error {
	err := p.tx.Commit()
	if err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	p.tx = nil
	return nil
}

func (p *Psql) RollbackTransaction() error {
	err := p.tx.Rollback()
	if err != nil {
		return fmt.Errorf("rolling back transaction: %w", err)
	}

	p.tx = nil
	return nil
}
