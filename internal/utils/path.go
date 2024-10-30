package utils

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const (
	// errOutputPathNotExist indicates that the specified output path does not exist.
	errOutputPathNotExist = "output path does not exist"

	// errOutputPathNotDir indicates that the specified output path is not a directory.
	errOutputPathNotDir = "output path is not a directory"

	// errGettingOutputPathInfo is used to signal that there was an error retrieving information about the output path.
	errGettingOutputPathInfo = "error getting output path info: %v"
)

// ValidateOutputPath checks if the given output directory path is valid, exists, and is a directory.
// Returns the absolute path if valid, otherwise returns an error.
func ValidateOutputPath(output string) (string, error) {
	info, err := os.Stat(output)
	if err != nil {
		if os.IsNotExist(err) {
			return "", errors.New(errOutputPathNotExist)
		}
		return "", fmt.Errorf(errGettingOutputPathInfo, err)
	}
	if !info.IsDir() {
		return "", errors.New(errOutputPathNotDir)
	}
	absPath, err := filepath.Abs(output)
	if err != nil {
		return "", fmt.Errorf("error getting absolute path: %v", err)
	}
	return absPath, nil
}

// GetCurrentDir retrieves and returns the current working directory. Returns an empty string and an error if any occurs.
func GetCurrentDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error getting current dir: %v", err)
	}
	return dir, nil
}

// GetResolvedOutputPath resolves the provided directory path; retrieves the current working directory if the path is empty.
func GetResolvedOutputPath(path string) (string, error) {
	if path == "" {
		return handleEmptyPath()
	}
	return handleProvidedPath(path)
}

// handleEmptyPath retrieves and returns the current working directory. If it fails, it returns an empty string and an error.
func handleEmptyPath() (string, error) {
	currentDir, err := GetCurrentDir()
	if err != nil {
		return "", err
	}
	return currentDir, nil
}

// handleProvidedPath validates the provided directory path and returns its absolute path if valid, otherwise returns an error.
func handleProvidedPath(path string) (string, error) {
	validatedPath, err := ValidateOutputPath(path)
	if err != nil {
		return "", err
	}
	return validatedPath, nil
}
