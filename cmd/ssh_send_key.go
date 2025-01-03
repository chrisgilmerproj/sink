package cmd

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

const (
	sshFlagKeyPath      = "key-path"
	sshFlagKeyPathShort = "k"
)

func initSSHSendKeyFlags(flag *pflag.FlagSet) {
	flag.StringP(sshFlagKeyPath, sshFlagKeyPathShort, "", "The path of the ssh key")
}

func validateSSHSendKeyFlags(v *viper.Viper) error {
	keyPath := v.GetString(sshFlagName)

	if keyPath != "" {
		if _, err := os.Stat(keyPath); os.IsNotExist(err) {
			return fmt.Errorf("Key path does not exist: %s", keyPath)
		}
	}
	return nil
}

func sshSendKeyCmd(cmd *cobra.Command, args []string) error {

	v, errViper := initViper(cmd)
	if errViper != nil {
		return fmt.Errorf("error initializing viper: %w", errViper)
	}

	err := validateSSHSendKeyFlags(v)
	if err != nil {
		return fmt.Errorf("error validating flags: %w", err)
	}

	remoteHost := args[0]
	publicKeyPath := v.GetString(sshFlagKeyPath)

	if publicKeyPath == "" {
		publicKeyPath = fmt.Sprintf("%s/.ssh/id_rsa.pub", os.Getenv("HOME"))
	}

	// Get the current user's name and the public key file path
	user := os.Getenv("USER")
	if user == "" {
		return fmt.Errorf("Error: USER environment variable not set")
	}

	// Read the public key file
	publicKey, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return fmt.Errorf("Failed to read public key file: %v\n", err)
	}

	if _, errIsPublicKey := isPublicKey(publicKey); errIsPublicKey != nil {
		return fmt.Errorf("Error checking if public key: %v", errIsPublicKey)
	}

	// Set up SSH client configuration
	sshAgentClient, err := getSSHAgentClient()
	if err != nil {
		return fmt.Errorf("Failed to create SSH agent authentication: %v\n", err)
	}

	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeysCallback(sshAgentClient.Signers),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // Skips host key verification
	}

	// Connect to the remote server
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", remoteHost), sshConfig)
	if err != nil {
		return fmt.Errorf("Failed to dial SSH connection: %v\n", err)
	}
	defer client.Close()

	// Create ~/.ssh and authorized_keys file on the remote host
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("Failed to create SSH session: %v\n", err)
	}
	defer session.Close()

	// Run the command to transfer the authorized key
	remoteCmd := fmt.Sprintf(`
		mkdir -p ~/.ssh
		&& chmod 700 ~/.ssh
		&& touch ~/.ssh/authorized_keys
		&& chmod 600 ~/.ssh/authorized_keys
		&& echo "%s" >> ~/.ssh/authorized_keys
	`, strings.TrimSpace(string(publicKey)))

	if out, err := session.CombinedOutput(remoteCmd); err != nil {
		return fmt.Errorf("Failed to send the key: %v\n%s", err, out)
	}

	fmt.Println("Public key successfully added to the remote server.")

	return nil
}

// getSSHAgentClient returns an SSH authentication method using the SSH agent
func getSSHAgentClient() (agent.ExtendedAgent, error) {
	// Connect to the SSH agent
	// Get the SSH_AUTH_SOCK environment variable which contains the path to the socket
	socketPath := os.Getenv("SSH_AUTH_SOCK")
	if socketPath == "" {
		return nil, fmt.Errorf("SSH_AUTH_SOCK environment variable is not set")
	}

	// Dial the Unix socket
	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SSH agent socket: %v", err)
	}

	// Create a new SSH agent client using the connection
	client := agent.NewClient(conn)
	return client, nil
}

func isPublicKey(publicKey []byte) (bool, error) {
	// Check for common private key headers
	privateKeyHeaders := []string{
		"-----BEGIN OPENSSH PRIVATE KEY-----",
		"-----BEGIN RSA PRIVATE KEY-----",
		"-----BEGIN DSA PRIVATE KEY-----",
		"-----BEGIN EC PRIVATE KEY-----",
	}
	for _, header := range privateKeyHeaders {
		if strings.Contains(string(publicKey), header) {
			return false, errors.New("file contains a private key")
		}
	}

	// Try to parse the file as a public key
	if _, _, _, _, err := ssh.ParseAuthorizedKey(publicKey); err != nil {
		return false, fmt.Errorf("file is not a valid public key: %v", err)
	}

	return true, nil
}
