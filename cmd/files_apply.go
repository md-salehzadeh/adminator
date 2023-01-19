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
	if item.Source == "" {
		item.Source = filepath.Join(configDir, "files", item.Path)
	}

	if _, err := os.Stat(item.Source); os.IsNotExist(err) {
		fmt.Printf("Source path '%s' doesn't exists\n", item.Source)

		return
	}

	var cmd *exec.Cmd

	if needsSudo {
		cmd = exec.Command("sudo", "rm", "-f", item.Path)
	} else {
		cmd = exec.Command("rm", "-f", item.Path)
	}

	cmd.Run()

	if needsSudo {
		cmd = exec.Command("sudo", "ln", "-s", item.Source, item.Path)
	} else {
		cmd = exec.Command("ln", "-s", item.Source, item.Path)
	}

	cmd.Run()
}
