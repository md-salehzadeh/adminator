package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var filesCmd = &cobra.Command{
	Use:   "files",
	Short: "files command",
	Long:  "files command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("files called")
	},
}

func init() {
	rootCmd.AddCommand(filesCmd)
}
