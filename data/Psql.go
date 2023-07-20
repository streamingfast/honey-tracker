package data

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
	pb "github.com/streamingfast/honey-tracker/data/pb/hivemapper/v1"
	pbsubstreams "github.com/streamingfast/substreams/pb/sf/substreams/v1"
)

type Psql struct {
	db *sql.DB
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

func NewPostgreSQL(psqlInfo *PsqlInfo) *Psql {
	db, err := sql.Open("postgres", psqlInfo.GetPsqlInfo())
	if err != nil {
		panic(err)
	}
	return &Psql{
		db: db,
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
	result, err := p.db.Exec("INSERT INTO hivemapper.clock (block_num, block_id, block_time) VALUES ($1, $2, $3)", clock.Number, clock.Id, clock.Timestamp)
	if err != nil {
		return 0, fmt.Errorf("inserting clock: %w", err)
	}

	return result.LastInsertId()
}

func (p *Psql) handleTransaction(dbBlockID int64, transactionHash string) (dbTransactionID int64, err error) {
	//todo: create a transaction cache
	rows, err := p.db.Query("SELECT id FROM hivemapper.transactions WHERE hash = $1", transactionHash)
	if err != nil {
		return 0, fmt.Errorf("selecting transaction: %w", err)
	}
	if rows.Next() {
		err = rows.Scan(&dbTransactionID)
		return
	}

	if err != nil {
		result, err := p.db.Exec("INSERT INTO hivemapper.transactions (hash, block_id) VALUES ($1, $2)", transactionHash, dbBlockID)
		if err != nil {
			return 0, fmt.Errorf("inserting transaction: %w", err)
		}
		return result.LastInsertId()
	}
	return
}

func (p *Psql) HandleInitializedAccount(dbBlockID int64, initializedAccounts []*pb.InitializedAccount) (err error) {
	for _, initializedAccount := range initializedAccounts {
		dbTransactionID, err := p.handleTransaction(dbBlockID, initializedAccount.TrxHash)
		if err != nil {
			return fmt.Errorf("handling transaction: %w", err)
		}
		_, err = p.db.Exec("INSERT INTO hivemapper.deriveAddresses (transaction_id, address, deriveAddress) VALUES ($1, $2, $3)", dbTransactionID, initializedAccount.Owner, initializedAccount.Account)
		if err != nil {
			return fmt.Errorf("inserting deriveAddresses: %w", err)
		}
	}
	return nil
}

var NotFound = errors.New("Not found")

func (p *Psql) resolveAddress(derivedAddress string) (string, error) {
	resolvedAddress := ""
	rows, err := p.db.Query("SELECT address FROM hivemapper.derived_addresses WHERE derivedAddress = $1", derivedAddress)
	if err != nil {
		return "", fmt.Errorf("selecting derivedAddresses: %w", err)
	}
	if rows.Next() {
		err = rows.Scan(&resolvedAddress)
		return resolvedAddress, nil
	}

	return "", NotFound
}

func (p *Psql) handleDriver(dbTransactionID int64, driverAddress string) (dbDriverID int64, err error) {
	rows, err := p.db.Query("SELECT id FROM hivemapper.drivers WHERE address = $1", driverAddress)
	if rows.Next() {
		err = rows.Scan(&dbDriverID)
		return
	}

	if err != nil {
		result, err := p.db.Exec("INSERT INTO hivemapper.drivers (address, transaction_id) VALUES ($1, $2)", driverAddress, dbTransactionID)
		if err != nil {
			return 0, fmt.Errorf("inserting driver: %w", err)
		}
		return result.LastInsertId()
	}
	return
}

func (p *Psql) handleFleet(dbTransactionID int64, fleetAddress string) (dbDriverID int64, err error) {
	rows, err := p.db.Query("SELECT id FROM hivemapper.fleets WHERE address = $1", fleetAddress)
	if rows.Next() {
		err = rows.Scan(&dbDriverID)
		return
	}

	if err != nil {
		result, err := p.db.Exec("INSERT INTO hivemapper.fleets (address, transaction_id) VALUES ($1, $2)", fleetAddress, dbTransactionID)
		if err != nil {
			return 0, fmt.Errorf("inserting driver: %w", err)
		}
		return result.LastInsertId()
	}
	return
}

func (p *Psql) handleFleetDriver(dbTransactionID int64, dbFleetID int64, dbDriverID int64) (dbFleetDriverID int64, err error) {
	rows, err := p.db.Query("SELECT id FROM hivemapper.fleet_drivers WHERE fleet_id = $1 and driver_id = $2", dbFleetID, dbDriverID)
	if rows.Next() {
		err = rows.Scan(&dbDriverID)
		return
	}

	if err != nil {
		result, err := p.db.Exec("INSERT INTO hivemapper.fleet_drivers (transaction_id, fleet_id, driver_id) VALUES ($1, $2, $3)", dbTransactionID, dbFleetID, dbDriverID)
		if err != nil {
			return 0, fmt.Errorf("inserting driver: %w", err)
		}
		return result.LastInsertId()
	}
	return
}

func (p *Psql) HandlePayments(dbBlockID int64, payments []*pb.Payment) error {
	for _, payment := range payments {
		dbTransactionID, err := p.handleTransaction(dbBlockID, payment.Mint.TrxHash)
		if err != nil {
			return fmt.Errorf("inserting transaction: %w", err)
		}

		_, err = p.insertMint(dbTransactionID, payment.Mint)
		if err != nil {
			return fmt.Errorf("inserting mint: %w", err)
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

		//todo: detect drive vs fleet from backend api
		_, err = p.insertMint(dbTransactionID, payment.Mint)
		if err != nil {
			return fmt.Errorf("inserting mint: %w", err)
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
			return fmt.Errorf("resolving address: %w", err)
		}

		fleetAddress, err := p.resolveAddress(payment.ManagerMint.To)
		if err != nil {
			return fmt.Errorf("resolving address: %w", err)
		}

		fleetID, err := p.handleFleet(dbTransactionID, fleetAddress)
		if err != nil {
			return fmt.Errorf("handling driver: %w", err)
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

		_, err = p.db.Exec("INSERT INTO hivemapper.split_payments (driver_id, fleet_id, transaction_id, fleet_mint_id, driver_mint_id) VALUES ($1, $2, $3, $4, $5)", driverID, fleetID, dbTransactionID, fleetMintID, driverMintID)
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

		_, err = p.db.Exec("INSERT INTO hivemapper.transfers (transaction_id, from_address, to_address, amount) VALUES ($1, $2, $3, $4, $5, $6)", dbTransactionID, transfer.From, transfer.To, transfer.Amount)
		if err != nil {
			return fmt.Errorf("inserting transfer: %w", err)
		}
	}
	return nil
}

func (p *Psql) insertMint(dbTransactionID int64, mint *pb.Mint) (dbMintID int64, err error) {
	result, err := p.db.Exec("INSERT INTO hivemapper.mints (transaction_id, to_address, amount) VALUES ($1, $2, $3)", dbTransactionID, mint.To, mint.Amount)
	if err != nil {
		return 0, fmt.Errorf("inserting mint: %w", err)
	}
	return result.LastInsertId()
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

func (p *Psql) insertBurns(dbTransactionID int64, burn *pb.Burn) (dbMintID int64, err error) {
	result, err := p.db.Exec("INSERT INTO hivemapper.burns (transaction_id, from_address, amount) VALUES ($1, $2, $3)", dbTransactionID, burn.From, burn.Amount)
	if err != nil {
		return 0, fmt.Errorf("inserting burn: %w", err)
	}
	return result.LastInsertId()
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
		_, err = p.db.Exec("INSERT INTO hivemapper.ai_payments (transaction_id, mint_id) VALUES ($1, $2)", dbTransactionID, dbMintID)
		if err != nil {
			return fmt.Errorf("inserting ai payment: %w", err)
		}
	}

	return nil
}
