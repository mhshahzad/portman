package scanner

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/mhshahzad/portman/internal/ports"
)

// LsofScanner implements PortScanner using the 'lsof' command.
type LsofScanner struct{}

func (s *LsofScanner) Name() string { return "lsof" }

var lsofRegex = regexp.MustCompile(`^(\S+)\s+(\d+)\s+\S+\s+\S+\s+\S+\s+\S+\s+\S+\s+\S+\s+(\S+)\s+\(LISTEN\)$`)

func (s *LsofScanner) Scan() ([]ports.Port, error) {
	// -i: internet files
	// -P: inhibit conversion of port numbers to port names
	// -n: inhibit conversion of network numbers to host names
	// +c 0: show full command name
	cmd := exec.Command("lsof", "-i", "-P", "-n", "+c", "0")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("lsof command failed: %w", err)
	}

	return s.parse(string(output)), nil
}

func (s *LsofScanner) parse(output string) []ports.Port {
	lines := strings.Split(output, "\n")
	var result []ports.Port
	seen := make(map[string]bool)

	for _, line := range lines {
		// COMMAND PID USER FD TYPE DEVICE SIZE/OFF NODE NAME
		// sshd 123 root 3u IPv4 0x... 0t0 TCP *:22 (LISTEN)
		if !strings.Contains(line, "(LISTEN)") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 9 {
			continue
		}

		command := strings.ReplaceAll(fields[0], "\\x20", " ")
		pid, _ := strconv.Atoi(fields[1])
		protocol := strings.ToLower(fields[7])
		name := fields[8] // e.g., *:22 or 127.0.0.1:8080

		portStr := name
		if idx := strings.LastIndex(name, ":"); idx != -1 {
			portStr = name[idx+1:]
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
			Number:      portNum,
			Protocol:    protocol,
			ProcessName: command,
			PID:         pid,
			Source:      "lsof",
		})
	}

	return result
}
