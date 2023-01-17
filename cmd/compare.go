package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var compareCmd = &cobra.Command{
	Use:   "compare",
	Short: "Compare packages on the system versus packages on the file",
	Long:  `Compare packages on the system versus packages on the file`,
	Run: func(cmd *cobra.Command, args []string) {
		runCompare()
	},
}

func init() {
	pkgsCmd.AddCommand(compareCmd)
}

func runCompare() {
	installPkgs, uninstallPkgs := compare()

	fmt.Println("Install Packages", installPkgs)

	fmt.Println("Uninstall Packages", uninstallPkgs)
}

func compare() (installPkgs []Package, uninstallPkgs []Package) {
	systemPackages := SystemPackages()
	filePackages := FilePackages()

	for _, filePkg := range filePackages {
		isInstalled := false

		for _, sysPkg := range systemPackages {
			if sysPkg.Name == filePkg.Name {
				isInstalled = true
				break
			}
		}

		if !isInstalled {
			installPkgs = append(installPkgs, filePkg)
		}
	}

	for _, sysPkg := range systemPackages {
		shouldUninstall := true

		for _, filePkg := range filePackages {
			if filePkg.Name == sysPkg.Name {
				shouldUninstall = false
				break
			}
		}

		if shouldUninstall {
			uninstallPkgs = append(uninstallPkgs, sysPkg)
		}
	}

	return
}
