package repo

import (
	"fmt"
	"github.com/Miguel-Dorta/gkup-core/pkg/repo/snapshot"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func NewSnapshot(snapshotName string, t time.Time) (*snapshot.Writer, error) {
	snapPath := filepath.Join(path, snapshot.SnapshotDirName)
	if snapshotName != "" {
		snapPath = filepath.Join(snapPath, snapshotName)
	}
	snapPath = filepath.Join(snapPath, strconv.FormatInt(t.Unix(), 10)+snapshot.Extension)

	w, err := snapshot.NewWriter(snapPath)
	if err != nil {
		return nil, fmt.Errorf("error creating snapshot %s: %w", snapPath, err)
	}
	return w, nil
}

func ListSnapshots() (map[string][]int64, error) {
	snapDir := filepath.Join(path, snapshot.SnapshotDirName)
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