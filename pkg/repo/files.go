package repo

import (
	"fmt"
	"github.com/Miguel-Dorta/gkup-core/pkg/common"
	"github.com/Miguel-Dorta/gkup-core/pkg/fileUtils"
	"os"
	"path/filepath"
)

const filesDirName = "files"

// AddFiles add the file provided to the repository.
func AddFile(f *common.File) error {
	if err := f.CheckAbsPath(); err != nil {
		return err
	}

	fRepoPath, err := getFilePathInRepo(f)
	if err != nil {
		return err
	}

	exists, err := fileUtils.FileExists(fRepoPath)
	if err != nil {
		return fmt.Errorf("error checking existence of file %s: %w", fRepoPath, err)
	}
	if exists {
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(fRepoPath), 0755); err != nil {
		return fmt.Errorf("error creating parent directories for file %s: %w", fRepoPath, err)
	}

	if err := fileUtils.CopyFile(f.AbsPath, fRepoPath); err != nil {
		return fmt.Errorf("error copying file from %s to %s: %w", f.AbsPath, fRepoPath, err)
	}
	return nil
}

// RestoreFile takes the file provided and restores it in the destination path provided.
func RestoreFile(f *common.File, destination string) error {
	if err := f.CheckRelPath(); err != nil {
		return err
	}

	pathInRepo, err := getFilePathInRepo(f)
	if err != nil {
		return err
	}

	destination = filepath.Join(destination, filepath.FromSlash(f.RelPath))

	if err := os.MkdirAll(filepath.Dir(destination), 0755); err != nil {
		return fmt.Errorf("error making parent directory for file %s: %w", destination, err)
	}

	if err := fileUtils.CopyFile(pathInRepo, destination); err != nil {
		return fmt.Errorf("error copying file from %s to %s: %w", getFilePathInRepo(f), destination, err)
	}
	return nil
}

// getFilePathInRepo returns the real path that the provided object should have if stored in the repository.
// It returns errors if hash or size are empty.
func getFilePathInRepo(f *common.File) (string, error) {
	if err := f.CheckHash(); err != nil {
		return "", err
	}
	if err := f.CheckSize(); err != nil {
		return "", err
	}
	return filepath.Join(path, filesDirName, f.Hash[:2], f.Hash[2:4], fmt.Sprintf("%s-%d", f.Hash[4:], f.Size)), nil
}
