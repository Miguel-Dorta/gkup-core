package repo

import (
	"github.com/Miguel-Dorta/gkup-core/internal"
	"github.com/Miguel-Dorta/gkup-core/pkg/hash"
	"github.com/Miguel-Dorta/gkup-core/pkg/repo/settings"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

var (
	repoPath         = filepath.Join(os.TempDir(), "gkup-core_pkg_repo_repoTest_"+strconv.FormatInt(time.Now().UnixNano(), 10))
	testCreateFailed = false
)

func init() {
	internal.Version = "v0.1.0-test"
}

func TestCreate(t *testing.T) {
	if err := Create(repoPath, &settings.Settings{HashAlgorithm: hash.SHA3_256}); err != nil {
		testCreateFailed = true
		os.RemoveAll(repoPath)
		t.Fatalf("error creating repo: %s", err)
	}
}

func TestInit(t *testing.T) {
	if testCreateFailed {
		t.Skip("TestCreate failed")
	}

	if err := Init(repoPath); err != nil {
		os.RemoveAll(repoPath)
		t.Fatalf("error initializing package repo: %s", err)
	}

	if path != repoPath {
		t.Errorf("path was not set to the expected one:\n--> Expected: %s\n--> Found:    %s", repoPath, path)
	}
	if (Sett.Version != internal.Version) || (Sett.HashAlgorithm != hash.SHA3_256) {
		t.Errorf("settings doesn't match:\n--> Expected: %+v\n--> Found:    %+v", settings.Settings{
			Version:       internal.Version,
			HashAlgorithm: hash.SHA3_256,
		}, *Sett)
	}
}
