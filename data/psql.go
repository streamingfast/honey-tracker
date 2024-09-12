package data

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"

	_ "github.com/lib/pq"
	pb "github.com/streamingfast/honey-tracker/data/pb/hivemapper/v1"
	sink "github.com/streamingfast/substreams-sink"
	pbsubstreams "github.com/streamingfast/substreams/pb/sf/substreams/v1"
)

type PreparedStatement struct {
	insertMint               *sql.Stmt
	insertTransaction        *sql.Stmt
	insertBlock              *sql.Stmt
	insertDerivedAddress     *sql.Stmt
	selectDerivedAddress     *sql.Stmt
	selectDriver             *sql.Stmt
	insertDriver             *sql.Stmt
	selectFleet              *sql.Stmt
	insertFleet              *sql.Stmt
	selectFleetDriver        *sql.Stmt
	insertFleetDriver        *sql.Stmt
	insertPayment            *sql.Stmt
	insertSlipPayment        *sql.Stmt
	insertTransfer           *sql.Stmt
	insertAIPayment          *sql.Stmt
	insertOperationalPayment *sql.Stmt
	insertRewardPayment      *sql.Stmt
	insertMapCreate          *sql.Stmt
	insertMapConsumption     *sql.Stmt
	insertBurn               *sql.Stmt
	insertCursor             *sql.Stmt
	insertPrice              *sql.Stmt
}

var preparedStatement *PreparedStatement

