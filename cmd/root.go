package cmd

import (
	"context"
	"fmt"
	"strings"

	gs "github.com/chrisgilmerproj/goshell"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	CliName      = "sink"
	ArtifactName = "sink"

	// Global Flags
	flagVerbose      = "verbose"
	flagVerboseShort = "v"
)

var CC = &gs.CommandChain{} // CommandChain shortcut

type contextVersionKey string

const versionKey contextVersionKey = "version"

func initViper(cmd *cobra.Command) (*viper.Viper, error) {
	v := viper.New()
	errBind := v.BindPFlags(cmd.Flags())
	if errBind != nil {
		return v, fmt.Errorf("error binding flag set to viper: %w", errBind)
	}
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.SetEnvPrefix(strings.ToUpper(CliName))
	v.AutomaticEnv()
	return v, nil
}

// Initialize the flags for the root command
func initRootFlags(flag *pflag.FlagSet) {
	flag.CountP(flagVerbose, flagVerboseShort, "Use verbose output (multiple instances ok)")
}

func checkRootFlags(cmd *cobra.Command, args []string) error {
	_, errViper := initViper(cmd)
	if errViper != nil {
		return fmt.Errorf("error initializing viper: %w", errViper)
	}
	return nil
}

func CreateCommands(version string) *cobra.Command {
	rootCommand := &cobra.Command{
		Use:                   fmt.Sprintf("%s [flags]", CliName),
		DisableFlagsInUseLine: true,
		Short:                 fmt.Sprintf("%s is a CLI tool for commonly used tools", CliName),
		PersistentPreRunE:     checkRootFlags,
	}
	initRootFlags(rootCommand.PersistentFlags())

	datetimeCommand := &cobra.Command{
		Use:                   `datetime [flags]`,
		Args:                  cobra.ExactArgs(0),
		DisableFlagsInUseLine: true,
		Short:                 "datetime commands",
		SilenceErrors:         true,
		SilenceUsage:          true,
	}

	datetimeUTCCommand := &cobra.Command{
		Use:                   `utc [flags] [timestamp]`,
		Args:                  cobra.MaximumNArgs(1),
		DisableFlagsInUseLine: true,
		Short:                 "utc commands",
		SilenceErrors:         true,
		SilenceUsage:          true,
		RunE:                  datetimeUTCCmd,
	}
	initDatetimeUTCFlags(datetimeUTCCommand.Flags())

	datetimeCommand.AddCommand(
		datetimeUTCCommand,
	)

	dockerCommand := &cobra.Command{
		Use:                   `docker [flags]`,
		Args:                  cobra.ExactArgs(0),
		DisableFlagsInUseLine: true,
		Short:                 "docker commands",
		SilenceErrors:         true,
		SilenceUsage:          true,
	}

	dockerRemoveImagesCommand := &cobra.Command{
		Use:                   `remove-images`,
		Args:                  cobra.ExactArgs(0),
		DisableFlagsInUseLine: true,
		Short:                 "remove exited docker images",
		SilenceErrors:         true,
		SilenceUsage:          true,
		RunE:                  dockerRemoveImagesCmd,
	}

	dockerCommand.AddCommand(
		dockerRemoveImagesCommand,
	)

	gpgCommand := &cobra.Command{
		Use:                   `gpg [flags]`,
		Args:                  cobra.ExactArgs(0),
		DisableFlagsInUseLine: true,
		Short:                 "gpg commands",
		SilenceErrors:         true,
		SilenceUsage:          true,
	}

	gpgListKeysCommand := &cobra.Command{
		Use:                   `list-keys`,
		Args:                  cobra.ExactArgs(0),
		DisableFlagsInUseLine: true,
		Short:                 "list gpg keys",
		SilenceErrors:         true,
		SilenceUsage:          true,
		RunE:                  gpgListKeysCmd,
	}

	gpgCommand.AddCommand(
		gpgListKeysCommand,
	)

	networkCommand := &cobra.Command{
		Use:                   `network [flags]`,
		Args:                  cobra.ExactArgs(0),
		DisableFlagsInUseLine: true,
		Short:                 "network commands",
		SilenceErrors:         true,
		SilenceUsage:          true,
	}

	networkClearDNSCommand := &cobra.Command{
		Use:                   `cleardns`,
		Args:                  cobra.ExactArgs(0),
		DisableFlagsInUseLine: true,
		Short:                 "cleardns the local DNS",
		SilenceErrors:         true,
		SilenceUsage:          true,
		RunE:                  networkClearDNSCmd,
	}

	networkIPCommand := &cobra.Command{
		Use:                   `ip [flags]`,
		Args:                  cobra.ExactArgs(0),
		DisableFlagsInUseLine: true,
		Short:                 "get IP values",
		SilenceErrors:         true,
		SilenceUsage:          true,
		RunE:                  networkIPCmd,
	}
	initNetworkIPFlags(networkIPCommand.Flags())

	networkCommand.AddCommand(
		networkClearDNSCommand,
		networkIPCommand,
	)

	randomCommand := &cobra.Command{
		Use:                   `random [flags]`,
		Args:                  cobra.ExactArgs(0),
		DisableFlagsInUseLine: true,
		Short:                 "random string generator",
		SilenceErrors:         true,
		SilenceUsage:          true,
		RunE:                  randomCmd,
	}
	initRandomFlags(randomCommand.Flags())

	sshCommand := &cobra.Command{
		Use:                   `ssh [flags]`,
		Args:                  cobra.ExactArgs(0),
		DisableFlagsInUseLine: true,
		Short:                 "ssh commands",
		SilenceErrors:         true,
		SilenceUsage:          true,
	}

	sshNewKeyCommand := &cobra.Command{
		Use:                   `new-key`,
		Args:                  cobra.ExactArgs(0),
		DisableFlagsInUseLine: true,
		Short:                 "create a new ssh key",
		SilenceErrors:         true,
		SilenceUsage:          true,
		RunE:                  sshNewKeyCmd,
	}
	initSSHNewKeyFlags(sshNewKeyCommand.Flags())

	sshSendKeyCommand := &cobra.Command{
		Use:                   `send-key`,
		Args:                  cobra.ExactArgs(1),
		DisableFlagsInUseLine: true,
		Short:                 "send an ssh key to a host",
		SilenceErrors:         true,
		SilenceUsage:          true,
		RunE:                  sshSendKeyCmd,
	}
	initSSHSendKeyFlags(sshSendKeyCommand.Flags())

	sshCommand.AddCommand(
		sshNewKeyCommand,
		sshSendKeyCommand,
	)

	uuidCommand := &cobra.Command{
		Use:                   `uuid`,
		Args:                  cobra.ExactArgs(0),
		DisableFlagsInUseLine: true,
		Short:                 "uuid4 generator",
		SilenceErrors:         true,
		SilenceUsage:          true,
		RunE:                  uuidCmd,
	}

	serverCommand := &cobra.Command{
		Use:                   `server [flags]`,
		Args:                  cobra.ExactArgs(0),
		DisableFlagsInUseLine: true,
		Short:                 "run simple server on local directory",
		SilenceErrors:         true,
		SilenceUsage:          true,
		RunE:                  serverCmd,
	}
	initServerFlags(serverCommand.Flags())

	versionCommand := &cobra.Command{
		Use:                   `version`,
		Args:                  cobra.ExactArgs(0),
		DisableFlagsInUseLine: true,
		Short:                 "display the version and quit",
		SilenceErrors:         true,
		SilenceUsage:          true,
		PreRun: func(cmd *cobra.Command, args []string) {
			cmd.SetContext(context.WithValue(cmd.Context(), versionKey, version))
		},
		RunE: versionCmd,
	}
	initVersionFlags(versionCommand.Flags())

	rootCommand.AddCommand(
		datetimeCommand,
		dockerCommand,
		gpgCommand,
		networkCommand,
		randomCommand,
		serverCommand,
		sshCommand,
		uuidCommand,
		versionCommand,
	)
	return rootCommand
}
