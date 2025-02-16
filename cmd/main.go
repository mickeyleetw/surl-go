package main

import (
	dbCmd "shorten_url/cmd/db"
	serverCmd "shorten_url/cmd/server"

	"log"

	"github.com/spf13/cobra"
)

// rootCmd is the root command for the application and & stands for the new command object pointer
var rootCmd = &cobra.Command{
	Use:   "sUrl",
	Short: "Shorten URL",
	Long:  "Shorten URL",
}

func init() {
	rootCmd.AddCommand(dbCmd.ResetDBCmd)
	rootCmd.AddCommand(dbCmd.ResetRedisCmd)
	rootCmd.AddCommand(serverCmd.RunServerCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Failed to execute root command: %v", err)
	}
}
