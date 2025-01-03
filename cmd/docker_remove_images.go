package cmd

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

func dockerRemoveImagesCmd(cmd *cobra.Command, args []string) error {

	_, errViper := initViper(cmd)
	if errViper != nil {
		return fmt.Errorf("error initializing viper: %w", errViper)
	}

	cmdList := exec.Command("docker", "ps", "-a", "-f", "status=exited", "-f", "status=created", "-q")
	var containerIDs bytes.Buffer
	cmdList.Stdout = &containerIDs
	cmdList.Stderr = &containerIDs

	if err := cmdList.Run(); err != nil {
		return fmt.Errorf("Failed to list containers: %v\nOutput: %s\n", err, containerIDs.String())
	}

	ids := strings.Fields(containerIDs.String())
	if len(ids) == 0 {
		return fmt.Errorf("No containers found with 'exited' or 'created' status.")
	}

	dockerArgs := append([]string{"rm"}, ids...)
	cmdRemove := exec.Command("docker", dockerArgs...)
	var output bytes.Buffer
	cmdRemove.Stdout = &output
	cmdRemove.Stderr = &output

	if err := cmdRemove.Run(); err != nil {
		return fmt.Errorf("Failed to remove containers: %v\nOutput: %s\n", err, output.String())
	}

	fmt.Println("Successfully removed the following containers:")
	fmt.Println(output.String())

	return nil
}
