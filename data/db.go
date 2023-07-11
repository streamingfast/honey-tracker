package data

import (
	"github.com/streamingfast/honey-tracker/data/proto"
)

type DB interface {
	Init()
	HandlePayment(payment *proto.Payment) error
	HandleSplitPayment(splitPayment *proto.SplitPayment) error
	HandleTransfer(transfer *proto.Transfer) error
	HandleMint(transfer *proto.Mint) error
	HandleBurn(transfer *proto.Mint) error
	//todo: what else?
}