type Psql struct {
	db             *sql.DB
	tx             *sql.Tx
	logger         *zap.Logger
	TransactionIDs map[string]int64
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

	insertMing, err := db.Prepare("INSERT INTO hivemapper.mints (transaction_id, to_address, amount) VALUES ($1, $2, $3) RETURNING id")
	if err != nil {
		panic(err)
	}

	insertTransaction, err := db.Prepare("INSERT INTO hivemapper.transactions (hash, block_id) VALUES ($1, $2) RETURNING id")
	if err != nil {
		panic(err)
	}

	insertBlock, err := db.Prepare("INSERT INTO hivemapper.blocks (number, hash, timestamp) VALUES ($1, $2, $3) RETURNING id")
	if err != nil {
		panic(err)
	}
	insertDerivedAddress, err := db.Prepare("INSERT INTO hivemapper.derived_addresses (transaction_id, address, derivedAddress) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING")
	if err != nil {
		panic(err)
	}
	selectDerivedAddress, err := db.Prepare("SELECT address FROM hivemapper.derived_addresses WHERE derivedAddress = $1")
	if err != nil {
		panic(err)
	}
	selectDriver, err := db.Prepare("SELECT id FROM hivemapper.drivers WHERE address = $1")
	if err != nil {
		panic(err)
	}
	insertDriver, err := db.Prepare("INSERT INTO hivemapper.drivers (address, transaction_id) VALUES ($1, $2) RETURNING id")
	if err != nil {
		panic(err)
	}
	selectFleet, err := db.Prepare("SELECT id FROM hivemapper.fleets WHERE address = $1")
	if err != nil {
		panic(err)
	}
	insertFleet, err := db.Prepare("INSERT INTO hivemapper.fleets (address, transaction_id) VALUES ($1, $2) RETURNING id")
	if err != nil {
		panic(err)
	}
	selectFleetDriver, err := db.Prepare("SELECT id FROM hivemapper.fleet_drivers WHERE fleet_id = $1 and driver_id = $2")
	if err != nil {
		panic(err)
	}
	insertFleetDriver, err := db.Prepare("INSERT INTO hivemapper.fleet_drivers (transaction_id, fleet_id, driver_id) VALUES ($1, $2, $3) RETURNING id")
	if err != nil {
		panic(err)
	}
	insertPayment, err := db.Prepare("INSERT INTO hivemapper.payments (mint_id) VALUES ($1) RETURNING id")
	if err != nil {
		panic(err)
	}
	insertSlipPayment, err := db.Prepare("INSERT INTO hivemapper.split_payments (transaction_id, fleet_mint_id, driver_mint_id) VALUES ($1, $2, $3)")
	if err != nil {
		panic(err)
	}
	insertTransfer, err := db.Prepare("INSERT INTO hivemapper.transfers (transaction_id, from_address, to_address, amount) VALUES ($1, $2, $3, $4)")
	if err != nil {
		panic(err)
	}
	insertAIPayment, err := db.Prepare("INSERT INTO hivemapper.ai_payments (mint_id) VALUES ($1) RETURNING id")
	if err != nil {
		panic(err)
	}
	insertOperationalPayment, err := db.Prepare("INSERT INTO hivemapper.operational_payments (mint_id) VALUES ($1) RETURNING id")
	if err != nil {
		panic(err)
	}
	insertRewardPayment, err := db.Prepare("INSERT INTO hivemapper.reward_payments (mint_id) VALUES ($1) RETURNING id")
	if err != nil {
		panic(err)
	}
	insertMapCreate, err := db.Prepare("INSERT INTO hivemapper.map_create (burn_id) VALUES ($1) RETURNING id")
	if err != nil {
		panic(err)
	}
	insertMapConsumption, err := db.Prepare("INSERT INTO hivemapper.map_consumption_reward (mint_id) VALUES ($1) RETURNING id")
	if err != nil {
		panic(err)
	}
	insertBurn, err := db.Prepare("INSERT INTO hivemapper.burns (transaction_id, from_address, amount) VALUES ($1, $2, $3) RETURNING id")
	if err != nil {
		panic(err)
	}
	insertCursor, err := db.Prepare("INSERT INTO hivemapper.cursor (name, cursor) VALUES ($1, $2) ON CONFLICT (name) DO UPDATE SET cursor = $2")
	if err != nil {
		panic(err)
	}

	insertPrice, err := db.Prepare("INSERT INTO hivemapper.prices (timestamp, price) VALUES ($1, $2) ON CONFLICT (timestamp) DO NOTHING")
	if err != nil {
		panic(err)
	}

	preparedStatement = &PreparedStatement{
		insertMint:               insertMing,
		insertTransaction:        insertTransaction,
		insertBlock:              insertBlock,
		insertDerivedAddress:     insertDerivedAddress,
		selectDerivedAddress:     selectDerivedAddress,
		selectDriver:             selectDriver,
		insertDriver:             insertDriver,
		selectFleet:              selectFleet,
		insertFleet:              insertFleet,
		selectFleetDriver:        selectFleetDriver,
		insertFleetDriver:        insertFleetDriver,
		insertPayment:            insertPayment,
		insertSlipPayment:        insertSlipPayment,
		insertTransfer:           insertTransfer,
		insertAIPayment:          insertAIPayment,
		insertOperationalPayment: insertOperationalPayment,
		insertRewardPayment:      insertRewardPayment,
		insertMapCreate:          insertMapCreate,
		insertMapConsumption:     insertMapConsumption,
		insertBurn:               insertBurn,
		insertCursor:             insertCursor,
		insertPrice:              insertPrice,
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
	row := p.tx.Stmt(preparedStatement.insertBlock).QueryRow(clock.Number, clock.Id, clock.Timestamp.AsTime())
	err = row.Err()
	if err != nil {
		return 0, fmt.Errorf("inserting clock: %w", err)
	}

	err = row.Scan(&dbBlockID)
	return
}

func (p *Psql) handleTransaction(dbBlockID int64, transactionHash string) (dbTransactionID int64, err error) {
	if id, found := p.TransactionIDs[transactionHash]; found {
		return id, nil
	}

	row := p.tx.Stmt(preparedStatement.insertTransaction).QueryRow(transactionHash, dbBlockID)
	err = row.Err()
	if err != nil {
		return 0, fmt.Errorf("inserting transaction: %w", err)
	}

	err = row.Scan(&dbTransactionID)
	p.TransactionIDs[transactionHash] = dbTransactionID
	return
}

func (p *Psql) HandleInitializedAccount(dbBlockID int64, initializedAccounts []*pb.InitializedAccount) (err error) {
	for _, initializedAccount := range initializedAccounts {
		dbTransactionID, err := p.handleTransaction(dbBlockID, initializedAccount.TrxHash)
		if err != nil {
			return fmt.Errorf("handling transaction: %w", err)
		}
		_, err = p.tx.Stmt(preparedStatement.insertDerivedAddress).Exec(dbTransactionID, initializedAccount.Owner, initializedAccount.Account)
		if err != nil {
			return fmt.Errorf("trx_hash: %d inserting derived_addresses: %w", dbBlockID, err)
		}
	}
	return nil
}

func (p *Psql) HandleBlocksUndo(lastValidBlockNum uint64) error {
	_, err := p.tx.Exec("DELETE CASCADE FROM solana_tokens.blocks WHERE num > $1", lastValidBlockNum)
	if err != nil {
		return fmt.Errorf("deleting block from %d: %w", lastValidBlockNum, err)
	}
	return nil
}

var NotFound = errors.New("not found")

func (p *Psql) resolveAddress(derivedAddress string) (string, error) {
	resolvedAddress := ""
	rows, err := p.tx.Stmt(preparedStatement.selectDerivedAddress).Query(derivedAddress)
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
	rows, err := p.tx.Stmt(preparedStatement.selectDriver).Query(driverAddress)
	if err != nil {
		return 0, fmt.Errorf("selecting drivers %q : %w", driverAddress, err)
	}

	if rows.Next() {
		err = rows.Scan(&dbDriverID)
		rows.Close()
		return
	}
	rows.Close()

	row := p.tx.Stmt(preparedStatement.insertDriver).QueryRow(driverAddress, dbTransactionID)
	err = row.Err()
	if err != nil {
		return 0, fmt.Errorf("inserting driver: %w", err)
	}

	err = row.Scan(&dbDriverID)
	return
}

func (p *Psql) handleFleet(dbTransactionID int64, fleetAddress string) (dbDriverID int64, err error) {
	rows, err := p.tx.Stmt(preparedStatement.selectFleet).Query(fleetAddress)
	if err != nil {
		return 0, fmt.Errorf("selecting fleets %q: %w", fleetAddress, err)
	}
	if rows.Next() {
		err = rows.Scan(&dbDriverID)
		rows.Close()
		return
	}

	rows.Close()

	row := p.tx.Stmt(preparedStatement.insertFleet).QueryRow(fleetAddress, dbTransactionID)
	err = row.Err()
	if err != nil {
		return 0, fmt.Errorf("inserting driver: %w", err)
	}

	err = row.Scan(&dbDriverID)
	return
}

func (p *Psql) handleFleetDriver(dbTransactionID int64, dbFleetID int64, dbDriverID int64) (dbFleetDriverID int64, err error) {
	rows, err := p.tx.Stmt(preparedStatement.selectFleetDriver).Query(dbFleetID, dbDriverID)
	if err != nil {
		return 0, fmt.Errorf("selecting fleet_drivers: %w", err)
	}

	if rows.Next() {
		err = rows.Scan(&dbDriverID)
		rows.Close()
		return
	}
	rows.Close()

	row := p.tx.Stmt(preparedStatement.insertFleetDriver).QueryRow(dbTransactionID, dbFleetID, dbDriverID)
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

		_, err = p.tx.Stmt(preparedStatement.insertPayment).Exec(mintDbID)
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
		_, err = p.tx.Stmt(preparedStatement.insertPayment).Exec(mintDbID)
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

		_, err = p.tx.Stmt(preparedStatement.insertSlipPayment).Exec(dbTransactionID, fleetMintID, driverMintID)
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
		_, err = p.tx.Stmt(preparedStatement.insertTransfer).Exec(dbTransactionID, transfer.From, transfer.To, transfer.Amount)
		if err != nil {
			fmt.Println("processing transfer: ", transfer.From, transfer.To, transfer.Amount, transfer.TrxHash)
			return fmt.Errorf("inserting transfer: %w", err)
		}
	}
	return nil
}

