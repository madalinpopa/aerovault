package container

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"time"
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
	backupName := fmt.Sprintf(backupTmpl, volumeName, time.Now().Unix())
	cmd := generateTarCommand(backupName, destinationPath)

	config := &container.Config{
		Image: "busybox",
		Tty:   false,
		Cmd:   []string{"sh", "-c", cmd},
	}

	hostConfig := &container.HostConfig{
		AutoRemove:  true,
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

func generateTarCommand(backupName, destinationPath string) string {
	return fmt.Sprintf(tarCmdTmpl, backupDir, backupName, destinationPath)
}
