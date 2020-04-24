package repo_test

import (
	"fmt"
	"github.com/Miguel-Dorta/gkup-core/pkg/common"
	"github.com/Miguel-Dorta/gkup-core/pkg/hash"
	"github.com/Miguel-Dorta/gkup-core/pkg/repo"
	"github.com/Miguel-Dorta/gkup-core/pkg/repo/settings"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

var (
	nowNano = strconv.FormatInt(time.Now().UnixNano(), 10)
	repoPath = filepath.Join(os.TempDir(), "gkup-core_pkg_repo_fileTest_repo_" + nowNano)
	tmpFile = filepath.Join(os.TempDir(), "gkup-core_pkg_repo_fileTest_file_" + nowNano)
	h = hash.NewHasher(hash.SHA3_512, 128*1024)
	f *common.File
	testAddFileFailed = false
)

func prepareFile() error {
	size := 4*1024*1024
	data := make([]byte, size)
	rand.Seed(time.Now().UnixNano())
	rand.Read(data)
	if err := ioutil.WriteFile(tmpFile, data, 0777); err != nil {
		return fmt.Errorf("error writing tmp file: %s", err)
	}

	f = &common.File{
		AbsPath: tmpFile,
		RelPath: "test/dir/tmpFile",
		Size:    int64(size),
	}
	if err := h.HashFile(f); err != nil {
		return fmt.Errorf("error hashing tmpFile: %s", err)
	}
	return nil
}

func TestAddFile(t *testing.T) {
	if err := repo.Create(repoPath, &settings.Settings{HashAlgorithm: hash.SHA3_512}); err != nil {
		testAddFileFailed = true
		t.Fatalf("error creating repo: %s", err)
	}

	if err := prepareFile(); err != nil {
		testAddFileFailed = true
		os.RemoveAll(repoPath)
		os.Remove(tmpFile)
		t.Fatal(err)
	}

	if err := repo.AddFile(f); err != nil {
		testAddFileFailed = true
		os.RemoveAll(repoPath)
		t.Fatalf("failed to add file: %s", err)
	}

	os.Remove(tmpFile)

	fInRepo := &common.File{AbsPath: filepath.Join(repoPath, "files", f.Hash[:2], f.Hash[2:4], f.Hash[4:] + "-" + strconv.Itoa(4*1024*1024))}
	if err := h.HashFile(fInRepo); err != nil {
		testAddFileFailed = true
		os.RemoveAll(repoPath)
		t.Fatalf("error hashing file in repo: %s", err)
	}

	if f.Hash != fInRepo.Hash {
		testAddFileFailed = true
		os.RemoveAll(repoPath)
		t.Fatalf("hashes don't match:\n--> EXPECTED: %s\n--> FOUND:    %s", f.Hash, fInRepo.Hash)
	}
}

func TestRestoreFile(t *testing.T) {
	if testAddFileFailed {
		t.Skip("TestAddFile failed")
	}

	restorePath := filepath.Join(os.TempDir(), "gkup-core_pkg_repo_fileTest_restoreDir_" + nowNano)
	if err := repo.RestoreFile(f, restorePath); err != nil {
		os.RemoveAll(repoPath)
		t.Fatalf("error restoring file: %s", err)
	}
	os.RemoveAll(repoPath)

	fRestoredPath := filepath.Join(restorePath, filepath.FromSlash(f.RelPath))
	stat, err := os.Stat(fRestoredPath)
	if err != nil {
		os.RemoveAll(restorePath)
		t.Fatalf("error getting stat of restored file: %s", err)
	}

	if stat.Size() != f.Size {
		os.RemoveAll(restorePath)
		t.Fatalf("sizes don't match:\n--> EXPECTED: %d\n--> FOUND:    %d", f.Size, stat.Size())
	}

	fRestored := &common.File{AbsPath: fRestoredPath}
	if err := h.HashFile(fRestored); err != nil {
		os.RemoveAll(restorePath)
		t.Fatalf("error hashing restored file: %s", err)
	}

	if fRestored.Hash != f.Hash {
		os.RemoveAll(restorePath)
		t.Fatalf("hashes don't match:\n--> EXPECTED: %s\n--> FOUND:    %s", f.Hash, fRestored.Hash)
	}
	os.RemoveAll(restorePath)
}
