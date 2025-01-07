package cmd

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/spf13/cobra"

	"github.com/chrisgilmerproj/sink/v2/pkg/clip"
)

func uuidCmd(cmd *cobra.Command, args []string) error {

	v, errViper := initViper(cmd)
	if errViper != nil {
		return fmt.Errorf("error initializing viper: %w", errViper)
	}
	verbose := v.GetInt(flagVerbose)

	// Generate a new UUID
	newUUID := uuid.New()

	// Convert UUID to lowercase
	lowercaseUUID := strings.ToLower(newUUID.String())

	// Print the UUID to the terminal
	fmt.Println(lowercaseUUID)
	clip.CopyToClipboard(lowercaseUUID, verbose)

	return nil
}
