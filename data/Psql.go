package data

import (
	"database/sql"
	"fmt"

	"github.com/streamingfast/honey-tracker/data/proto"
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
	_, err := p.db.Exec(fleetsCreateTable)
	if err != nil {
		return fmt.Errorf("creating fleets table: %w", err)
	}

	_, err = p.db.Exec(driversCreateTable)
	if err != nil {
		return fmt.Errorf("creating drivers table: %w", err)
	}

	_, err = p.db.Exec(fleetDriversCreateTable)
	if err != nil {
		return fmt.Errorf("creating fleet_drivers table: %w", err)
	}

	_, err = p.db.Exec(paymentsCreateTable)
	if err != nil {
		return fmt.Errorf("creating payments table: %w", err)
	}

	_, err = p.db.Exec(splitPaymentsCreateTable)
	if err != nil {
		return fmt.Errorf("creating split_payments table: %w", err)
	}

	_, err = p.db.Exec(transfersCreateTable)
	if err != nil {
		return fmt.Errorf("creating transfers table: %w", err)
	}

	//todo: missing MINT
	//todo: missing BURN

	return nil
}

func (p *Psql) HandlePayment(payment proto.Payment) error {
	//TODO implement me
	panic("implement me")
}

func (p *Psql) HandleSplitPayment(splitPayment proto.SplitPayment) error {
	//TODO implement me
	panic("implement me")
}

func (p *Psql) HandleTransfer(transfer proto.Transfer) error {
	//TODO implement me
	panic("implement me")
}
