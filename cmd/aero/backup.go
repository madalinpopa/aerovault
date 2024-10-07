package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/client"
	"github.com/madalinpopa/aerovault/container"
	"github.com/madalinpopa/aerovault/internal/utils"
	"github.com/spf13/cobra"
	"log"
	"os"
)

// backupCmd represents the command to create a backup tar file for a given container and its volume.
var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Create a backup tar file for given container container",
	Run: func(cmd *cobra.Command, args []string) {
		containerName := getStringFlag(cmd, "container")
		volumeName := getStringFlag(cmd, "volume")
		outputPath := getStringFlag(cmd, "output")

		if err := backup(containerName, volumeName, outputPath); err != nil {
			_, err := fmt.Fprintf(os.Stderr, "Backup failed: %v\n", err)
			if err != nil {
				return
			}
			os.Exit(1)
		}
	},
}

// init initializes the backup command by setting up flags and marking required ones. Adds the command to rootCmd.
func init() {
	var containerName string
	var volumeName string
	var outputPath string

	backupCmd.Flags().StringVarP(&containerName, "container", "c", "", "Container name (required)")
	markFlagRequired(backupCmd, "container")
	backupCmd.Flags().StringVarP(&volumeName, "volume", "v", "", "Volume name (required)")
	markFlagRequired(backupCmd, "volume")
	backupCmd.Flags().StringVarP(&outputPath, "output", "o", ".", "Output path")

	rootCmd.AddCommand(backupCmd)
}

// backup creates a backup of the specified volume in a given container and writes it to the specified output path.
// Takes containerName as the name of the container, volumeName as the name of the volume, and outputPath as the output file path.
// Returns an error if the backup operation fails.
func backup(containerName string, volumeName string, outputPath string) error {
	outputPath, err := utils.GetResolvedOutputPath(outputPath)
	if err != nil {
		return err
	}

	cli, err := createDockerClient()
	if err != nil {
		return fmt.Errorf("failed to create Docker client: %v", err)
	}
	defer closeDockerClient(cli)

	ctx := context.Background()
	bm := container.NewBackupManager(cli, ctx)
	return bm.BackupVolume(containerName, volumeName, outputPath)
}

// getStringFlag retrieves the string value of the specified flag from the given command.
// It exits the program if an error occurs while fetching the flag.
func getStringFlag(cmd *cobra.Command, name string) string {
	value, err := cmd.Flags().GetString(name)
	if err != nil {
		if _, err := fmt.Fprintf(os.Stderr, "Error getting %s flag: %v\n", name, err); err != nil {
			fmt.Println("Failed to print error message")
		}
		os.Exit(1)
	}
	return value
}

// markFlagRequired marks a flag as required for a given Cobra command. Logs and exits on error.
func markFlagRequired(cmd *cobra.Command, name string) {
	if err := cmd.MarkFlagRequired(name); err != nil {
		log.Fatal("Failed to mark flag as required")
	}
}

// createDockerClient creates and returns a Docker client using environment variables and API version negotiation.
func createDockerClient() (*client.Client, error) {
	return client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
}

// closeDockerClient closes the provided Docker client and logs an error if the close operation fails.
func closeDockerClient(cli *client.Client) {
	if err := cli.Close(); err != nil {
		fmt.Println("Error closing Docker client:", err)
	}
}
