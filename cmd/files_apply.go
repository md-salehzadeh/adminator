package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var filesApplyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply command",
	Long:  `Apply command`,
	Run: func(cmd *cobra.Command, args []string) {
		runFilesApplyCmd(cmd)
	},
}

func init() {
	filesCmd.AddCommand(filesApplyCmd)
}

func runFilesApplyCmd(cmd *cobra.Command) {
	configDir := ConfigDir()

	filesConfig := FilesConfig()

	for _, item := range filesConfig {
		needsSudo := true

		if strings.HasPrefix(item.Path, "~/") {
			item.Path = filepath.Join(os.Getenv("HOME"), strings.Replace(item.Path, "~/", "", 1))

			needsSudo = false
		}

		if item.Type == "create" {
			item.createFile(needsSudo)
		} else if item.Type == "link" {
			item.linkFile(configDir, needsSudo)
		}
	}
}

func (item File) createFile(needsSudo bool) {
	if _, err := os.Stat(item.Path); os.IsNotExist(err) {
		var cmd *exec.Cmd

		if needsSudo {
			cmd = exec.Command("sudo", "mkdir", item.Path)
		} else {
			cmd = exec.Command("mkdir", item.Path)
		}

		cmd.Run()

		if _, err := os.Stat(item.Path); os.IsNotExist(err) {
			fmt.Printf("couldn't create directory '%s'", item.Path)
		}
	}
}

func (item File) linkFile(configDir string, needsSudo bool) {
	source := filepath.Join(configDir, "files", item.Path)

	createLink := false

	if info, err := os.Lstat(item.Path); os.IsNotExist(err) {
		createLink = true
	} else if err != nil {
		fmt.Printf("Error checking if %s exists: %v\n", item.Path, err)
	} else if info.Mode()&os.ModeSymlink != 0 {
		var cmd *exec.Cmd

		if needsSudo {
			cmd = exec.Command("sudo", "rm", "-f", item.Path)
		} else {
			cmd = exec.Command("rm", "-f", item.Path)
		}

		cmd.Run()

		createLink = true
	}

	if _, err := os.Stat(source); !os.IsNotExist(err) && createLink {
		var cmd *exec.Cmd

		if needsSudo {
			cmd = exec.Command("sudo", "ln", "-s", source, item.Path)
		} else {
			cmd = exec.Command("ln", "-s", source, item.Path)
		}

		cmd.Run()
	}
}
