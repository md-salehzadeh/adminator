package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"io/ioutil"
	"path/filepath"
	"encoding/json"

	"github.com/spf13/cobra"
)

type Package struct {
	Name string
	Type string
	Args string
}

// pkgsCmd represents the pkgs command
var pkgsCmd = &cobra.Command{
	Use:   "pkgs",
	Short: "pkgs command",
	Long: "pkgs command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("pkgs called")
	},
}

func init() {
	rootCmd.AddCommand(pkgsCmd)
}

func SystemPackages() (packages []Package) {
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

		packages = append(packages, Package{Name: line, Type: "Pacman", Args: "S"})
	}

	return
}

func FilePackages() (packages []Package) {
	configDir := filepath.Join(os.Getenv("HOME"), ".config", "packages")
	
	configFile := filepath.Join(configDir, "packages.json")

	jsonData, err := ioutil.ReadFile(configFile)

	if err != nil {
		panic(err)
	}

	json.Unmarshal(jsonData, &packages)

	return
}
