package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"

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

	var wg sync.WaitGroup

	var mutex = &sync.Mutex{}

	for _, pkg := range lines {
		wg.Add(1)

		go func(pkg string) {
			defer wg.Done()

			pkg = strings.TrimSpace(pkg)

			if pkg == "" {
				return
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

			mutex.Lock()

			packages = append(packages, Package{
				Name: pkg,
				Type: pkgType,
				Args: pkgArgs,
			})

			mutex.Unlock()
		}(pkg)
	}

	wg.Wait()

	sort.Slice(packages, func(i, j int) bool {
		return packages[i].Name < packages[j].Name
	})

	return
}

func StoredPkgs() (packages []Package) {
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
	storedPkgs := StoredPkgs()

	for _, storedPkg := range storedPkgs {
		isInstalled := false

		for _, sysPkg := range systemPkgs {
			if sysPkg.Name == storedPkg.Name {
				isInstalled = true
				break
			}
		}

		if !isInstalled {
			installPkgs = append(installPkgs, storedPkg)
		}
	}

	for _, sysPkg := range systemPkgs {
		shouldUninstall := true

		for _, storedPkg := range storedPkgs {
			if storedPkg.Name == sysPkg.Name {
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
		fmt.Println("files.json does not exist please make the config file first")

		return
	}

	json.Unmarshal(jsonData, &items)

	return
}

func ServicesConfig() (items []Service) {
	configFile := filepath.Join(ConfigDir(), "services.json")

	jsonData, err := os.ReadFile(configFile)

	if err != nil {
		fmt.Println("services.json does not exist please make the config file first")

		return
	}

	json.Unmarshal(jsonData, &items)

	return
}

func IsPkgDependency(packageName string) bool {
	cmd := exec.Command("pacman", "-Qi", packageName)

	output, err := cmd.Output()

	if err != nil {
		return false
	}

	outputStr := string(output)

	return strings.Contains(outputStr, "Depends On")
}

func Contains(list []string, str string) bool {
	for _, item := range list {
		if item == str {
			return true
		}
	}

	return false
}
