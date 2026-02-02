package cmd

import (
	"errors"

	"github.com/spf13/cobra"
)

var namespaceOverride string

func registerNamespaceFlag(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&namespaceOverride, "namespace", "n", "", "Workspace namespace to target")
}

func currentNamespace() (string, error) {
	if namespaceOverride != "" {
		return namespaceOverride, nil
	}
	if currentConfig == nil || currentConfig.CurrentWorkspace == "" {
		return "", errors.New("no active workspace: run 'shoulders workspace use <name>' or pass --namespace")
	}
	return currentConfig.CurrentWorkspace, nil
}
