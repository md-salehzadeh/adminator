package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create the packages file",
	Long:  `Create the packages file`,
	Run: func(cmd *cobra.Command, args []string) {
		runCreateCmd()
	},
}

func init() {
	pkgsCmd.AddCommand(createCmd)
}

func runCreateCmd() {
	packages := SystemPkgs(true)

	fmt.Printf("Number of %d packages stored in '~/.config/%s/packages.json'", len(packages), CONFIG_NAME)

	writeToFile(packages)
}

func writeToFile(packages []Package) {
	jsonData, err := json.MarshalIndent(packages, "", "    ")

	if err != nil {
		panic(err)
	}

	configFile := filepath.Join(ConfigDir(), "packages.json")

	err = os.WriteFile(configFile, jsonData, 0644)

	if err != nil {
		panic(err)
	}
}
