package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/chrisgilmerproj/sink/v2/pkg/clip"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	networkFlagLocal      = "local"
	networkFlagLocalShort = "l"
)

func initNetworkIPFlags(flag *pflag.FlagSet) {
	flag.BoolP(networkFlagLocal, networkFlagLocalShort, false, "Get the local network IP")
}

func networkIPCmd(cmd *cobra.Command, args []string) error {

	v, errViper := initViper(cmd)
	if errViper != nil {
		return fmt.Errorf("error initializing viper: %w", errViper)
	}
	local := v.GetBool(networkFlagLocal)
	verbose := v.GetInt(flagVerbose)

	var output string
	var err error
	if local {
		output, err = CC.Run([][]string{
			{"ifconfig", "en0"},
			{"grep", "inet "},
			{"cut", "-d", " ", "-f2"},
		})
	} else {
		output, err = CC.Run([][]string{
			{"dig", "+short", "myip.opendns.com", "@resolver1.opendns.com"},
		})
	}
	if err != nil {
		log.Fatalf("Error running command chain: %v", err)
	}
	fmt.Print(strings.TrimSpace(output))
	clip.CopyToClipboard(output, verbose)

	return nil
}
