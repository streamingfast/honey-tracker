package data

import (
	"github.com/streamingfast/honey-tracker/data/proto"
)

type DB interface {
	Init()
	HandlePayment(payment *proto.Payment) error
	HandleSplitPayment(splitPayment *proto.SplitPayment) error
	HandleTransfer(transfer *proto.Transfer) error
	HandleMint(mint *proto.Mint) error
	HandleBurn(burn *proto.Mint) error
	TransferChecked(transferCheck *proto.TransferCheck) error
	//todo: what else?
	//todo: MintToChecked
	//todo: BurnChecked

	//todo: double check https://github.com/solana-labs/solana-program-library/blob/master/token/program/src/instruction.rs
}
