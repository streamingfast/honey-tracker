package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/streamingfast/cli/sflags"

	"go.uber.org/zap"

	"github.com/streamingfast/honey-tracker/data"
	"github.com/streamingfast/logging"
	sink "github.com/streamingfast/substreams-sink"
	"github.com/streamingfast/substreams/client"
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

	// Manifest
	RootCmd.Flags().String("output-module-type", "proto:hivemapper.types.v1.Output", "Expected output module type")
}

func rootRun(cmd *cobra.Command, args []string) error {
	apiToken := os.Getenv("SUBSTREAMS_API_TOKEN")
	if apiToken == "" {
		return fmt.Errorf("missing SUBSTREAMS_API_TOKEN environment variable")
	}

	logger, tracer := logging.ApplicationLogger("honey-tracker", "honey-tracker")

	endpoint := args[0]
	manifestPath := args[1]
	outputModuleName := args[2]
	expectedOutputModuleType := sflags.MustGetString(cmd, "output-module-type")

	flagInsecure := sflags.MustGetBool(cmd, "insecure")
	flagPlaintext := sflags.MustGetBool(cmd, "plaintext")

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

	s, err := sink.New(
		sink.SubstreamsModeProduction,
		pkg,
		module,
		outputModuleHash,
		clientConfig,
		zap.NewNop(),
		tracer,
		sink.WithBlockRange(br),
	)
	checkError(err)

	ctx := context.Background()
	sinker := data.NewSinker(logger, s, db)
	sinker.OnTerminating(func(err error) {
		logger.Error("sinker terminating", zap.Error(err))
	})
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

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
