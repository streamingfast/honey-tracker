package data

import (
	pb "github.com/streamingfast/honey-tracker/data/pb/hivemapper/v1"
	pbsubstreams "github.com/streamingfast/substreams/pb/sf/substreams/v1"
)

type DB interface {
	Init() error

	HandleClock(clock *pbsubstreams.Clock) (dbBlockID int64, err error)
	HandleInitializedAccount(dbBlockID int64, initializedAccount []*pb.InitializedAccount) error
	HandlePayments(dbBlockID int64, payments []*pb.Payment) error
	HandleSplitPayments(dbBlockID int64, splitPayments []*pb.TokenSplittingPayment) error
	HandleTransfers(dbBlockID int64, transfers []*pb.Transfer) error
	HandleMints(dbBlockID int64, mints []*pb.Mint) error
	HandleBurns(dbBlockID int64, burns []*pb.Burn) error
	HandleTransferCheckeds(dbBlockID int64, transferChecks []*pb.TransferChecked) error
	HandleMintCheckeds(dbBlockID int64, transferChecks []*pb.MintToChecked) error
	HandleBurnChecks(dbBlockID int64, transferChecks []*pb.BurnChecked) error
	//todo: double check https://github.com/solana-labs/solana-program-library/blob/master/token/program/src/instruction.rs
}
