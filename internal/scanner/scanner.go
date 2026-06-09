package scanner

import (
	"fmt"
	"runtime"

	"github.com/mhshahzad/portman/internal/ports"
)

// PortScanner defines the interface for different platform-specific scanners.
type PortScanner interface {
	Scan() ([]ports.Port, error)
	Name() string
}

// GetScanner returns the appropriate scanner for the current platform with fallback logic.
func GetScanner() PortScanner {
	switch runtime.GOOS {
	case "linux":
		return &LinuxStrategy{}
	case "darwin":
		return &DarwinStrategy{}
	default:
		return &UnsupportedScanner{}
	}
}

// LinuxStrategy handles the scanning priority for Linux: ss -> /proc
type LinuxStrategy struct{}

func (s *LinuxStrategy) Name() string { return "linux-strategy" }
func (s *LinuxStrategy) Scan() ([]ports.Port, error) {
	// Try ss first
	ss := &SSScanner{}
	ports, err := ss.Scan()
	if err == nil {
		return ports, nil
	}

	// Fallback to /proc/net/tcp
	proc := &ProcScanner{}
	return proc.Scan()
}

// DarwinStrategy handles the scanning priority for macOS: lsof -> netstat
type DarwinStrategy struct{}

func (s *DarwinStrategy) Name() string { return "darwin-strategy" }
func (s *DarwinStrategy) Scan() ([]ports.Port, error) {
	// Try lsof first
	lsof := &LsofScanner{}
	ports, err := lsof.Scan()
	if err == nil {
		return ports, nil
	}

	// Fallback to netstat
	ns := &NetstatScanner{}
	return ns.Scan()
}

// UnsupportedScanner is a fallback for unsupported operating systems
type UnsupportedScanner struct{}

func (s *UnsupportedScanner) Name() string { return "unsupported" }
func (s *UnsupportedScanner) Scan() ([]ports.Port, error) {
	return nil, fmt.Errorf("operating system %s is not supported", runtime.GOOS)
}

// Global scan function for backward compatibility and easy use
func ScanActivePorts() ([]ports.Port, error) {
	scanner := GetScanner()
	return scanner.Scan()
}
