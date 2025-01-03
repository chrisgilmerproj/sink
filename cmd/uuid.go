package cmd

import (
	"fmt"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

func uuidCmd(cmd *cobra.Command, args []string) error {

	_, errViper := initViper(cmd)
	if errViper != nil {
		return fmt.Errorf("error initializing viper: %w", errViper)
	}

	// Generate a new UUID
	newUUID := uuid.New()

	// Convert UUID to lowercase
	lowercaseUUID := strings.ToLower(newUUID.String())

	// Copy UUID to clipboard
	err := clipboard.WriteAll(lowercaseUUID)
	if err != nil {
		return fmt.Errorf("failed to copy UUID to clipboard: %v", err)
	}

	// Print the UUID to the terminal
	fmt.Println(lowercaseUUID)

	return nil
}
