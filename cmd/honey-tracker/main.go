package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/streamingfast/honey-tracker/web"

	"github.com/spf13/cobra"
	"github.com/streamingfast/bstream"
	"github.com/streamingfast/cli/sflags"
	"github.com/streamingfast/honey-tracker/data"
	"github.com/streamingfast/honey-tracker/price"
	"github.com/streamingfast/logging"
	sink "github.com/streamingfast/substreams-sink"
	"github.com/streamingfast/substreams/client"
	"go.uber.org/zap"
)

var RootCmd = &cobra.Command{
	Use:   "honey-tracker <endpoint> <manifest> <module>",
	Short: "Hivemapper Honey Tracker",
	RunE:  rootRun,
	Args:  cobra.ExactArgs(3),
}

func init() {
	RootCmd.Flags().Bool("insecure", false, "Skip TLS certificate verification")
	RootCmd.Flags().Bool("plaintext", false, "Use plaintext connection")

	// Database
	RootCmd.Flags().String("db-host", "localhost", "PostgreSQL host endpoint")
	RootCmd.Flags().Int("db-port", 5432, "PostgreSQL port")
	RootCmd.Flags().String("db-user", "user", "PostgreSQL user")
	RootCmd.Flags().String("db-password", "secureme", "PostgreSQL password")
	RootCmd.Flags().String("db-name", "postgres", "PostgreSQL database name")
	RootCmd.Flags().Uint64("start-block", 0, "start block number (0 means no start block)")
	RootCmd.Flags().Uint64("stop-block", 0, "stop block number (0 means no stop block)")
	RootCmd.Flags().Duration("startup-delay", time.Duration(0), "stop block number (0 means no stop block)")

	// Manifest
	RootCmd.Flags().String("output-module-type", "proto:hivemapper.types.v1.Output", "Expected output module type")
}

func rootRun(cmd *cobra.Command, args []string) error {
	apiToken := os.Getenv("SUBSTREAMS_API_TOKEN")
	if apiToken == "" {
		return fmt.Errorf("missing SUBSTREAMS_API_TOKEN environment variable")
	}

	if sflags.MustGetDuration(cmd, "startup-delay") != time.Duration(0) {
		time.Sleep(sflags.MustGetDuration(cmd, "startup-delay"))
	}

	logger, tracer := logging.ApplicationLogger("honey-tracker", "honey-tracker")

	endpoint := args[0]
	manifestPath := args[1]
	outputModuleName := args[2]
	expectedOutputModuleType := sflags.MustGetString(cmd, "output-module-type")

	flagInsecure := sflags.MustGetBool(cmd, "insecure")
	flagPlaintext := sflags.MustGetBool(cmd, "plaintext")
	startBlock := sflags.MustGetUint64(cmd, "start-block")
	stopBlock := sflags.MustGetUint64(cmd, "start-block")

	db := data.NewPostgreSQL(
		&data.PsqlInfo{
			Host:     sflags.MustGetString(cmd, "db-host"),
			Port:     sflags.MustGetInt(cmd, "db-port"),
			User:     sflags.MustGetString(cmd, "db-user"),
			Password: sflags.MustGetString(cmd, "db-password"),
			Dbname:   sflags.MustGetString(cmd, "db-name"),
		},
		logger,
	)
	err := db.Init()
	checkError(err)

	clientConfig := client.NewSubstreamsClientConfig(
		endpoint,
		apiToken,
		flagInsecure,
		flagPlaintext,
	)

	pkg, module, outputModuleHash, br, err := sink.ReadManifestAndModuleAndBlockRange(manifestPath, nil, outputModuleName, expectedOutputModuleType, false, "", logger)
	checkError(err)

	options := []sink.Option{
		sink.WithBlockRange(br),
		sink.WithAverageBlockSec("average received block second", 30),
		sink.WithAverageBlockTimeProcessing("average block processing time", 1000),
	}

	if startBlock > 0 && stopBlock > 0 {
		blockRange, err := bstream.NewRangeContaining(startBlock, stopBlock)
		if err != nil {
			return fmt.Errorf("creating block range: %w", err)
		}
		options = append(options, sink.WithBlockRange(blockRange))
	}

	s, err := sink.New(
		sink.SubstreamsModeProduction,
		pkg,
		module,
		outputModuleHash,
		clientConfig,
		logger,
		tracer,
		options...,
	)
	checkError(err)

	ctx := context.Background()
	sinker := data.NewSinker(logger, s, db)
	sinker.OnTerminating(func(err error) {
		logger.Error("sinker terminating", zap.Error(err))
	})

	go func() {
		err := trackPrice(db, logger)
		panic(err)
	}()

	go func() {
		err := runDBT(logger)
		panic(err)
	}()

	go func() {
		web.ServeHttp()
	}()

	err = sinker.Run(ctx)
	if err != nil {
		return fmt.Errorf("runnning sinker:%w", err)
	}

	return nil
}
func main() {
	if err := RootCmd.Execute(); err != nil {
		panic(err)
	}

	fmt.Println("Goodbye!")
}

func runDBT(logger *zap.Logger) error {
	for {
		logger.Info("running dbt")
		cmd := exec.Command("/usr/local/bin/dbt", "run", "--profiles-dir", "/app/hivemapper", "--project-dir", "/app/hivemapper")
		output, err := cmd.CombinedOutput()
		if err != nil {
			logger.Error("running dbt", zap.Error(err), zap.ByteString("output", output))
			return fmt.Errorf("running dbt: %w", err)
		}

		time.Sleep(5 * time.Minute)
	}
}

func trackPrice(db *data.Psql, logger *zap.Logger) error {
	for {
		from, _, err := db.FetchLastPrice()
		if err != nil {
			return fmt.Errorf("fetching last price: %w", err)
		}

		prices := make(chan *price.HistoricalPrice)
		go func() {
			err := price.Fetch(from, time.Now(), prices, logger)
			if err != nil {
				panic(err)
			}
		}()

		for historicalPrice := range prices {
			t := time.Unix(historicalPrice.UnixTime, 0)
			err := db.InsertPrice(t, historicalPrice.Value)
			if err != nil {
				return fmt.Errorf("inserting price: %w", err)
			}
		}
		time.Sleep(2 * time.Minute)
	}
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
