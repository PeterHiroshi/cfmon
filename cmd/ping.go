// Added by forge-dispatch e2e test
package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Print pong",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("pong! 🏓🔨 forge-dispatch works!")
		fmt.Println("forge-dispatch e2e: notification test")
		fmt.Println("at:", time.Now().Format(time.RFC3339))
	},
}

func init() {
	rootCmd.AddCommand(pingCmd)
}
