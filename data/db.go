package data

import (
	pb "github.com/streamingfast/honey-tracker/data/pb/hivemapper/v1"
	sink "github.com/streamingfast/substreams-sink"
	pbsubstreams "github.com/streamingfast/substreams/pb/sf/substreams/v1"
)

type DB interface {
	Init() error

	HandleClock(clock *pbsubstreams.Clock) (dbBlockID int64, err error)
	HandleInitializedAccount(dbBlockID int64, initializedAccount []*pb.InitializedAccount) error
	HandleRegularDriverPayments(dbBlockID int64, payments []*pb.RegularDriverPayment) error
	HandleAiPayments(dbBlockID int64, payments []*pb.AiTrainerPayment) error
	HandleSplitPayments(dbBlockID int64, splitPayments []*pb.TokenSplittingPayment) error
	HandleNoneSplitPayments(dbBlockID int64, payments []*pb.NoSplitPayment) error
	HandleTransfers(dbBlockID int64, transfers []*pb.Transfer) error
	HandleMints(dbBlockID int64, mints []*pb.Mint) error
	HandleBurns(dbBlockID int64, burns []*pb.Burn) error
	StoreCursor(cursor *sink.Cursor) error
	FetchCursor() (*sink.Cursor, error)
}
