package ports

import "sort"

// Port represents an active listening port on the system.
type Port struct {
	Number      int
	Protocol    string // tcp or udp
	ProcessName string // Name of the process owning the port
	PID         int    // Process ID owning the port
	Source      string // ss, lsof, netstat, proc, docker
	Service     string // human-readable service name
	ContainerID string // ID of the container (if applicable)
}

// SuggestNext returns the first available port in the range 3000-9999.
// If all ports in the range are occupied, it returns -1.
func SuggestNext(occupiedPorts map[int]bool) int {
	const startPort = 3000
	const endPort = 9999

	for p := startPort; p <= endPort; p++ {
		if !occupiedPorts[p] {
			return p
		}
	}

	return -1
}

// SortPorts sorts a slice of Port by port number.
func SortPorts(ports []Port) {
	sort.Slice(ports, func(i, j int) bool {
		return ports[i].Number < ports[j].Number
	})
}
