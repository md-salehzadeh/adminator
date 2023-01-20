package cmd

import (
	"fmt"

	"github.com/pterm/pterm"
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
}

func runCompareCmd(cmd *cobra.Command) {
	installPkgs, uninstallPkgs := ComparePkgs()

	pterm.DefaultSection.Println("Install Packages")

	var bullteItems []pterm.BulletListItem

	if len(installPkgs) == 0 {
		pterm.FgLightMagenta.Println("No packages need to be installed")
	} else {
		for _, pkg := range installPkgs {
			bullteItems = append(bullteItems, pterm.BulletListItem{
				Level: 0,
				Text:  fmt.Sprintf("paru -%s %s", pkg.Args, pkg.Name),
			})
		}
	}

	pterm.DefaultBulletList.WithItems(bullteItems).Render()

	pterm.DefaultSection.Println("Uninstall Packages")

	bullteItems = make([]pterm.BulletListItem, 0)

	if len(uninstallPkgs) == 0 {
		pterm.FgLightMagenta.Println("No packages need to be uninstalled")
	} else {
		for _, pkg := range uninstallPkgs {
			bullteItems = append(bullteItems, pterm.BulletListItem{
				Level: 0,
				Text:  fmt.Sprintf("sudo pacman -D --asdeps %s", pkg.Name),
			})
		}
	}

	pterm.DefaultBulletList.WithItems(bullteItems).Render()

	if len(installPkgs) == 0 && len(uninstallPkgs) == 0 {
		return
	}

	pterm.FgLightBlue.Println("\nnow tun -> sudo pacman -Rns $(pacman -Qtdq)")
}
