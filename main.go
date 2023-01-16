package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func main()  {
	cmd := exec.Command("ls")
	
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	lines := strings.Split(string(stdout), "\n")

	var packages []string

	for _, line := range lines {
		if line == "" {
			continue
		}

		packages = append(packages, line)
	}

	for _, pkg := range packages {
		fmt.Printf("Package: %s\n", pkg)
	}

	fmt.Printf("Length is %d Capacity is %d\n", len(packages), cap(packages))
}