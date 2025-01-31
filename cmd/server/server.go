package serverCmd

import (
	server "shorten_url/pkg/core/server"

	"github.com/spf13/cobra"
)

var (
	Cmd = &cobra.Command{
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
