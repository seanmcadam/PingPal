package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

// CheckLatency measures the round-trip time to the specified IP address using ICMP echo
// It returns the latency in milliseconds and any error encountered
func CheckLatencyICMP(ipAddr string, timeout time.Duration) (float64, error) {
	// Create a new ICMP connection
	// On most systems, you need to run as root/administrator to use raw sockets
	conn, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		return 0, fmt.Errorf("error creating ICMP listener: %w", err)
	}
	defer conn.Close()

	// Resolve the IP address
	dst, err := net.ResolveIPAddr("ip4", ipAddr)
	if err != nil {
		return 0, fmt.Errorf("error resolving IP address: %w", err)
	}

	// Create an ICMP message
	msg := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff, // Use process ID as identifier
			Seq:  1,                    // Sequence number
			Data: []byte("ping test"),  // Payload data
		},
	}

	// Marshal the message into binary
	msgBytes, err := msg.Marshal(nil)
	if err != nil {
		return 0, fmt.Errorf("error marshaling ICMP message: %w", err)
	}

	// Set a deadline on the connection
	if err := conn.SetDeadline(time.Now().Add(timeout)); err != nil {
		return 0, fmt.Errorf("error setting deadline: %w", err)
	}

	// Record the start time
	startTime := time.Now()

	// Send the ICMP packet
	_, err = conn.WriteTo(msgBytes, dst)
	if err != nil {
		return 0, fmt.Errorf("error sending ICMP packet: %w", err)
	}

	// Prepare to receive the reply
	reply := make([]byte, 1500) // Buffer size for reply
	n, _, err := conn.ReadFrom(reply)
	if err != nil {
		return 0, fmt.Errorf("error receiving ICMP reply: %w", err)
	}

	// Calculate elapsed time
	elapsed := time.Since(startTime)

	// Parse the reply message
	parsedMsg, err := icmp.ParseMessage(ipv4.ICMPTypeEcho.Protocol(), reply[:n])
	if err != nil {
		return 0, fmt.Errorf("error parsing ICMP reply: %w", err)
	}

	// Check if it's an echo reply
	if parsedMsg.Type != ipv4.ICMPTypeEchoReply {
		return 0, fmt.Errorf("received non-echo reply message: %v", parsedMsg.Type)
	}

	return float64(elapsed.Microseconds()) / 1000.0, nil // Return latency in milliseconds
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

func main() {
	addrs := ParseFlags()

	timeout := 2 * time.Second

	output := ""

	for _, addr := range addrs {
		latency, err := CheckLatencyICMP(addr, timeout)
		if err != nil {
			fmt.Printf("Error checking latency: %v\n", err)
			os.Exit(1)
		}
		output += fmt.Sprintf("Ping to %s: latency = %.3f ms\n", addr, latency)
	}

	fmt.Printf(output)
}
