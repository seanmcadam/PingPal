package config

import (
	"flag"
	"os"
	"strings"
)

// Alternative implementation with more explicit validation
func ParseFlagsWithValidation() (*Input, error) {
	settings := &SessionSettings{}
	var aFlag StringSliceFlag

	// Create custom FlagSet for more control
	fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	fs.Uint64Var(&settings.DisplayRefreshTimeS, "d", 1, "Display refresh rate in seconds (uint64)")
	fs.Uint64Var(&settings.PktDropTimeS, "p", 300, "Time to average latency and packet loss across in seconds (uint64)")
	fs.Uint64Var(&settings.LatencyCheckIntervalS, "l", 5, "Latency check interval seconds (uint64)")
	fs.Uint64Var(&settings.ConnectionTimeoutS, "c", 120, "Connection timeout in seconds (uint64)")
	fs.Var(&aFlag, "a", "IP address to monitor (string, repeatable)")

	// Parse with custom error handling
	if err := fs.Parse(os.Args[1:]); err != nil {
		return nil, err
	}

	input := Input{Addresses: []string(aFlag), Settings: *settings}

	return &input, nil
}

type Input struct {
	Settings  SessionSettings
	Addresses []string
}

// The settings for this session including update time preferences
type SessionSettings struct {
	DisplayRefreshTimeS   uint64 // how often to refresh the display in seconds
	PktDropTimeS          uint64 // time to retain packets and average latency and loss across in seconds
	LatencyCheckIntervalS uint64 // how often to check latency in seconds
	ConnectionTimeoutS    uint64 // how long to wait for a reply to the ICMP packet in seconds
}

// StringSliceFlag is a custom flag type that can collect multiple values
type StringSliceFlag []string

// String returns the string representation of the flag values
func (s *StringSliceFlag) String() string {
	return strings.Join(*s, ", ")
}

// Set appends each value to the slice
// This is called each time the flag appears on the command line
func (s *StringSliceFlag) Set(value string) error {
	*s = append(*s, value)
	return nil
}

// ParseRepeatedFlag parses command line arguments and collects all values
// provided with the '-a' flag.
// Returns a slice containing all the values.
func ParseFlags() []string {
	// Create our custom flag
	var addresses StringSliceFlag

	// Register the flag with the flag package
	flag.Var(&addresses, "a", "IP address to monitor (can be repeated)")
	flag.Parse()
	return addresses
}
