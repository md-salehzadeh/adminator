package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func ConfigDir() (configDir string) {
	configDir = filepath.Join(os.Getenv("HOME"), ".config", CONFIG_NAME)

	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		os.MkdirAll(configDir, os.ModePerm)
	}

	return
}

func SystemPkgs() (packages []Package) {
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

		packages = append(packages, Package{
			Name: line,
			Type: "Pacman",
			Args: "S",
		})
	}

	return
}

func FilePkgs() (packages []Package) {
	configFile := filepath.Join(ConfigDir(), "packages.json")

	jsonData, err := ioutil.ReadFile(configFile)

	if err != nil {
		fmt.Println("packages.json does not exist please run `command` first")

		return
	}

	json.Unmarshal(jsonData, &packages)

	return
}

func ComparePkgs() (installPkgs []Package, uninstallPkgs []Package) {
	systemPkgs := SystemPkgs()
	filePkgs := FilePkgs()

	for _, filePkg := range filePkgs {
		isInstalled := false

		for _, sysPkg := range systemPkgs {
			if sysPkg.Name == filePkg.Name {
				isInstalled = true
				break
			}
		}

		if !isInstalled {
			installPkgs = append(installPkgs, filePkg)
		}
	}

	for _, sysPkg := range systemPkgs {
		shouldUninstall := true

		for _, filePkg := range filePkgs {
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

func ApplyPkgs(installPkgs []Package, uninstallPkgs []Package) {
	fmt.Println("\nnow applying changes")
}
