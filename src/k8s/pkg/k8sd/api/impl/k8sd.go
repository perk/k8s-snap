package impl

import (
	"context"
	"fmt"
	"log"

	apiv1 "github.com/canonical/k8s/api/v1"
	"github.com/canonical/k8s/pkg/snap"
	snaputil "github.com/canonical/k8s/pkg/snap/util"
	"github.com/canonical/k8s/pkg/utils"
	"github.com/canonical/k8s/pkg/utils/k8s"
	"github.com/canonical/microcluster/state"
)

// GetClusterStatus retrieves the status of the cluster, including information about its members.
func GetClusterStatus(ctx context.Context, s *state.State) (apiv1.ClusterStatus, error) {
	snap := snap.SnapFromContext(s.Context)

	client, err := k8s.NewClient(snap)
	if err != nil {
		return apiv1.ClusterStatus{}, fmt.Errorf("failed to create k8s client: %w", err)
	}

	if err := client.WaitApiServerReady(ctx); err != nil {
		return apiv1.ClusterStatus{}, fmt.Errorf("k8s api server did not become ready in time: %w", err)
	}

	ready, err := client.ClusterReady(ctx)
	if err != nil {
		return apiv1.ClusterStatus{}, fmt.Errorf("failed to get cluster components: %w", err)
	}

	members, err := GetClusterMembers(ctx, s)
	if err != nil {
		return apiv1.ClusterStatus{}, fmt.Errorf("failed to get cluster members: %w", err)
	}

	config, err := utils.GetUserFacingClusterConfig(ctx, s)
	if err != nil {
		return apiv1.ClusterStatus{}, fmt.Errorf("failed to get user-facing cluster config: %w", err)
	}

	return apiv1.ClusterStatus{
		Ready:   ready,
		Members: members,
		Config:  config,
	}, nil
}

// GetClusterMembers retrieves information about the members of the cluster.
func GetClusterMembers(ctx context.Context, s *state.State) ([]apiv1.NodeStatus, error) {
	c, err := s.Leader()
	if err != nil {
		return nil, fmt.Errorf("failed to get leader client: %w", err)
	}

	clusterMembers, err := c.GetClusterMembers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get cluster members: %w", err)
	}

	members := make([]apiv1.NodeStatus, len(clusterMembers))
	for i, clusterMember := range clusterMembers {
		members[i] = apiv1.NodeStatus{
			Name:          clusterMember.Name,
			Address:       clusterMember.Address.String(),
			ClusterRole:   apiv1.ClusterRoleControlPlane,
			DatastoreRole: utils.DatastoreRoleFromString(clusterMember.Role),
		}
	}

	return members, nil
}

// GetLocalNodeStatus retrieves the status of the local node, including its roles within the cluster.
// Unlike "GetClusterMembers" this also works on a worker node.
func GetLocalNodeStatus(ctx context.Context, s *state.State) (apiv1.NodeStatus, error) {
	snap := snap.SnapFromContext(s.Context)

	// Determine cluster role.
	clusterRole := apiv1.ClusterRoleUnknown
	isWorker, err := snaputil.IsWorker(snap)
	if err != nil {
		return apiv1.NodeStatus{}, fmt.Errorf("failed to check if node is a worker: %w", err)
	}
	if isWorker {
		clusterRole = apiv1.ClusterRoleWorker
	} else {
		node, err := utils.GetControlPlaneNode(ctx, s, s.Name())
		if err != nil {
			// The node is likely in a joining or leaving phase where the role is not yet settled.
			// Use the unknown role but still log this incident for debugging.
			log.Printf("Failed to check if node is control-plane. This is expected if the node is in a joining/leaving phase. %v", err)
			clusterRole = apiv1.ClusterRoleUnknown
		} else {
			if node != nil {
				return *node, nil
			}
		}

	}
	return apiv1.NodeStatus{
		Name:        s.Name(),
		Address:     s.Address().Hostname(),
		ClusterRole: clusterRole,
	}, nil

}
