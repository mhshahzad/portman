package output

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/mhshahzad/portman/internal/ports"
)

// PrintTable displays the list of ports in a formatted table.
func PrintTable(activePorts []ports.Port) {
	if len(activePorts) == 0 {
		fmt.Println("No active ports found.")
		return
	}

	// Sort ports by number for better readability
	ports.SortPorts(activePorts)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
	fmt.Fprintln(w, "PORT\tPROTOCOL\tPROCESS\tPID")

	for _, p := range activePorts {
		process := p.ProcessName
		if process == "" {
			process = "-"
		}
		pidStr := fmt.Sprintf("%d", p.PID)
		if p.PID == 0 {
			pidStr = "-"
		}
		fmt.Fprintf(w, "%d\t%s\t%s\t%s\n", p.Number, p.Protocol, process, pidStr)
	}
	w.Flush()
}
