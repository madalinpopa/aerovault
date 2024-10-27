package container

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
)

const (
	image      = "busybox"
	backupDir  = "/backup"
	backupTmpl = "%s-%d"
	tarCmdTmpl = "tar cvf %s/%s.tar %s"
	autoRemove = true
)

// BackupManager handles backup operations such as creating and inspecting container states.
type BackupManager struct {
	cli APIClient
	ctx context.Context
}

// NewBackupManager initializes and returns a new BackupManager with the provided APIClient and context.
func NewBackupManager(cli APIClient, ctx context.Context) *BackupManager {
	return &BackupManager{cli: cli, ctx: ctx}
}

// BackupVolume creates a backup of the specified volume in the given container and writes it to the specified output path.
func (bm *BackupManager) BackupVolume(container, volume, outputPath string) error {
	m, err := bm.getMountPoint(container, volume)
	if err != nil {
		return err
	}
	return bm.createBackupContainer(container, volume, m.Destination, outputPath)
}

// createBackupContainer creates a backup of the specified volume by creating a Docker container to tar its contents.
func (bm *BackupManager) createBackupContainer(volumeFrom, volumeName, destinationPath, hostPath string) error {
	cmd := generateTarCommand(volumeName, destinationPath)

	config, err := createContainerConfig(image, cmd)
	if err != nil {
		return err
	}

	hostConfig := &container.HostConfig{
		AutoRemove:  autoRemove,
		VolumesFrom: []string{volumeFrom},
		Binds:       []string{fmt.Sprintf("%s:/backup:rw", hostPath)},
	}

	cr, err := bm.cli.ContainerCreate(bm.ctx, config, hostConfig, nil, nil, "backup-"+volumeName)
	if err != nil {
		return fmt.Errorf("failed to create backup container: %w", err)
	}

	return bm.cli.ContainerStart(bm.ctx, cr.ID, container.StartOptions{})

}

// getMountPoint retrieves the mount point for a specified volume in a container.
func (bm *BackupManager) getMountPoint(containerName, volumeName string) (types.MountPoint, error) {

	c, err := bm.cli.ContainerInspect(bm.ctx, containerName)
	if err != nil {
		return types.MountPoint{}, fmt.Errorf("failed to inspect container %s: %w", containerName, err)
	}

	if len(c.Mounts) == 0 {
		return types.MountPoint{}, fmt.Errorf("no mounts found for container %s", containerName)
	}

	for _, m := range c.Mounts {
		if m.Name == volumeName {
			return m, nil
		}
	}

	return types.MountPoint{}, fmt.Errorf("no mount found for volume %s", volumeName)
}

// nowFunc returns the current time, used to generate timestamps for various operations. It can be overridden for testing purposes.
var nowFunc = time.Now

// generateTarCommand generates a tar command string for creating a tarball of a specified backup in a defined destination path.
func generateTarCommand(volumeName, destinationPath string) string {
	backupName := fmt.Sprintf(backupTmpl, volumeName, nowFunc().Unix())
	return fmt.Sprintf(tarCmdTmpl, backupDir, backupName, destinationPath)
}

// createContainerConfig generates a Docker container configuration object using the specified image and command.
// The function retrieves the user ID (uid) and group ID (gid) of the executing user, and constructs the configuration
// to ensure the backup file is created with these IDs.
func createContainerConfig(image, cmd string) (*container.Config, error) {
	uid, gid, err := getUserAndGroup()
	if err != nil {
		return nil, err
	}
	return &container.Config{
		Image: image,
		Tty:   false,
		User:  uid + ":" + gid,
		Cmd:   []string{"sh", "-c", cmd},
	}, nil
}

// getUID is a variable that holds the function os.Getuid, which retrieves the numeric user ID of the caller.
var getUID = os.Getuid

// getGID holds the function os.Getgid which returns the numeric group ID of the caller.
var getGID = os.Getgid

// getUserAndGroup retrieves the current user's UID and GID as strings and returns them along with any error encountered.
func getUserAndGroup() (string, string, error) {
	uid := strconv.Itoa(getUID())
	gid := strconv.Itoa(getGID())
	return uid, gid, nil
}
