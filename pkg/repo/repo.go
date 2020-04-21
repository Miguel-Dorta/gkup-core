package repo

import (
	"fmt"
	"github.com/Miguel-Dorta/gkup-core/pkg/common"
	"github.com/Miguel-Dorta/gkup-core/pkg/fileUtils"
	"github.com/Miguel-Dorta/gkup-core/pkg/repo/settings"
	"github.com/Miguel-Dorta/gkup-core/pkg/repo/snapshot"
	"os"
	"path/filepath"
	"time"
)

var (
	Sett *settings.Settings
	path string
)

func Init(repoPath string) error {
	path = repoPath

	s, err := settings.Load(filepath.Join(path, settings.Filename))
	if err != nil {
		return fmt.Errorf("error loading settings: %w", err)
	}
	Sett = s
	return nil
}

func AddFile(f *common.File) error {
	if f.Hash == "" {
		return fmt.Errorf("file %s hasn't been hashed", f.Path)
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

	if err := fileUtils.CopyFile(f.Path, fRepoPath); err != nil {
		return fmt.Errorf("error copying file from %s to %s: %w", f.Path, fRepoPath, err)
	}
	return nil
}

func NewSnapshot(snapshotName string, t time.Time) (*snapshot.Writer, error) {
	snapPath := filepath.Join(path, snapshot.SnapshotDirName)
	if snapshotName != "" {
		snapPath = filepath.Join(snapPath, snapshotName)
	}
	snapPath = filepath.Join(snapPath, fmt.Sprintf("%d.gkup", t.Unix()))

	w, err := snapshot.NewWriter(snapPath)
	if err != nil {
		return nil, fmt.Errorf("error creating snapshot %s: %w", snapPath, err)
	}
	return w, nil
}

func getFilePathInRepo(f *common.File) string {
	return filepath.Join(path, f.Hash[:2], f.Hash[2:4], fmt.Sprintf("%s-%d", f.Hash[4:], f.Size))
}
