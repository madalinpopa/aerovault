package utils_test

import (
	"os"
	"testing"

	"github.com/madalinpopa/aerovault/internal/utils"
)

const (
	tempDirPattern = "container-volume-backup-test-*"
)

func createTempFile(t *testing.T) *os.File {
	t.Helper()
	f, err := os.CreateTemp(os.TempDir(), tempDirPattern)
	if err != nil {
		t.Fatalf("failed to create temp file: %s", err)
	}
	return f
}

func removeTempFile(t *testing.T, fileName string) {
	t.Helper()
	err := os.Remove(fileName)
	if err != nil {
		t.Fatalf("failed to remove temp file: %s", err)
	}
}

func createTempD(t *testing.T) string {
	t.Helper()
	tempD, err := os.MkdirTemp(os.TempDir(), tempDirPattern)
	if err != nil {
		t.Fatalf("failed to create temp dir: %s", err)
	}
	return tempD
}

func removeTempD(t *testing.T, dirName string) {
	t.Helper()
	err := os.RemoveAll(dirName)
	if err != nil {
		t.Fatalf("failed to remove temp dir: %s", err)
	}
}

func TestValidateOutputPath_pathNotExists(t *testing.T) {

	_, err := utils.ValidateOutputPath("not-exists")
	if err == nil {
		t.Errorf("expected error, got nil")
	}

	if err != nil && err.Error() != "output path does not exist" {
		t.Errorf("expected error message: 'output path not exists', got: %s", err.Error())
	}
}

func TestValidateOutputPath_pathExists(t *testing.T) {

	// create a temporary directory
	tempD := createTempD(t)
	defer removeTempD(t, tempD)

	_, err := utils.ValidateOutputPath(tempD)
	if err != nil {
		t.Errorf("expected no error, got %s", err)
	}

}

func TestValidateOutputPath_pathIsFile(t *testing.T) {
	// create a temp file
	f := createTempFile(t)
	defer removeTempFile(t, f.Name())

	_, err := utils.ValidateOutputPath(f.Name())
	if err == nil {
		t.Errorf("expected error, got nil: %s", err)
	}

	if err.Error() != "output path is not a directory" {
		t.Errorf("expected error message: 'output path is a file', got: %s", err.Error())
	}
}

func TestGetCurrenDir(t *testing.T) {
	dir, err := utils.GetCurrentDir()
	if err != nil {
		t.Errorf("expected no error, got %s", err)
	}
	if dir == "" {
		t.Errorf("expected dir, got empty string")
	}

	osDir, osErr := os.Getwd()
	if osErr != nil {
		t.Fatalf("os.Getwd failed: %v", osErr)
	}

	if dir != osDir {
		t.Errorf("GetCurrentDir returned %s, but os.Getwd returned %s", dir, osDir)
	}
}

func TestGetResolvedOutputPath_validPath(t *testing.T) {
	// create a temporary directory
	tempD := createTempD(t)
	defer removeTempD(t, tempD)

	p, err := utils.GetResolvedOutputPath(tempD)
	if err != nil {
		t.Errorf("expected no error, got %s", err)
	}

	if p != tempD {
		t.Errorf("expected %s, got %s", tempD, p)
	}
}

func TestGetResolvedOutputPath_emptyPath(t *testing.T) {
	osDir, osErr := os.Getwd()
	if osErr != nil {
		t.Fatalf("os.Getwd failed: %v", osErr)
	}

	p, err := utils.GetResolvedOutputPath("")
	if err != nil {
		t.Errorf("expected no error, got %s", err)
	}
	if p != osDir {
		t.Errorf("expected string %s, got %s", osDir, p)
	}
}
