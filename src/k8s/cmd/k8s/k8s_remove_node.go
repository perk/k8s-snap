package k8s

import (
	"fmt"

	"github.com/canonical/k8s/pkg/k8s/client"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	removeNodeCmdOpts struct {
		force bool
	}

	removeNodeCmd = &cobra.Command{
		Use:    "remove-node <name>",
		Short:  "Remove a node from the cluster",
		Hidden: true,
		Args:   cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if rootCmdOpts.logDebug {
				logrus.SetLevel(logrus.TraceLevel)
			}

			name := args[0]
			client, err := client.NewClient(cmd.Context(), client.ClusterOpts{
				RemoteAddress: clusterCmdOpts.remoteAddress,
				StorageDir:    clusterCmdOpts.storageDir,
				Verbose:       rootCmdOpts.logVerbose,
				Debug:         rootCmdOpts.logDebug,
			})
			if err != nil {
				return fmt.Errorf("failed to create cluster client: %w", err)
			}

			err = client.RemoveNode(cmd.Context(), name, removeNodeCmdOpts.force)
			if err != nil {
				return fmt.Errorf("failed to remove node from cluster: %w", err)
			}
			logrus.Infof("Removed %s from cluster", name)
			return nil
		},
	}
)

func init() {
	removeNodeCmd.Flags().BoolVar(&removeNodeCmdOpts.force, "force", false, "Forcibly remove the cluster member")

	rootCmd.AddCommand(removeNodeCmd)
}