func (p *Psql) insertMint(dbTransactionID int64, mint *pb.Mint) (dbMintID int64, err error) {
	row := p.tx.Stmt(preparedStatement.insertMint).QueryRow(dbTransactionID, mint.To, mint.Amount)
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

		_, err = p.tx.Stmt(preparedStatement.insertAIPayment).Exec(mintDbID)
		if err != nil {
			return fmt.Errorf("inserting payment: %w", err)
		}
	}

	return nil
}

func (p *Psql) HandleOperationalPayments(dbBlockID int64, payments []*pb.OperationalPayment) error {
	for _, payment := range payments {
		dbTransactionID, err := p.handleTransaction(dbBlockID, payment.Mint.TrxHash)
		if err != nil {
			return fmt.Errorf("inserting transaction: %w", err)
		}

		mintDbID, err := p.insertMint(dbTransactionID, payment.Mint)
		if err != nil {
			return fmt.Errorf("inserting mint: %w", err)
		}

		_, err = p.tx.Stmt(preparedStatement.insertOperationalPayment).Exec(mintDbID)
		if err != nil {
			return fmt.Errorf("inserting payment: %w", err)
		}
	}

	return nil
}

func (p *Psql) HandleRewardPayments(dbBlockID int64, payments []*pb.RewardPayment) error {
	for _, payment := range payments {
		dbTransactionID, err := p.handleTransaction(dbBlockID, payment.Mint.TrxHash)
		if err != nil {
			return fmt.Errorf("inserting transaction: %w", err)
		}

		mintDbID, err := p.insertMint(dbTransactionID, payment.Mint)
		if err != nil {
			return fmt.Errorf("inserting mint: %w", err)
		}

		_, err = p.tx.Stmt(preparedStatement.insertRewardPayment).Exec(mintDbID)
		if err != nil {
			return fmt.Errorf("inserting reward_payments: %w", err)
		}
	}

	return nil
}
func (p *Psql) HandleMapCreate(dbBlockID int64, payments []*pb.MapCreate) error {
	for _, payment := range payments {
		dbTransactionID, err := p.handleTransaction(dbBlockID, payment.Burn.TrxHash)
		if err != nil {
			return fmt.Errorf("inserting transaction: %w", err)
		}

		burnDbID, err := p.insertBurns(dbTransactionID, payment.Burn)
		if err != nil {
			return fmt.Errorf("inserting mint: %w", err)
		}

		_, err = p.tx.Stmt(preparedStatement.insertMapCreate).Exec(burnDbID)
		if err != nil {
			return fmt.Errorf("inserting reward_payments: %w", err)
		}
	}

	return nil
}
func (p *Psql) HandleMapConsumptionReward(dbBlockID int64, payments []*pb.MapConsumptionReward) error {
	for _, payment := range payments {
		dbTransactionID, err := p.handleTransaction(dbBlockID, payment.Mint.TrxHash)
		if err != nil {
			return fmt.Errorf("inserting transaction: %w", err)
		}

		mintDbID, err := p.insertMint(dbTransactionID, payment.Mint)
		if err != nil {
			return fmt.Errorf("inserting mint: %w", err)
		}

		_, err = p.tx.Stmt(preparedStatement.insertMapConsumption).Exec(mintDbID)
		if err != nil {
			return fmt.Errorf("inserting reward_payments: %w", err)
		}
	}

	return nil
}

