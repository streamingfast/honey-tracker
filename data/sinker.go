package data

import (
	"context"
	"fmt"

	sink "github.com/streamingfast/substreams-sink"
	pbsubstreamsrpc "github.com/streamingfast/substreams/pb/sf/substreams/rpc/v2"
)

type Sinker struct {
	*sink.Sinker
	db *DB
}

func NewSinker(sink *sink.Sinker, db *DB) *Sinker {
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

	//todo: unmarshal proto data

	//todo: switch data := data.Data.(type)
	//todo then call db methods

	//todo: save cursor

	panic("implement me")
}

func (s *Sinker) HandleBlockUndoSignal(ctx context.Context, undoSignal *pbsubstreamsrpc.BlockUndoSignal, cursor *sink.Cursor) error {
	panic("should not called on solana")
}
