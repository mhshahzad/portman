package cmd

import (
	"fmt"

	"github.com/mhshahzad/portman/internal/ports"
	"github.com/mhshahzad/portman/internal/scanner"
	"github.com/spf13/cobra"
)

var suggestCmd = &cobra.Command{
	Use:     "suggest",
	Aliases: []string{"next"},
	Short:   "Suggest a safe deployment port (range 3000-9999)",
	RunE: func(cmd *cobra.Command, args []string) error {
		activePorts, err := scanner.ScanActivePorts()
		if err != nil {
			return fmt.Errorf("failed to scan ports: %w", err)
		}

		nextPort := ports.SuggestNext(activePorts)
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
