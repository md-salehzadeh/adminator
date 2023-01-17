package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var pkgsCmd = &cobra.Command{
	Use:   "pkgs",
	Short: "pkgs command",
	Long:  "pkgs command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("pkgs called")
	},
}

func init() {
	rootCmd.AddCommand(pkgsCmd)
}
