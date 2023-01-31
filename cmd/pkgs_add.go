package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var pkg string

var pkgsAddCmd = &cobra.Command{
	Use:   "add",
	Short: "add a package to the packages file",
	Long:  `add a package to the packages file`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			pkg = args[0]
		}

		runPkgsAddCmd(cmd)
	},
}

func init() {
	pkgsCmd.AddCommand(pkgsAddCmd)
}

func runPkgsAddCmd(cmd *cobra.Command) {
	pterm.Println()

	if pkg == "" {
		pterm.FgLightRed.Println("please specify package name to add")

		return
	}

	storedPkgs := StoredPkgs()

	doesExist := false

	for _, storedPkg := range storedPkgs {
		if storedPkg.Name == pkg {
			doesExist = true
			break
		}
	}

	if doesExist {
		pterm.FgLightRed.Println("package you specified already exists in the packages file")

		return
	}

	pkgEntity := PkgEntity(pkg, true)

	packages := append(storedPkgs, pkgEntity)

	sort.Slice(packages, func(i, j int) bool {
		return packages[i].Name < packages[j].Name
	})

	writeBackToFile(packages)

	pterm.Success.Println("package successfully added to the packages file")
}

func writeBackToFile(packages []Package) {
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
