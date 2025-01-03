package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"
)

func networkClearDNSCmd(cmd *cobra.Command, args []string) error {

	_, errViper := initViper(cmd)
	if errViper != nil {
		return fmt.Errorf("error initializing viper: %w", errViper)
	}

	var output string
	var err error
	output, err = CC.Run([][]string{
		{"sudo", "killall", "-HUP", "mDNSResponder"},
	})
	if err != nil {
		log.Fatalf("Error running command chain: %v", err)
	}
	fmt.Print(strings.TrimSpace(output))

	return nil
}
