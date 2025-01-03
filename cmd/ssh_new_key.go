package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	sshFlagName      = "name"
	sshFlagNameShort = "n"
)

func initSSHNewKeyFlags(flag *pflag.FlagSet) {
	flag.StringP(sshFlagName, sshFlagNameShort, "", "Set the name of the ssh key")
}

func validateSSHNewKeyFlags(v *viper.Viper) error {
	if v.GetString(sshFlagName) == "" {
		return fmt.Errorf("name cannot be empty")
	}
	return nil
}

func sshNewKeyCmd(cmd *cobra.Command, args []string) error {

	v, errViper := initViper(cmd)
	if errViper != nil {
		return fmt.Errorf("error initializing viper: %w", errViper)
	}

	err := validateSSHNewKeyFlags(v)
	if err != nil {
		return fmt.Errorf("error validating flags: %w", err)
	}

	name := v.GetString(sshFlagName)

	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		return fmt.Errorf("HOME environment variable not set")
	}

	// Construct SSH key path
	sshKey := fmt.Sprintf("%s/.ssh/%s_%s_ed25519", homeDir, name, os.Getenv("USER"))

	// Check if the file already exists
	if _, err := os.Stat(sshKey); os.IsNotExist(err) {
		// Generate the SSH key if it doesn't exist
		cmd := exec.Command("ssh-keygen", "-t", "ed25519", "-o", "-a", "100", "-f", sshKey, "-C", fmt.Sprintf("%s's %s key", os.Getenv("USER"), name))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		// Run the command to generate the key
		if err := cmd.Run(); err != nil {
			fmt.Printf("Error generating SSH key: %v\n", err)
		} else {
			fmt.Printf("SSH key generated successfully at %s.", sshKey)
		}
	} else {
		// If the key file exists, print a message
		fmt.Printf("Key exists at %s\n", sshKey)
	}

	return nil
}
