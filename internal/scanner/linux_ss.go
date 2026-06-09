package scanner

import (
	"fmt"
	"os/exec"

	"github.com/mhshahzad/portman/internal/ports"
)

// SSScanner implements PortScanner using the 'ss' command.
type SSScanner struct{}

func (s *SSScanner) Name() string { return "ss" }

func (s *SSScanner) Scan() ([]ports.Port, error) {
	cmd := exec.Command("ss", "-tulpn")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("ss command failed: %w", err)
	}

	return s.parse(string(output)), nil
}

func (s *SSScanner) parse(output string) []ports.Port {
	// Re-using the logic from ParseSSOutput but ensuring Source is set
	parsed := ParseSSOutput(output)
	for i := range parsed {
		parsed[i].Source = "ss"
	}
	return parsed
}
