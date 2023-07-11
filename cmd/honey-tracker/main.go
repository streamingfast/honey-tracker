package main

import (
	"context"

	"github.com/streamingfast/honey-tracker/data"
	sink "github.com/streamingfast/substreams-sink"
	"github.com/streamingfast/substreams/client"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()

	endpoint := ""
	apiToken := ""
	manifestPath := ""
	outputModuleName := ""
	expectedOutputModuleType := ""

	flagInsecure := false
	flagPlaintext := false

	clientConfig := client.NewSubstreamsClientConfig(
		endpoint,
		apiToken,
		flagInsecure,
		flagPlaintext,
	)

	pkg, module, outputModuleHash, err := sink.ReadManifestAndModule(manifestPath, outputModuleName, expectedOutputModuleType, logger)
	checkError(err)

	s, err := sink.New(
		sink.SubstreamsModeProduction,
		pkg,
		module,
		outputModuleHash,
		clientConfig,
		logger,
		nil,
	)
	checkError(err)

	ctx := context.Background()
	sinker := data.NewSinker(s)
	go func() {
		sinker.Run(ctx)
	}()
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
