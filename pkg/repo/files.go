package repo

import (
	"fmt"
	"github.com/Miguel-Dorta/gkup-core/pkg/common"
	"github.com/Miguel-Dorta/gkup-core/pkg/fileUtils"
	"os"
	"path/filepath"
)

func AddFile(f *common.File) error {
	if f.Hash == "" {
		return fmt.Errorf("file %s hasn't been hashed", f.AbsPath)
	}
	fRepoPath := getFilePathInRepo(f)

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

func RestoreFile(f *common.File, destination string) error {
	destination = filepath.Join(destination, filepath.FromSlash(f.RelPath))

	if err := os.MkdirAll(filepath.Dir(destination), 0755); err != nil {
		return fmt.Errorf("error making parent directory for file %s: %w", destination, err)
	}

	if err := fileUtils.CopyFile(getFilePathInRepo(f), destination); err != nil {
		return fmt.Errorf("error copying file from %s to %s: %w", getFilePathInRepo(f), destination, err)
	}
	return nil
}

func getFilePathInRepo(f *common.File) string {
	return filepath.Join(path, f.Hash[:2], f.Hash[2:4], fmt.Sprintf("%s-%d", f.Hash[4:], f.Size))
}
