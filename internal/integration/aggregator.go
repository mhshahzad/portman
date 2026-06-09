package integration

import (
	"fmt"
	"sort"

	"github.com/mhshahzad/portman/internal/ports"
)

// PortEntry represents a unified, resolved view of a port.
type PortEntry struct {
	Port     int
	Protocol string

	Service string // unified service name
	Backing string // docker | system | unknown

	Process string
	PID     int

	Source  string   // primary source (ss, lsof, docker, proc)
	Sources []string // all contributing sources
}

// Aggregator collects and merges results from multiple scanners.
type Aggregator struct {
	Verbose bool
}

// Aggregate merges multiple slices of ports into a single unified dataset.
// It follows the priority: Docker > System (lsof/ss) > Proc.
func (a *Aggregator) Aggregate(allSources ...[]ports.Port) []PortEntry {
	groups := make(map[string][]ports.Port)

	for _, sourcePorts := range allSources {
		for _, p := range sourcePorts {
			key := fmt.Sprintf("%d/%s", p.Number, p.Protocol)
			groups[key] = append(groups[key], p)
		}
	}

	var result []PortEntry
	for _, group := range groups {
		result = append(result, a.resolve(group))
	}

	// Deterministic sort: Port number, then Protocol
	sort.Slice(result, func(i, j int) bool {
		if result[i].Port != result[j].Port {
			return result[i].Port < result[j].Port
		}
		return result[i].Protocol < result[j].Protocol
	})

	return result
}

func (a *Aggregator) resolve(entries []ports.Port) PortEntry {
	// Priority map: higher is better
	priorities := map[string]int{
		"docker":  30,
		"lsof":    20,
		"ss":      20,
		"netstat": 15,
		"proc":    10,
	}

	var primary ports.Port
	maxPriority := -1
	allSources := make(map[string]bool)

	for _, e := range entries {
		allSources[e.Source] = true
		p := priorities[e.Source]
		if p > maxPriority {
			maxPriority = p
			primary = e
		} else if p == maxPriority {
			// Tie-breaker: prefer entries with more metadata
			if primary.ProcessName == "" && e.ProcessName != "" {
				primary = e
			}
		}
	}

	resolved := PortEntry{
		Port:     primary.Number,
		Protocol: primary.Protocol,
		Service:  primary.Service,
		Backing:  "system",
		Process:  primary.ProcessName,
		PID:      primary.PID,
		Source:   primary.Source,
	}

	if primary.Source == "docker" {
		resolved.Backing = "docker"
		// If primary is Docker, it might lack PID/ProcessName. 
		// Try to find them from other sources in the group (like lsof/ss).
		if resolved.Process == "" || resolved.Process == "-" {
			for _, e := range entries {
				if e.ProcessName != "" && e.ProcessName != "-" && e.ProcessName != "docker-proxy" && e.ProcessName != "com.docker.backend" {
					resolved.Process = e.ProcessName
					resolved.PID = e.PID
					break
				}
			}
		}
	}

	// Ensure we have a service name
	if resolved.Service == "" {
		resolved.Service = resolved.Process
	}
	if resolved.Service == "" || resolved.Service == "-" {
		// Try to find a process name from any other source in the group if primary lacks it
		for _, e := range entries {
			if e.ProcessName != "" && e.ProcessName != "-" && e.ProcessName != "docker-proxy" && e.ProcessName != "com.docker.backend" {
				resolved.Process = e.ProcessName
				resolved.PID = e.PID
				if resolved.Service == "" || resolved.Service == "-" {
					resolved.Service = e.ProcessName
				}
				break
			}
		}
	}

	// Clean up service name fallback
	if resolved.Service == "" {
		resolved.Service = "-"
	}

	for s := range allSources {
		resolved.Sources = append(resolved.Sources, s)
	}
	sort.Strings(resolved.Sources)

	return resolved
}
