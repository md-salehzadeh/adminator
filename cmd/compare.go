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
		runCompare(cmd)
	},
}

func init() {
	pkgsCmd.AddCommand(compareCmd)

	compareCmd.Flags().BoolP("apply", "a", false, "whether to apply the changes or not")
}

func runCompare(cmd *cobra.Command) {
	installPkgs, uninstallPkgs := compare()

	fmt.Println("----------------------")

	if len(installPkgs) == 0 {
		fmt.Println("No packages need to be installed")
	} else {
		fmt.Println("Install Packages:")

		for _, pkg := range installPkgs {
			fmt.Printf("\n%s", pkg.Name)
		}
	}

	fmt.Println("\n----------------------")

	if len(uninstallPkgs) == 0 {
		fmt.Println("No packages need to be uninstalled")
	} else {
		fmt.Println("Uninstall Packages:")

		for _, pkg := range uninstallPkgs {
			fmt.Printf("\n%s", pkg.Name)
		}
	}

	fmt.Println("\n----------------------")

	apply, err := cmd.Flags().GetBool("apply")

	if err != nil {
		panic(err)
	}

	if (apply) {
		fmt.Println("\nnow applying changes")
	}
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
