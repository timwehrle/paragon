package app

import (
	"encoding/json"
	"fmt"
	"syscall"

	"golang.org/x/sys/windows"
)

const bufferSize int = 256

type VolumeInfo struct {
	Drive        string `json:"drive"`
	Type         string `json:"type"`
	VolumeLabel  string `json:"volumeLabel"`
	FileSystem   string `json:"fileSystem"`
	SerialNumber string `json:"serialNumber"`
	MaxComponentSize uint32 `json:"maxComponentSize"`
	FileSystemFlags uint32 `json:"fileSystemFlags"`
}

func (a *App) ListVolumes() ([]VolumeInfo, error) {
	var volumes []VolumeInfo

	// Buffer for GetLogicalDriveStrings (e.g., C:\, D:\, E:\)
	buffer := make([]uint16, bufferSize)
	size, err := windows.GetLogicalDriveStrings(uint32(len(buffer)), &buffer[0])
	if err != nil {
		return nil, fmt.Errorf("failed to get logical drives: %v", err)
	}

	// Iterate through the drive letters
	for i := 0; i < int(size); i += 4 {
		if buffer[i] == 0 {
			break
		}
		drive := syscall.UTF16ToString(buffer[i : i+4])
		volume, err := getVolumeInfo(drive)
		if err != nil {
			return nil, fmt.Errorf("error fetching volume info for drive %s: %v", drive, err)
		}

		volumes = append(volumes, volume)
	}

	return volumes, nil
}

func getVolumeInfo(drive string) (VolumeInfo, error) {
	volume := VolumeInfo{Drive: drive}

	// Get the drive type (removable, fixed, etc.)
	drivePtr, err := syscall.UTF16PtrFromString(drive)
	if err != nil {
		return volume, fmt.Errorf("failed to convert drive string to UTF16: %v", err)
	}
	volume.Type = getDriveType(drivePtr)

	// Get volume information (label, filesystem, etc.)
	var volumeName [bufferSize]uint16
	var fileSystemName [bufferSize]uint16
	var serialNumber, maxComponentLength, fileSystemFlags uint32

	err = windows.GetVolumeInformation(
		drivePtr, 
		&volumeName[0], 
		uint32(len(volumeName)), 
		&serialNumber, 
		&maxComponentLength, 
		&fileSystemFlags, 
		&fileSystemName[0], 
		uint32(len(fileSystemName)),
	)

	if err != nil {
		return volume, fmt.Errorf("failed to get volume information for drive %s: %v", drive, err)
	}

	volume.VolumeLabel = syscall.UTF16ToString(volumeName[:])
	volume.FileSystem = syscall.UTF16ToString(fileSystemName[:])
	volume.SerialNumber = fmt.Sprintf("%X", serialNumber)
	volume.MaxComponentSize = maxComponentLength
	volume.FileSystemFlags = fileSystemFlags

	return volume, nil
}

func getDriveType(drivePtr *uint16) string {
	switch windows.GetDriveType(drivePtr) {
	case windows.DRIVE_FIXED:
		return "Fixed drive (HDD/SSD)"
	case windows.DRIVE_REMOVABLE:
		return "Removable drive (USB)"
	case windows.DRIVE_CDROM:
		return "CD-ROM"
	case windows.DRIVE_REMOTE: 
		return "Network drive"
	default:
		return "Unknown"
	}
}

// ListVolumesJSON returns the list of volumes in JSON format
func (a *App) ListVolumesJSON() (string, error) {
	volumes, err := a.ListVolumes()
	if err != nil {
		return "", err
	}

	jsonData, err := json.MarshalIndent(volumes, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal volumes: %v", err)
	}

	return string(jsonData), nil
}
