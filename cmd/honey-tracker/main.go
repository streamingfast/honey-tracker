package main

import (
	"context"
	"os"

	"go.uber.org/zap"

	"github.com/streamingfast/honey-tracker/data"
	"github.com/streamingfast/logging"
	sink "github.com/streamingfast/substreams-sink"
	"github.com/streamingfast/substreams/client"
)

func main() {
	apiToken := os.Getenv("SUBSTREAMS_API_TOKEN")
	if apiToken == "" {
		panic("Missing SUBSTREAMS_API_TOKEN environment variable")
	}

	//logger, _ := zap.NewProduction()
	logger, tracer := logging.ApplicationLogger("honey-tracker", "honey-tracker")

	endpoint := "mainnet.sol.streamingfast.io:443"
	manifestPath := "/Users/cbillett/devel/sf/substreams-hivemapper/substreams.yaml"
	outputModuleName := "map_outputs"
	expectedOutputModuleType := "proto:hivemapper.types.v1.Output"

	flagInsecure := false
	flagPlaintext := false

	db := data.NewPostgreSQL(&data.PsqlInfo{
		Host:     "localhost",
		Port:     5432,
		User:     "cbillett",
		Password: "secureme",
		Dbname:   "hivemapper",
	})
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
		logger,
		tracer,
		sink.WithBlockRange(br),
	)
	checkError(err)

	ctx := context.Background()
	sinker := data.NewSinker(s, db)
	sinker.OnTerminating(func(err error) {
		logger.Error("sinker terminating", zap.Error(err))
	})
	sinker.Run(ctx)
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
