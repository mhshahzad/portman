package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	Version = "dev"
	rootCmd = &cobra.Command{
		Use:   "portman",
		Short: "Portman is a CLI tool for identifying and suggesting ports on Linux hosts.",
		Long: `Portman helps DevOps and platform engineers identify allocated ports 
and safely recommend available ports for new deployments.`,
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
