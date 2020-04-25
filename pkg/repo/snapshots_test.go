package repo_test

import (
	"fmt"
	"github.com/Miguel-Dorta/gkup-core/pkg/hash"
	"github.com/Miguel-Dorta/gkup-core/pkg/repo"
	"github.com/Miguel-Dorta/gkup-core/pkg/repo/settings"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

var (
	snapTestRepoPath = filepath.Join(os.TempDir(), "gkup-core_pkg_repo_snapshotTest_" + strconv.FormatInt(time.Now().UnixNano(), 10))
	t1 = time.Now()
	testNewSnapshotFailed = false
)

func TestNewSnapshot(t *testing.T) {
	if err := prepare(); err != nil {
		testNewSnapshotFailed = true
		t.Fatalf("error creating snapshot: %s", err)
	}

	if err := emptySnapshot("", t1); err != nil {
		testNewSnapshotFailed = true
		t.Fatalf("error creating snapshot with no name: %s", err)
	}
	if err := emptySnapshot("test", t1); err != nil {
		testNewSnapshotFailed = true
		t.Fatalf("error creating snapshot with name: %s", err)
	}

	if _, err := os.Stat(filepath.Join(snapTestRepoPath, "snapshots", strconv.FormatInt(t1.UTC().Unix(), 10) + ".gkup")); err != nil {
		t.Errorf("error creating snapshot with empty name: %s", err)
	}
	if _, err := os.Stat(filepath.Join(snapTestRepoPath, "snapshots", "test", strconv.FormatInt(t1.UTC().Unix(), 10) + ".gkup")); err != nil {
		t.Errorf("error creating snapshot with name: %s", err)
	}
	testNewSnapshotFailed = t.Failed()
}

func TestOpenSnapshot(t *testing.T) {
	if testNewSnapshotFailed {
		t.Skip("TestOpenSnapshot failed")
	}

	if r, err := repo.OpenSnapshot("", t1); err != nil {
		t.Errorf("error opening snapshot with no name: %s", err)
	} else {
		r.Close()
	}

	if r, err := repo.OpenSnapshot("test", t1); err != nil {
		t.Errorf("error opening snapshot with name: %s", err)
	} else {
		r.Close()
	}
}

func TestListSnapshots(t *testing.T) {
	if testNewSnapshotFailed {
		t.Skip("TestOpenSnapshot failed")
	}

	snaps, err := repo.ListSnapshots()
	if err != nil {
		t.Fatalf("error listing snapshots: %s", err)
	}

	if len(snaps) != 2 {
		t.Fatalf("not the number of snapshots expected: %+v", snaps)
	}

	if (snaps[""][0] != t1.UTC().Unix()) || (snaps["test"][0] != t1.UTC().Unix()) {
		t.Fatalf("times and names doesn't match: %+v", snaps)
	}
}

func prepare() error {
	os.RemoveAll(snapTestRepoPath)
	return repo.Create(snapTestRepoPath, &settings.Settings{HashAlgorithm: hash.SHA1})
}

func emptySnapshot(name string, time time.Time) error {
	w, err := repo.NewSnapshot(name, time, 0)
	if err != nil {
		return fmt.Errorf("error creating snapshot: %s", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("error closing snapshot: %s", err)
	}
	return nil
}
