package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"encoding/json"
	"io/ioutil"
	"path/filepath"

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
	fmt.Println(packages)

	writeToFile(packages)
}

type Package struct {
	Name string
	Type string
	Args string
}

func getPackages() (packages []Package) {
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

func writeToFile(packages []Package) {
	jsonData, err := json.MarshalIndent(packages, "", "    ")

	if err != nil {
		panic(err)
	}

	configDir := filepath.Join(os.Getenv("HOME"), ".config", "packages")
	
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		os.MkdirAll(configDir, os.ModePerm)
	}

	configFile := filepath.Join(configDir, "packages.json")

	err = ioutil.WriteFile(configFile, jsonData, 0644)

	if err != nil {
		panic(err)
	}
}