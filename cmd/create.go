package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create the packages file",
	Long:  `Create the packages file`,
	Run: func(cmd *cobra.Command, args []string) {
		runCreate()
	},
}

func init() {
	pkgsCmd.AddCommand(createCmd)
}

func runCreate() {
	packages := getPackages()

	fmt.Printf("Length is %d Capacity is %d\n", len(packages), cap(packages))
}

func getPackages() (packages []string) {
	cmd := exec.Command("pacman", "-Qqe")

	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	lines := strings.Split(string(stdout), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		packages = append(packages, line)
	}

	return
}