func (p *Psql) insertBurns(dbTransactionID int64, burn *pb.Burn) (dbMintID int64, err error) {
	row := p.tx.Stmt(preparedStatement.insertBurn).QueryRow(dbTransactionID, burn.From, burn.Amount)
	err = row.Err()
	if err != nil {
		return 0, fmt.Errorf("inserting burn: %w", err)
	}

	err = row.Scan(&dbMintID)
	return
}

func (p *Psql) InsertPrice(timestamp time.Time, price float64) error {
	_, err := preparedStatement.insertPrice.Exec(timestamp, price)
	if err != nil {
		return fmt.Errorf("inserting price: %w", err)
	}
	return nil
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
		_, err = p.tx.Stmt(preparedStatement.insertAIPayment).Exec(dbTransactionID, dbMintID)
		if err != nil {
			return fmt.Errorf("inserting ai payment: %w", err)
		}
	}

	return nil
}

func (p *Psql) StoreCursor(cursor *sink.Cursor) error {
	_, err := p.tx.Stmt(preparedStatement.insertCursor).Exec("hivemapper", cursor.String())
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

func (p *Psql) FetchLastPrice() (timestamp time.Time, price float64, err error) {
	row := p.db.QueryRow("SELECT timestamp, price FROM hivemapper.prices ORDER BY timestamp DESC LIMIT 1")
	if row.Err() != nil {
		return time.Now(), 0, fmt.Errorf("selecting price: %w", row.Err())
	}

	err = row.Scan(&timestamp, &price)
	return timestamp, price, err
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
