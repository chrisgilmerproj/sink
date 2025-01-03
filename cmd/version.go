package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	verboseFlagShort      = "short"
	verboseFlagShortShort = "s"
)

func initVersionFlags(flag *pflag.FlagSet) {
	flag.BoolP(verboseFlagShort, verboseFlagShortShort, false, "Print the version number only")
}

func versionCmd(cmd *cobra.Command, args []string) error {

	v, errViper := initViper(cmd)
	if errViper != nil {
		return fmt.Errorf("error initializing viper: %w", errViper)
	}
	short := v.GetBool(verboseFlagShort)
	version := cmd.Context().Value(versionKey).(string)

	if short {
		fmt.Println(version)
	} else {
		fmt.Printf("%s version %s\n", CliName, version)
	}
	return nil
}
