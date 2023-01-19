package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var servicesApplyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply command",
	Long:  `Apply command`,
	Run: func(cmd *cobra.Command, args []string) {
		runServicesApplyCmd(cmd)
	},
}

func init() {
	servicesCmd.AddCommand(servicesApplyCmd)
}

func runServicesApplyCmd(cmd *cobra.Command) {
	servicesConfig := ServicesConfig()

	for _, item := range servicesConfig {
		_, err := item.getStatus()

		if err != nil {
			fmt.Println(err)

			continue
		}

		var cmd *exec.Cmd

		if item.Enabled {
			cmd = exec.Command("sudo", "systemctl", "enable", item.Name)
		} else {
			cmd = exec.Command("sudo", "systemctl", "disable", item.Name)
		}

		cmd.Run()

		// fmt.Printf("Service %s status: %s\n", item.Name, status)
	}
}

func (item Service) getStatus() (status string, error error) {
	cmd := exec.Command("systemctl", "status", item.Name)

	cmdOutput, err := cmd.Output()

	output := strings.TrimSpace(string(cmdOutput))

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if strings.Contains(output, "Loaded: not-found") {
				status = "not found"
			} else if strings.Contains(output, "Active: inactive") {
				status = "inactive"
			} else if strings.Contains(output, "Active: failed") {
				status = "failed"
			} else {
				error = fmt.Errorf("error: %s", exitError.Stderr)
			}
		} else {
			error = fmt.Errorf("error: %s", err.Error())
		}
	} else {
		if strings.Contains(output, "Active: active") {
			status = "active"
		} else {
			status = output
		}
	}

	return
}
