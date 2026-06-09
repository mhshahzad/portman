package output

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/mhshahzad/portman/internal/integration"
)

// PrintTable displays the list of ports in a formatted table.
func PrintTable(entries []integration.PortEntry, verbose bool) {
	if len(entries) == 0 {
		fmt.Println("No active ports found.")
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
	
	header := "PORT\tPROTOCOL\tSERVICE\tPROCESS\tPID\tBACKING"
	if verbose {
		header += "\tSOURCES"
	}
	fmt.Fprintln(w, header)

	for _, p := range entries {
		process := p.Process
		if process == "" || process == "docker-proxy" || process == "com.docker.backend" {
			process = "-"
		}

		pidStr := fmt.Sprintf("%d", p.PID)
		if p.PID == 0 {
			pidStr = "-"
		}

		line := fmt.Sprintf("%d\t%s\t%s\t%s\t%s\t%s", 
			p.Port, p.Protocol, p.Service, process, pidStr, p.Backing)
		
		if verbose {
			line += fmt.Sprintf("\t%s", strings.Join(p.Sources, ","))
		}
		
		fmt.Fprintln(w, line)
	}
	w.Flush()
}
