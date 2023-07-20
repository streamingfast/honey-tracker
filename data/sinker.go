package data

import (
	"context"
	"fmt"

	pb "github.com/streamingfast/honey-tracker/data/pb/hivemapper/v1"
	sink "github.com/streamingfast/substreams-sink"
	pbsubstreamsrpc "github.com/streamingfast/substreams/pb/sf/substreams/rpc/v2"
	"google.golang.org/protobuf/proto"
)

type Sinker struct {
	*sink.Sinker
	db DB
}

func NewSinker(sink *sink.Sinker, db DB) *Sinker {
	return &Sinker{
		Sinker: sink,
		db:     db,
	}
}

func (s *Sinker) Run(ctx context.Context) {
	//todo: get cursor
	var cursor *sink.Cursor

	s.Sinker.Run(ctx, cursor, s)
}

func (s *Sinker) HandleBlockScopedData(ctx context.Context, data *pbsubstreamsrpc.BlockScopedData, isLive *bool, cursor *sink.Cursor) error {
	output := data.Output
	if output.Name != s.OutputModuleName() {
		return fmt.Errorf("received data from wrong output module, expected to received from %q but got module's output for %q", s.OutputModuleName(), output.Name)
	}

	//todo: start db transaction

	//todo: handle block clock
	dbBlockID, err := s.db.HandleClock(data.Clock)
	if err != nil {
		return fmt.Errorf("handle block clock: %w", err)
	}
	moduleOutput := &pb.Output{}
	err = proto.Unmarshal(output.GetMapOutput().GetValue(), moduleOutput)
	if err != nil {
		return fmt.Errorf("unmarshal module output changes: %w", err)
	}

	if err := s.db.HandlePayments(dbBlockID, moduleOutput.Payments); err != nil {
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

	//todo: save cursor
	//todo: commit transaction
	return nil
}

func (s *Sinker) HandleBlockUndoSignal(ctx context.Context, undoSignal *pbsubstreamsrpc.BlockUndoSignal, cursor *sink.Cursor) error {
	panic("should not called on solana")
}
