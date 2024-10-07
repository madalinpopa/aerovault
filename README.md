# AeroVault

**AeroVault** is a CLI tool designed to easily backup, restore, and sync Docker volumes to popular cloud storage services like AWS and Azure. With future plans to include features like volume restoration and syncing, AeroVault simplifies the process of managing container volume backups in the cloud.

## Features

- [x] **Backup**: Create backups of Docker container volumes.
- [ ] **Restore**: Restore Docker container volumes from backups.
- [ ] **Cloud Provider Support**: Backup to multiple cloud providers, including Azure and AWS.
  - [ ] **Azure Storage Account**
- [ ] **Sync**: Sync Docker volume backups across different cloud storage services.

## Usage

1. Install **AeroVault** by cloning the repository or downloading the latest release.
2. Ensure you have Docker and cloud credentials (AWS, Azure) set up.
3. Run the following command to back up a Docker volume:

```bash
aero backup --container <container-name> --volume <volume-name> --output <save-path>
```
