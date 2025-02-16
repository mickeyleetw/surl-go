package servercmd

import (
	server "shorten_url/pkg/core/server"

	"github.com/spf13/cobra"
)

// RunServerCmd is the command to run the server
var (
	RunServerCmd = &cobra.Command{
		Use:   "run",
		Short: "Server commands",
		Long:  "Server commands",
		Run: func(cmd *cobra.Command, args []string) {
			runServer()
		},
	}
)

func runServer() {
	server.InitServer()
}
