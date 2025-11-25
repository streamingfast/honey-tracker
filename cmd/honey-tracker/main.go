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
	Args:  cobra.ExactArgs(3),
}

func init() {
}

func rootRun(cmd *cobra.Command, args []string) error {
	server := &web.Server{}
	go func() {
		server.ServeHttp()
	}()

	return nil
}
func main() {
	if err := RootCmd.Execute(); err != nil {
		panic(err)
	}

	fmt.Println("Goodbye!")
}
