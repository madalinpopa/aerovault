package dockerbackup

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

// mockTimeNow function to return a fixed time for testing
var mockTimeNow = func() time.Time {
	return time.Unix(1609459200, 0) // Fixed timestamp (January 1, 2021)
}

// mockUID is a function that returns a mock user ID for testing purposes.
var mockUID = func() int {
	return 1001 // Mocked UID
}

// mockGID is a mock function that always returns a fixed GID value of 1002.
var mockGID = func() int {
	return 1002 // Mocked GID
}

// APIClientStub is a stub implementation of the ContainerManager interface for testing purposes.
type APIClientStub struct{}

// ContainerInspect retrieves detailed information about a container specified by its containerID.
func (api *APIClientStub) ContainerInspect(_ context.Context, containerID string) (types.ContainerJSON, error) {
	if containerID == "nginx" {
		return types.ContainerJSON{
			Mounts: []types.MountPoint{
				{
					Name:        "nginx",
					Destination: "/var/www/data",
				},
			},
		}, nil
	}
	return types.ContainerJSON{}, nil
}

// ContainerCreate creates a new container with the provided configuration and returns a creation response or an error.
func (api *APIClientStub) ContainerCreate(_ context.Context, _ *container.Config, _ *container.HostConfig, _ *network.NetworkingConfig, _ *ocispec.Platform, _ string) (container.CreateResponse, error) {
	return container.CreateResponse{}, nil
}

// ContainerStart starts an existing container based on the provided container ID and start options.
func (api *APIClientStub) ContainerStart(_ context.Context, _ string, _ container.StartOptions) error {
	return nil
}

// TestNewBackupManager tests the creation of a new BackupManager instance with a stubbed APIClient.
func TestNewBackupManager(t *testing.T) {
	ctx := context.Background()
	cms := &APIClientStub{}
	dcm := NewBackupManager(cms, ctx)
	if dcm == nil {
		t.Errorf("expected DockerContainerManager instance, got nil")
	}
}

// TestGetMountPoint_noMountsFound tests the GetMountPoint method when no mounts are found for the container.
func TestGetMountPoint_noMountsFound(t *testing.T) {
	ctx := context.Background()
	cli := &APIClientStub{}
	bm := NewBackupManager(cli, ctx)

	_, err := bm.getMountPoint("container", "volume")
	if err == nil {
		t.Errorf("expected error, got nil %s", err)
	}
}

// TestGetMountPoint_noVolumeFound tests the GetMountPoint function when no volume is found in the container's mounts.
func TestGetMountPoint_noVolumeFound(t *testing.T) {
	ctx := context.Background()
	cli := &APIClientStub{}
	bm := NewBackupManager(cli, ctx)

	_, err := bm.getMountPoint("nginx", "volume")
	if err == nil {
		t.Errorf("expected error, got nil: %s", err)
	}
}

// TestGetMountPoint_volumeFound verifies that the GetMountPoint method correctly finds the specified volume in a container inspection.
func TestGetMountPoint_volumeFound(t *testing.T) {
	ctx := context.Background()
	cli := &APIClientStub{}
	bm := NewBackupManager(cli, ctx)

	m, err := bm.getMountPoint("nginx", "nginx")
	if err != nil {
		t.Errorf("expected no error, got %s", err)
	}

	if m.Name != "nginx" {
		t.Errorf("expected mount point name to be 'nginx', got %s", m.Name)
	}
	if m.Destination != "/var/www/data" {
		t.Errorf("expected mount point destination to be '/var/www/data', got %s", m.Destination)
	}
}

// TestGenerateTarCommand tests the generateTarCommand function
func TestGenerateTarCommand(t *testing.T) {

	volumeName := "test_volume"
	destinationPath := "/tmp/destination"

	expectedBackupName := fmt.Sprintf(backupTmpl, volumeName, mockTimeNow().Unix())
	expectedCommand := fmt.Sprintf(tarCmdTmpl, backupDir, expectedBackupName, destinationPath)

	nowFunc = mockTimeNow
	actualCommand := generateTarCommand(volumeName, destinationPath)

	if actualCommand != expectedCommand {
		t.Errorf("Expected command %s, but got %s", expectedCommand, actualCommand)
	}
}

func TestGetUserAndGroup(t *testing.T) {
	// Override getUID and getGID for the test
	getUID = mockUID
	getGID = mockGID

	defer func() {
		// Restore original behavior after the test
		getUID = os.Getuid
		getGID = os.Getgid
	}()

	// Arrange
	expectedUID := strconv.Itoa(mockUID())
	expectedGID := strconv.Itoa(mockGID())

	// Act
	actualUID, actualGID, err := getUserAndGroup()

	// Assert
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}
	if actualUID != expectedUID {
		t.Errorf("Expected UID %s, but got %s", expectedUID, actualUID)
	}
	if actualGID != expectedGID {
		t.Errorf("Expected GID %s, but got %s", expectedGID, actualGID)
	}
}
