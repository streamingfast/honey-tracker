package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/streamingfast/honey-tracker/web"
)

var RootCmd = &cobra.Command{
	Use:   "honey-tracker",
	Short: "Hivemapper Honey Tracker",
	RunE:  rootRun,
}

func init() {
}

func rootRun(cmd *cobra.Command, args []string) error {
	server := &web.Server{}
	server.ServeHttp()

	return nil
}
func main() {
	if err := RootCmd.Execute(); err != nil {
		panic(err)
	}

	fmt.Println("Goodbye!")
}
