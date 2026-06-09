package scanner

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/mhshahzad/portman/internal/ports"
)

// ProcScanner implements PortScanner by reading /proc/net/tcp and /proc/net/udp.
type ProcScanner struct{}

func (s *ProcScanner) Name() string { return "proc" }

func (s *ProcScanner) Scan() ([]ports.Port, error) {
	var allPorts []ports.Port

	tcpPorts, err := s.scanFile("/proc/net/tcp", "tcp")
	if err == nil {
		allPorts = append(allPorts, tcpPorts...)
	}

	udpPorts, err := s.scanFile("/proc/net/udp", "udp")
	if err == nil {
		allPorts = append(allPorts, udpPorts...)
	}

	// Also check IPv6
	tcp6Ports, err := s.scanFile("/proc/net/tcp6", "tcp6")
	if err == nil {
		allPorts = append(allPorts, tcp6Ports...)
	}

	if len(allPorts) == 0 && err != nil {
		return nil, fmt.Errorf("could not read any /proc/net files: %w", err)
	}

	return allPorts, nil
}

func (s *ProcScanner) scanFile(path, protocol string) ([]ports.Port, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var result []ports.Port
	scanner := bufio.NewScanner(file)
	// Skip header
	if scanner.Scan() {
		for scanner.Scan() {
			line := scanner.Text()
			fields := strings.Fields(line)
			if len(fields) < 4 {
				continue
			}

			// Local address field is fields[1], format is HEX_IP:HEX_PORT
			localAddr := fields[1]
			parts := strings.Split(localAddr, ":")
			if len(parts) != 2 {
				continue
			}

			portHex := parts[1]
			portDec, err := strconv.ParseUint(portHex, 16, 16)
			if err != nil {
				continue
			}

			// Remote address field is fields[2]. 
			// State field is fields[3]. "0A" is LISTEN for TCP.
			state := fields[3]
			if protocol != "udp" && state != "0A" {
				continue
			}

			result = append(result, ports.Port{
				Number:   int(portDec),
				Protocol: protocol,
				Source:   "proc",
			})
		}
	}

	return result, nil
}

// helper to parse hex IP (not used currently but good for future)
func parseHexIP(h string) (string, error) {
	b, err := hex.DecodeString(h)
	if err != nil {
		return "", err
	}
	if len(b) == 4 {
		return fmt.Sprintf("%d.%d.%d.%d", b[3], b[2], b[1], b[0]), nil
	}
	return "", fmt.Errorf("unsupported IP length")
}
