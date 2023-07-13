package main

import (
	"github.com/streamingfast/honey-tracker/data"
)

func main() {
	//logger, _ := zap.NewProduction()
	//
	//endpoint := ""
	//apiToken := ""
	//manifestPath := ""
	//outputModuleName := ""
	//expectedOutputModuleType := ""
	//
	//flagInsecure := false
	//flagPlaintext := false

	db := data.NewPostgreSQL(&data.PsqlInfo{
		Host:     "localhost",
		Port:     5432,
		User:     "cbillett",
		Password: "secureme",
		Dbname:   "hivemapper",
	})
	err := db.Init()
	checkError(err)

	//clientConfig := client.NewSubstreamsClientConfig(
	//	endpoint,
	//	apiToken,
	//	flagInsecure,
	//	flagPlaintext,
	//)
	//
	//pkg, module, outputModuleHash, err := sink.ReadManifestAndModule(manifestPath, outputModuleName, expectedOutputModuleType, logger)
	//checkError(err)
	//
	//s, err := sink.New(
	//	sink.SubstreamsModeProduction,
	//	pkg,
	//	module,
	//	outputModuleHash,
	//	clientConfig,
	//	logger,
	//	nil,
	//)
	//checkError(err)

	//ctx := context.Background()
	//sinker := data.NewSinker(s, db)
	//go func() {
	//	sinker.Run(ctx)
	//}()
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
