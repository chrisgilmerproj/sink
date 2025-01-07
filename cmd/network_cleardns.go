package cmd

import (
	"fmt"
	"log"
	"runtime"
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

	var command []string
	switch runtime.GOOS {
	case "darwin":
		command = []string{"sudo", "killall", "-HUP", "mDNSResponder"}
	case "linux":
		command = []string{"sudo", "/usr/local/sbin/dnsmasq-reload"}
	}
	output, err = CC.Run([][]string{command})
	if err != nil {
		log.Fatalf("Error running command chain: %v", err)
	}
	fmt.Print(strings.TrimSpace(output))

	return nil
}
