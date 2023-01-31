package cmd

import (
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

	writeToFile(packages)

	pterm.Success.Println("package successfully added to the packages file")
}
