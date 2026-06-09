package scanner

import (
	"fmt"
	"os/exec"

	"github.com/mhshahzad/portman/internal/ports"
)

// ScanActivePorts executes 'ss -tulpn' and returns a list of active ports.
func ScanActivePorts() ([]ports.Port, error) {
	// We use 'ss -tulpn' as specified in the requirements.
	// -t: tcp
	// -u: udp
	// -l: listening
	// -p: processes
	// -n: numeric
	cmd := exec.Command("ss", "-tulpn")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to execute 'ss': %w (ensure you are on Linux and 'iproute2' is installed)", err)
	}

	return ParseSSOutput(string(output)), nil
}
