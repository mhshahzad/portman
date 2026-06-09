package scanner

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/mhshahzad/portman/internal/ports"
)

var (
	// regex to match port from "address:port" or "[address]:port"
	portRegex = regexp.MustCompile(`:(\d+)$`)
	// regex to match process name and PID from users:(("name",pid=123,fd=4))
	processRegex = regexp.MustCompile(`users:\(\("([^"]+)",pid=(\d+),`)
)

// ParseSSOutput parses the output of `ss -tulpn`.
func ParseSSOutput(output string) []ports.Port {
	lines := strings.Split(output, "\n")
	var result []ports.Port

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 5 {
			continue
		}

		// Skip header
		if fields[0] == "Netid" || fields[0] == "State" {
			continue
		}

		protocol := fields[0]
		localAddr := fields[4]

		// Extract port
		portMatch := portRegex.FindStringSubmatch(localAddr)
		if len(portMatch) < 2 {
			continue
		}
		portNum, _ := strconv.Atoi(portMatch[1])

		port := ports.Port{
			Number:   portNum,
			Protocol: protocol,
		}

		// Extract Process and PID if available (usually last field)
		if len(fields) >= 7 {
			usersField := fields[len(fields)-1]
			processMatch := processRegex.FindStringSubmatch(usersField)
			if len(processMatch) >= 3 {
				port.ProcessName = processMatch[1]
				pid, _ := strconv.Atoi(processMatch[2])
				port.PID = pid
			}
		}

		result = append(result, port)
	}

	return result
}
