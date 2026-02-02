package cmd

import (
	"context"
	"fmt"

	"github.com/jherreros/shoulders/shoulders-cli/internal/config"
	"github.com/jherreros/shoulders/shoulders-cli/internal/kube"
	"github.com/jherreros/shoulders/shoulders-cli/internal/output"
	"github.com/jherreros/shoulders/shoulders-cli/pkg/api/v1alpha1"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/yaml"
)

var workspaceCmd = &cobra.Command{
	Use:   "workspace",
	Short: "Manage workspace contexts",
}

var workspaceCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a Workspace",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		workspace := v1alpha1.Workspace{
			TypeMeta:   v1alpha1.TypeMeta("Workspace"),
			ObjectMeta: v1alpha1.ObjectMeta(name, ""),
		}

		content, err := yaml.Marshal(workspace)
		if err != nil {
			return err
		}

		dynamicClient, err := kube.NewDynamicClient(kubeconfig)
		if err != nil {
			return err
		}
		obj := &unstructured.Unstructured{}
		if err := yaml.Unmarshal(content, obj); err != nil {
			return err
		}

		gvr := schema.GroupVersionResource{Group: v1alpha1.Group, Version: v1alpha1.Version, Resource: "workspaces"}
		if err := kube.Apply(context.Background(), dynamicClient, gvr, "", obj); err != nil {
			return err
		}

		fmt.Printf("Workspace %s created\n", name)
		return nil
	},
}

var workspaceUseCmd = &cobra.Command{
	Use:   "use <name>",
	Short: "Set the active workspace",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		dynamicClient, err := kube.NewDynamicClient(kubeconfig)
		if err != nil {
			return err
		}
		gvr := schema.GroupVersionResource{Group: v1alpha1.Group, Version: v1alpha1.Version, Resource: "workspaces"}
		if _, err := dynamicClient.Resource(gvr).Get(context.Background(), name, metav1.GetOptions{}); err != nil {
			return err
		}

		currentConfig.CurrentWorkspace = name
		if err := config.Save(currentConfig); err != nil {
			return err
		}
		fmt.Printf("Active workspace set to %s\n", name)
		return nil
	},
}

var workspaceListCmd = &cobra.Command{
	Use:   "list",
	Short: "List Workspaces",
	RunE: func(cmd *cobra.Command, args []string) error {
		format, err := outputOption()
		if err != nil {
			return err
		}
		dynamicClient, err := kube.NewDynamicClient(kubeconfig)
		if err != nil {
			return err
		}
		gvr := schema.GroupVersionResource{Group: v1alpha1.Group, Version: v1alpha1.Version, Resource: "workspaces"}
		list, err := dynamicClient.Resource(gvr).List(context.Background(), metav1.ListOptions{})
		if err != nil {
			return err
		}
		if format == output.Table {
			rows := [][]string{}
			for _, item := range list.Items {
				rows = append(rows, []string{item.GetName()})
			}
			return output.PrintTable([]string{"Name"}, rows)
		}

		payload, err := output.Render(list.Items, format)
		if err != nil {
			return err
		}
		fmt.Println(string(payload))
		return nil
	},
}

func init() {
	workspaceCmd.AddCommand(workspaceCreateCmd)
	workspaceCmd.AddCommand(workspaceUseCmd)
	workspaceCmd.AddCommand(workspaceListCmd)
}
