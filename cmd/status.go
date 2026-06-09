package cmd

import (
	"fmt"

	"github.com/mhshahzad/portman/internal/output"
	"github.com/mhshahzad/portman/internal/scanner"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Display all listening TCP and UDP ports",
	RunE: func(cmd *cobra.Command, args []string) error {
		activePorts, err := scanner.ScanActivePorts()
		if err != nil {
			return fmt.Errorf("failed to scan ports: %w", err)
		}

		output.PrintTable(activePorts)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
