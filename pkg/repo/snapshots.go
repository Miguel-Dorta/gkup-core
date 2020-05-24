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

// NewSnapshot creates a new snapshot with the args given, and returns its writer.
func NewSnapshot(snapGroup string, t time.Time, numberOfFiles uint64) (*snapshot.Writer, error) {
	snapPath := getPath(snapGroup, t)

	if err := os.MkdirAll(filepath.Dir(snapPath), 0755); err != nil {
		return nil, fmt.Errorf("error creating parent directories for new snapshot %s: %w", snapPath, err)
	}
	w, err := snapshot.NewWriter(snapPath, numberOfFiles)
	if err != nil {
		return nil, fmt.Errorf("error creating snapshot %s: %w", snapPath, err)
	}
	return w, nil
}

// OpenSnapshot opens a snapshot with the args given, and returns its reader.
func OpenSnapshot(snapGroup string, t time.Time) (*snapshot.Reader, error) {
	snapPath := getPath(snapGroup, t)
	r, err := snapshot.NewReader(snapPath)
	if err != nil {
		return nil, fmt.Errorf("error reading snapshot %s: %w", snapPath, err)
	}
	return r, nil
}

// ListSnapshots lists all the snapshots of a repository, returning a map where the key is the snapshot group,
// and the value is a slice of those snapshots in unix-time UTC.
func ListSnapshots() (map[string][]int64, error) {
	snapDir := filepath.Join(path, snapshotDirName)
	fs, err := ioutil.ReadDir(snapDir)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string][]int64{}, nil
		}
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

// listSnapChild returns a slice of the times of the snapshots contained in the directory provided. It is NOT recursive.
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

// getTime gets a valid filename of a snapshot and returns its time. It returns nil if the filename is invalid.
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

// getPath gets a group name and a time.Time, and returns the path of the corresponding snapshot.
func getPath(snapGroup string, t time.Time) string {
	snapPath := filepath.Join(path, snapshotDirName)
	if snapGroup != "" {
		snapPath = filepath.Join(snapPath, snapGroup)
	}
	return filepath.Join(snapPath, fmt.Sprintf("%d%s", t.UTC().Unix(), snapshot.Extension))
}
