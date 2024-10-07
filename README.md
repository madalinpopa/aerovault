# AeroVault

**AeroVault** is a CLI tool designed to easily backup, restore, and sync Docker volumes to popular cloud storage services like AWS and Azure. With future plans to include features like volume restoration and syncing, AeroVault simplifies the process of managing container volume backups in the cloud.

## Features

- [x] **Backup**: Create backups of Docker container volumes.
- [ ] **Restore**: Restore Docker container volumes from backups.
- [ ] **Cloud Provider Support**: Backup to multiple cloud providers, including Azure and AWS.
  - [ ] **Azure Storage Account**
- [ ] **Sync**: Sync Docker volume backups across different cloud storage services.

## Usage

**Backup Volume**

```bash
aero backup -c my-container -v my-volume
```
