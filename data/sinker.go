package data

import (
	"context"
	"fmt"

	pb "github.com/streamingfast/honey-tracker/data/pb/hivemapper/v1"
	sink "github.com/streamingfast/substreams-sink"
	pbsubstreamsrpc "github.com/streamingfast/substreams/pb/sf/substreams/rpc/v2"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

type Sinker struct {
	logger *zap.Logger
	*sink.Sinker
	db DB
}

func NewSinker(logger *zap.Logger, sink *sink.Sinker, db DB) *Sinker {
	return &Sinker{
		logger: logger,
		Sinker: sink,
		db:     db,
	}
}

func (s *Sinker) Run(ctx context.Context) error {
	//todo: get cursor
	//var cursor *sink.Cursor

	cursor, err := s.db.FetchCursor()
	if err != nil {
		return fmt.Errorf("fetch cursor: %w", err)
	}
	s.Sinker.Run(ctx, cursor, s)
	return nil
}

func (s *Sinker) HandleBlockScopedData(ctx context.Context, data *pbsubstreamsrpc.BlockScopedData, isLive *bool, cursor *sink.Cursor) error {
	output := data.Output
	if output.Name != s.OutputModuleName() {
		return fmt.Errorf("received data from wrong output module, expected to received from %q but got module's output for %q", s.OutputModuleName(), output.Name)
	}

	if len(output.GetMapOutput().GetValue()) == 0 {
		if data.Clock.Number%100 == 0 {
			s.logger.Info("progress_block", zap.Uint64("block", data.Clock.Number))
		}
		return s.db.StoreCursor(cursor)
	}

	moduleOutput := &pb.Output{}
	err := proto.Unmarshal(output.GetMapOutput().GetValue(), moduleOutput)
	if err != nil {
		return fmt.Errorf("unmarshal module output changes: %w", err)
	}

	dbBlockID, err := s.db.HandleClock(data.Clock)
	if err != nil {
		return fmt.Errorf("handle block clock: %w", err)
	}

	err = s.db.HandleInitializedAccount(dbBlockID, moduleOutput.InitializedAccount)
	if err != nil {
		return fmt.Errorf("handle initialized accounts: %w", err)
	}

	if err := s.db.HandleRegularDriverPayments(dbBlockID, moduleOutput.RegularDriverPayments); err != nil {
		return fmt.Errorf("handle payments: %w", err)
	}

	if err := s.db.HandleSplitPayments(dbBlockID, moduleOutput.TokenSplittingPayments); err != nil {
		return fmt.Errorf("handle split payments: %w", err)
	}

	if err := s.db.HandleTransfers(dbBlockID, moduleOutput.Transfers); err != nil {
		return fmt.Errorf("handle transfers: %w", err)
	}

	if err := s.db.HandleMints(dbBlockID, moduleOutput.Mints); err != nil {
		return fmt.Errorf("handle mints: %w", err)
	}

	if err := s.db.HandleBurns(dbBlockID, moduleOutput.Burns); err != nil {
		return fmt.Errorf("handle burns: %w", err)
	}

	err = s.db.StoreCursor(cursor)
	if err != nil {
		return fmt.Errorf("store cursor: %w", err)
	}
	//todo: commit transaction
	return nil
}

func (s *Sinker) HandleBlockUndoSignal(ctx context.Context, undoSignal *pbsubstreamsrpc.BlockUndoSignal, cursor *sink.Cursor) error {
	panic("should not be called on solana")
}
