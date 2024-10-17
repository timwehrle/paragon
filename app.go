package main

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/sys/windows"
	"syscall"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called at application startup
func (a *App) startup(ctx context.Context) {
	// Perform your setup here
	a.ctx = ctx
}

// domReady is called after front-end resources have been loaded
func (a App) domReady(ctx context.Context) {
	// Add your action here
}

// beforeClose is called when the application is about to quit,
// either by clicking the window close button or calling runtime.Quit.
// Returning true will cause the application to continue, false will continue shutdown as normal.
func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	return false
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {
	// Perform your teardown here
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

type VolumeInfo struct {
	Drive        string `json:"drive"`
	Type         string `json:"type"`
	VolumeLabel  string `json:"volumeLabel"`
	FileSystem   string `json:"fileSystem"`
	SerialNumber string `json:"serialNumber"`
}

func (a *App) ListVolumes() ([]VolumeInfo, error) {
	var volumes []VolumeInfo

	// Buffer for GetLogicalDriveStrings (e.g., C:\, D:\, E:\)
	buffer := make([]uint16, 256)
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
		volume := VolumeInfo{Drive: drive}

		// Get the drive type (removable, fixed, etc.)
		driveType := windows.GetDriveType(syscall.StringToUTF16Ptr(drive))
		switch driveType {
		case windows.DRIVE_FIXED:
			volume.Type = "Fixed drive (HDD/SSD)"
		case windows.DRIVE_REMOVABLE:
			volume.Type = "Removable drive (USB)"
		case windows.DRIVE_CDROM:
			volume.Type = "CD-ROM"
		case windows.DRIVE_REMOTE:
			volume.Type = "Network drive"
		default:
			volume.Type = "Unknown"
		}

		// Get volume information (label, filesystem, etc.)
		var volumeName [256]uint16
		var fileSystemName [256]uint16
		var serialNumber, maxComponentLength, fileSystemFlags uint32

		err := windows.GetVolumeInformation(
			syscall.StringToUTF16Ptr(drive),
			&volumeName[0],
			uint32(len(volumeName)),
			&serialNumber,
			&maxComponentLength,
			&fileSystemFlags,
			&fileSystemName[0],
			uint32(len(fileSystemName)),
		)
		if err == nil {
			volume.VolumeLabel = syscall.UTF16ToString(volumeName[:])
			volume.FileSystem = syscall.UTF16ToString(fileSystemName[:])
		}

		volumes = append(volumes, volume)
	}

	return volumes, nil
}

func (a *App) ListVolumesJSON() (string, error) {
	volumes, err := a.ListVolumes()
	if err != nil {
		return "", err
	}

	jsonData, err := json.Marshal(volumes)
	if err != nil {
		return "", fmt.Errorf("failed to marshal volumes: %v", err)
	}

	return string(jsonData), nil
}
