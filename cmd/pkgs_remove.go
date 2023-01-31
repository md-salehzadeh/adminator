package cmd

import (
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var pkgsRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "remove a package from the packages file",
	Long:  `remove a package from the packages file`,
	Run: func(cmd *cobra.Command, args []string) {
		var pkg string

		if len(args) > 0 {
			pkg = args[0]
		}

		runPkgsRemoveCmd(cmd, pkg)
	},
}

func init() {
	pkgsCmd.AddCommand(pkgsRemoveCmd)
}

func runPkgsRemoveCmd(cmd *cobra.Command, pkg string) {
	pterm.Println()

	if pkg == "" {
		pterm.FgLightRed.Println("please specify package name to remove")

		return
	}

	storedPkgs := StoredPkgs()

	doesExist := false
	var packages []Package

	for _, storedPkg := range storedPkgs {
		if storedPkg.Name == pkg {
			doesExist = true
		} else {
			packages = append(packages, storedPkg)
		}
	}

	if !doesExist {
		pterm.FgLightRed.Println("package you specified does not exist in the packages file")

		return
	}

	writeToFile(packages)

	pterm.Success.Println("package successfully removed from the packages file")
}
