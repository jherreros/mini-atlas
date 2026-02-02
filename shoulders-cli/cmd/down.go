package cmd

import (
	"github.com/jherreros/shoulders/shoulders-cli/internal/bootstrap"
	"github.com/spf13/cobra"
)

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Delete the local kind cluster",
	RunE: func(cmd *cobra.Command, args []string) error {
		return bootstrap.DeleteKindCluster(bootstrap.DefaultClusterName)
	},
}
