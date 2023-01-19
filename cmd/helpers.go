package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
)

func ConfigDir() (configDir string) {
	configDir = filepath.Join(os.Getenv("HOME"), ".config", CONFIG_NAME)

	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		os.MkdirAll(configDir, os.ModePerm)
	}

	return
}

func SystemPkgs(getType bool) (packages []Package) {
	cmd := exec.Command("pacman", "-Qqe")

	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())

		return
	}

	lines := strings.Split(string(stdout), "\n")

	for _, pkg := range lines {
		pkg = strings.TrimSpace(pkg)

		if pkg == "" {
			continue
		}

		var pkgType string

		if getType {
			cmd = exec.Command("pacman", "-Ss", fmt.Sprintf("^%s$", pkg))

			stdout, err := cmd.Output()

			if err != nil {
				if err.Error() != "exit status 1" {
					fmt.Println(err.Error())

					return
				}
			}

			searchResult := strings.TrimSpace(string(stdout))

			if len(searchResult) > 0 {
				pkgType = "Pacman"
			} else {
				pkgType = "AUR"
			}
		}

		pkgsWithoutDep := []string{"composer", "phpmyadmin"}

		var pkgArgs string

		if Contains(pkgsWithoutDep, pkg) {
			pkgArgs = "Sdd"
		} else {
			pkgArgs = "S"
		}

		packages = append(packages, Package{
			Name: pkg,
			Type: pkgType,
			Args: pkgArgs,
		})
	}

	return
}

func FilePkgs() (packages []Package) {
	configFile := filepath.Join(ConfigDir(), "packages.json")

	jsonData, err := os.ReadFile(configFile)

	if err != nil {
		fmt.Println("packages.json does not exist please run `command` first")

		return
	}

	json.Unmarshal(jsonData, &packages)

	return
}

func ComparePkgs() (installPkgs []Package, uninstallPkgs []Package) {
	systemPkgs := SystemPkgs(false)
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

	for _, pkg := range installPkgs {
		fmt.Printf("paru -%s %s\n\n", pkg.Args, pkg.Name)

		prompt := promptui.Prompt{
			Label: "Proceed with installation? [Y/n]",
			Validate: func(input string) error {
				if input != "Y" && input != "n" {
					return errors.New("not a valid option [Y/n]")
				}

				return nil
			},
		}

		confirmResult, _ := prompt.Run()

		if confirmResult != "Y" {
			continue
		}

		// cmd := exec.Command("paru", "-S", "calc")

		// // some command output will be input into stderr
		// // e.g.
		// // cmd := exec.Command("../../bin/master_build")
		// // stderr, err := cmd.StderrPipe()

		// stdout, err := cmd.StdoutPipe()

		// if err != nil {
		// 	fmt.Println(err)
		// }

		// err = cmd.Start()

		// fmt.Println("The command is running")

		// if err != nil {
		// 	fmt.Println(err)
		// }

		// // print the output of the subprocess
		// scanner := bufio.NewScanner(stdout)

		// for scanner.Scan() {
		// 	m := scanner.Text()

		// 	fmt.Println(m)
		// }

		// cmd.Wait()
	}
}

func FilesConfig() (items []File) {
	configFile := filepath.Join(ConfigDir(), "files.json")

	jsonData, err := os.ReadFile(configFile)

	if err != nil {
		fmt.Println("files.json does not exist please make the file first")

		return
	}

	json.Unmarshal(jsonData, &items)

	return
}

func Contains(list []string, str string) bool {
	for _, item := range list {
		if item == str {
			return true
		}
	}

	return false
}
