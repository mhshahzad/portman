package cmd

import (
	"github.com/mhshahzad/portman/internal/integration"
	"github.com/mhshahzad/portman/internal/output"
	"github.com/mhshahzad/portman/internal/ports"
	"github.com/mhshahzad/portman/internal/scanner"
	"github.com/spf13/cobra"
)

var (
	verbose bool
	statusCmd = &cobra.Command{
		Use:   "status",
		Short: "Display all listening TCP and UDP ports",
		RunE: func(cmd *cobra.Command, args []string) error {
			// 1. Collect all scanner outputs
			var allSourcePorts [][]ports.Port

			// System scanner
			systemPorts, err := scanner.ScanActivePorts()
			if err == nil {
				allSourcePorts = append(allSourcePorts, systemPorts)
			}

			// Docker scanner
			dockerScanner, err := scanner.NewDockerScanner()
			if err == nil {
				dockerPorts, err := dockerScanner.Scan()
				if err == nil {
					allSourcePorts = append(allSourcePorts, dockerPorts)
				}
			}

			// 2. Aggregate
			agg := &integration.Aggregator{Verbose: verbose}
			entries := agg.Aggregate(allSourcePorts...)

			// 3. Output
			output.PrintTable(entries, verbose)
			return nil
		},
	}
)

func init() {
	statusCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show detailed source information")
	rootCmd.AddCommand(statusCmd)
}
