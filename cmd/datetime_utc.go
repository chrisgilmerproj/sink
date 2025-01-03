package cmd

import (
	"fmt"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	datetimeFlagUnix      = "unix"
	datetimeFlagUnixShort = "u"
)

func initDatetimeUTCFlags(flag *pflag.FlagSet) {
	flag.BoolP(datetimeFlagUnix, datetimeFlagUnixShort, false, "Print the current time in Unix format")
}

func datetimeUTCCmd(cmd *cobra.Command, args []string) error {

	v, errViper := initViper(cmd)
	if errViper != nil {
		return fmt.Errorf("error initializing viper: %w", errViper)
	}
	unix := v.GetBool(datetimeFlagUnix)

	if unix {
		fmt.Println(time.Now().Unix())
		return nil
	}

	// Define the time zones
	timeZones := map[string]string{
		"PST": "America/Los_Angeles",
		"EST": "America/New_York",
		"UTC": "UTC",
	}

	var t time.Time
	if len(args) == 1 {
		timestamp, err := strconv.ParseInt(args[0], 10, 64)
		if err != nil {
			return err
		}
		t = time.Unix(timestamp, 0)

	}

	// Set the format
	datetimeFormat := "2006-01-02T15:04:05-07:00"

	// Iterate through each time zone and print the formatted time
	for zone, location := range timeZones {
		loc, err := time.LoadLocation(location)
		if err != nil {
			fmt.Printf("Error loading location %s: %v\n", location, err)
			continue
		}

		var currentTime string
		if len(args) == 1 {
			currentTime = t.In(loc).Format(datetimeFormat)
		} else {
			// Format the time to RFC 3339 with T separator
			currentTime = time.Now().In(loc).Format(datetimeFormat)
		}
		fmt.Printf("%s %s\n", zone, currentTime)
	}

	return nil
}
