package data

import (
	"database/sql"
	"fmt"
	"honey-tracker/data/proto"
)

type DB interface {
	Init()
	HandlePayment(payment proto.Payment) error
	HandleSplitPayment(splitPayment proto.SplitPayment) error
	HandleTransfer(transfer proto.Transfer) error
}

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

func (p *Psql) Init() {
	//TODO implement me
	panic("implement me")
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
