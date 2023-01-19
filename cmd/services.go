package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var servicesCmd = &cobra.Command{
	Use:   "services",
	Short: "services command",
	Long:  "services command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("services called")
	},
}

func init() {
	rootCmd.AddCommand(servicesCmd)
}
