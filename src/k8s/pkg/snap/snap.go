package snap

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/canonical/k8s/pkg/utils"
)

func Path(parts ...string) string {
	return filepath.Join(append([]string{os.Getenv("SNAP")}, parts...)...)
}

func DataPath(parts ...string) string {
	return filepath.Join(append([]string{os.Getenv("SNAP_DATA")}, parts...)...)
}
func CommonPath(parts ...string) string {
	return filepath.Join(append([]string{os.Getenv("SNAP_COMMON")}, parts...)...)
}

// StartService starts a k8s service. The name can be either prefixed or not.
func StartService(ctx context.Context, name string) error {
	return utils.RunCommand(ctx, "snapctl", "start", serviceName(name))
}

// StopService stops a k8s service. The name can be either prefixed or not.
func StopService(ctx context.Context, name string) error {
	return utils.RunCommand(ctx, "snapctl", "stop", serviceName(name))
}

// WaitServiceActive waits until a snap service reports to be active.
// or a timeout is reached.
func WaitServiceActive(ctx context.Context, name string) error {
	startTime := time.Now()
	timeout := 30 * time.Second

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("context cancelled")
		default:
		}

		cmd := exec.Command("snapctl", "services", serviceName(name))
		output, err := cmd.CombinedOutput()

		if err != nil {
			return fmt.Errorf("failed to check service status: %w", err)
		}
		// Check if the output contains "active"
		if strings.Contains(string(output), "active") {
			return nil
		}

		if time.Since(startTime) >= timeout {
			break
		}

		// Wait for 1 second before retrying
		time.Sleep(time.Second)
	}

	return fmt.Errorf("service %q did not become active within time.", name)
}

// serviceName infers the name of the snapctl daemon from the service name.
// if the serviceName is the snap name `k8s` (=referes to all services) it will return it as is.
func serviceName(serviceName string) string {
	if strings.HasPrefix(serviceName, "k8s.") || serviceName == "k8s" {
		return serviceName
	}
	return fmt.Sprintf("k8s.%s", serviceName)
}
