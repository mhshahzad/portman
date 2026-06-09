package cmd

import (
	"fmt"

	"github.com/mhshahzad/portman/internal/integration"
	"github.com/mhshahzad/portman/internal/ports"
	"github.com/mhshahzad/portman/internal/scanner"
	"github.com/spf13/cobra"
)

var suggestCmd = &cobra.Command{
	Use:     "suggest",
	Aliases: []string{"next"},
	Short:   "Suggest a safe deployment port (range 3000-9999)",
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
		agg := &integration.Aggregator{}
		entries := agg.Aggregate(allSourcePorts...)

		// 3. Suggest
		occupied := make(map[int]bool)
		for _, e := range entries {
			occupied[e.Port] = true
		}

		nextPort := ports.SuggestNext(occupied)
		if nextPort == -1 {
			return fmt.Errorf("no available ports in range 3000-9999")
		}

		fmt.Println(nextPort)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(suggestCmd)
}
