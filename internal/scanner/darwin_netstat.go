package scanner

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/mhshahzad/portman/internal/ports"
)

// NetstatScanner implements PortScanner using the 'netstat' command.
type NetstatScanner struct{}

func (s *NetstatScanner) Name() string { return "netstat" }

func (s *NetstatScanner) Scan() ([]ports.Port, error) {
	// -a: all
	// -n: numeric
	// -v: verbose (sometimes needed for more info)
	cmd := exec.Command("netstat", "-an")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("netstat command failed: %w", err)
	}

	return s.parse(string(output)), nil
}

func (s *NetstatScanner) parse(output string) []ports.Port {
	lines := strings.Split(output, "\n")
	var result []ports.Port
	seen := make(map[string]bool)

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}

		protocol := strings.ToLower(fields[0])
		if !strings.HasPrefix(protocol, "tcp") && !strings.HasPrefix(protocol, "udp") {
			continue
		}

		// Local address is usually field 3
		localAddr := fields[3]
		if !strings.Contains(line, "LISTEN") && !strings.HasPrefix(protocol, "udp") {
			continue
		}

		portStr := localAddr
		if idx := strings.LastIndex(localAddr, "."); idx != -1 {
			portStr = localAddr[idx+1:]
		} else if idx := strings.LastIndex(localAddr, ":"); idx != -1 {
			portStr = localAddr[idx+1:]
		}

		portNum, err := strconv.Atoi(portStr)
		if err != nil {
			continue
		}

		// Deduplicate
		key := fmt.Sprintf("%d-%s", portNum, protocol)
		if seen[key] {
			continue
		}
		seen[key] = true

		result = append(result, ports.Port{
			Number:   portNum,
			Protocol: protocol,
			Source:   "netstat",
		})
	}

	return result
}
