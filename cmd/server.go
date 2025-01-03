package cmd

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	serverFlagDirectory      = "directory"
	serverFlagDirectoryShort = "d"
	serverFlagPort           = "port"
	serverFlagPortShort      = "p"
)

func initServerFlags(flag *pflag.FlagSet) {
	flag.StringP(serverFlagDirectory, serverFlagDirectoryShort, ".", "Set the directory to serve")
	flag.StringP(serverFlagPort, serverFlagPortShort, "8000", "Set the port of the server")
}

func serverCmd(cmd *cobra.Command, args []string) error {

	v, errViper := initViper(cmd)
	if errViper != nil {
		return fmt.Errorf("error initializing viper: %w", errViper)
	}
	directory := v.GetString(serverFlagDirectory)
	port := v.GetString(serverFlagPort)

	// Set up a channel to listen for OS signals
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGINT, syscall.SIGTERM) // Listen for SIGINT (Ctrl+C) and SIGTERM

	// Serve files from the current directory
	http.Handle("/", http.FileServer(http.Dir(directory)))

	go func() {
		// Start the server on port 8000
		fmt.Printf("Starting server on port %s...", port)
		err := http.ListenAndServe(":"+port, nil)
		if err != nil {
			fmt.Println("Error starting server:", err)
		}
	}()

	// Wait for an interrupt signal (Ctrl+C)
	<-signalChannel
	fmt.Println("\nShutting down server...")

	return nil
}
