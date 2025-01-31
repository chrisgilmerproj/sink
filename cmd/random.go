package cmd

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/chrisgilmerproj/sink/v2/pkg/clip"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	randomFlagLength      = "length"
	randomFlagLengthShort = "l"
)

func initRandomFlags(flag *pflag.FlagSet) {
	flag.IntP(randomFlagLength, randomFlagLengthShort, 32, "Set the length of the random string")
}

func randomCmd(cmd *cobra.Command, args []string) error {

	v, errViper := initViper(cmd)
	if errViper != nil {
		return fmt.Errorf("error initializing viper: %w", errViper)
	}
	length := v.GetInt(randomFlagLength)
	verbose := v.GetInt(flagVerbose)

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	result := make([]byte, length)
	for i := 0; i < length; i++ {
		// Generate a random index in the charset
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return err
		}
		// Add the character at the generated index to the result
		result[i] = charset[index.Int64()]
	}
	fmt.Println(string(result))
	clip.CopyToClipboard(string(result), verbose)

	return nil
}
