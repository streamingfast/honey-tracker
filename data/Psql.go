package data

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	pb "github.com/streamingfast/honey-tracker/data/pb/hivemapper/v1"
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

func (p *Psql) HandlePayment(payment *pb.DriverPayment) error {
	//TODO implement me
	panic("implement me")
}

func (p *Psql) HandleSplitPayment(splitPayment *pb.TokenSplittingPayment) error {
	//TODO implement me
	panic("implement me")
}

func (p *Psql) HandleTransfer(transfer *pb.Transfer) error {
	//TODO implement me
	panic("implement me")
}

func (p *Psql) HandleMint(mint *pb.Mint) error {
	//TODO implement me
	panic("implement me")
}

func (p *Psql) HandleBurn(burn *pb.Burn) error {
	//TODO implement me
	panic("implement me")
}

func (p *Psql) HandleTransferChecked(transferCheck *pb.TransferChecked) error {
	//TODO implement me
	panic("implement me")
}

func (p *Psql) HandleMintChecked(transferCheck *pb.MintToChecked) error {
	//TODO implement me
	panic("implement me")
}

func (p *Psql) HandleBurnCheck(transferCheck *pb.BurnChecked) error {
	//TODO implement me
	panic("implement me")
}
