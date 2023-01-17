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
		runCompareCmd(cmd)
	},
}

func init() {
	pkgsCmd.AddCommand(compareCmd)

	compareCmd.Flags().BoolP("apply", "a", false, "whether to apply the changes or not")
}

func runCompareCmd(cmd *cobra.Command) {
	installPkgs, uninstallPkgs := ComparePkgs()

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

	if len(installPkgs) == 0 && len(uninstallPkgs) == 0 {
		return
	}

	apply, _ := cmd.Flags().GetBool("apply")

	if apply {
		ApplyPkgs(installPkgs, uninstallPkgs)
	}
}
