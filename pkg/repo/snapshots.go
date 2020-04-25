package repo

import (
	"fmt"
	"github.com/Miguel-Dorta/gkup-core/pkg/repo/snapshot"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const snapshotDirName = "snapshots"

func NewSnapshot(snapshotName string, t time.Time, numberOfFiles uint64) (*snapshot.Writer, error) {
	snapPath := getPath(snapshotName, t)

	if err := os.MkdirAll(filepath.Dir(snapPath), 0755); err != nil {
		return nil, fmt.Errorf("error creating parent directories for new snapshot %s: %w", snapPath, err)
	}
	w, err := snapshot.NewWriter(snapPath, numberOfFiles)
	if err != nil {
		return nil, fmt.Errorf("error creating snapshot %s: %w", snapPath, err)
	}
	return w, nil
}

func OpenSnapshot(snapshotName string, t time.Time) (*snapshot.Reader, error) {
	snapPath := getPath(snapshotName, t)
	r, err := snapshot.NewReader(snapPath)
	if err != nil {
		return nil, fmt.Errorf("error reading snapshot %s: %w", snapPath, err)
	}
	return r, nil
}

func ListSnapshots() (map[string][]int64, error) {
	snapDir := filepath.Join(path, snapshotDirName)
	fs, err := ioutil.ReadDir(snapDir)
	if err != nil {
		return nil, fmt.Errorf("error listing directory %s: %w", snapDir, err)
	}

	m := make(map[string][]int64, len(fs))

	for _, f := range fs {
		fPath := filepath.Join(snapDir, f.Name())
		if f.IsDir() {
			childList, err := listSnapChild(fPath)
			if err != nil {
				return nil, err
			}
			m[f.Name()] = childList
			continue
		}

		if t := getTime(f.Name()); t != nil {
			m[""] = append(m[""], *t)
		}
	}
	return m, nil
}

func listSnapChild(path string) ([]int64, error) {
	fs, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("error listing directory %s: %w", path, err)
	}

	list := make([]int64, 0, len(fs))
	for _, f := range fs {
		if !f.Mode().IsRegular() {
			continue
		}
		if t := getTime(f.Name()); t != nil {
			list = append(list, *t)
		}
	}
	return list, nil
}

func getTime(filename string) *int64 {
	if !strings.HasSuffix(filename, snapshot.Extension) {
		return nil
	}
	i, err := strconv.ParseInt(strings.TrimSuffix(filename, snapshot.Extension), 10, 64)
	if err != nil {
		return nil
	}
	return &i
}

func getPath(snapshotName string, t time.Time) string {
	snapPath := filepath.Join(path, snapshotDirName)
	if snapshotName != "" {
		snapPath = filepath.Join(snapPath, snapshotName)
	}
	return filepath.Join(snapPath, fmt.Sprintf("%d%s", t.UTC().Unix(), snapshot.Extension))
}
