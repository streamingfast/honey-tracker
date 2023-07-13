package data

import (
	pb "github.com/streamingfast/honey-tracker/data/pb/hivemapper/v1"
)

type DB interface {
	Init() error
	HandlePayment(payment *pb.DriverPayment) error
	HandleSplitPayment(splitPayment *pb.TokenSplittingPayment) error
	HandleTransfer(transfer *pb.Transfer) error
	HandleMint(mint *pb.Mint) error
	HandleBurn(burn *pb.Burn) error
	HandleTransferChecked(transferCheck *pb.TransferChecked) error
	HandleMintChecked(transferCheck *pb.MintToChecked) error
	HandleBurnCheck(transferCheck *pb.BurnChecked) error
	//todo: double check https://github.com/solana-labs/solana-program-library/blob/master/token/program/src/instruction.rs
}
