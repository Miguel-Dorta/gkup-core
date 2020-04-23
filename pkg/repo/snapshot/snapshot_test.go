package snapshot_test

import (
	"github.com/Miguel-Dorta/gkup-core/internal"
	"github.com/Miguel-Dorta/gkup-core/pkg/common"
	"github.com/Miguel-Dorta/gkup-core/pkg/repo/snapshot"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

var (
	writeTestFailed bool
	fPath = filepath.Join(os.TempDir(), "gkup-core_pkg_repo_snapshot_Test_" + strconv.FormatInt(time.Now().UnixNano(), 10))
	fileList = []*common.File{
	{
			AbsPath: "C:\\User\\a",
			RelPath: "a",
			Hash: "0123456789abcdef",
			Size: 0,
		},
	{
			AbsPath: "C:\\User\\a\\b",
			RelPath: "a/b",
			Hash: "fedcba9876543210",
			Size: 1,
		},
	{
			AbsPath: "C:\\Programs\\c",
			RelPath: "c",
			Hash: "784328432289bdfe",
			Size: 2000000,
		},
	{
			AbsPath: "C:\\User\\a\\b\\d",
			RelPath: "d",
			Hash: "30278952",
			Size: 7500000,
		},
	{
			AbsPath: "D:\\a\\b\\c\\d\\e",
			RelPath: "a/b/c/d/e",
			Hash: "932482798bafe",
			Size: (1 << 63) - 1,
		},
	}
)

func init() {
	internal.Version = "v1.0.0"
}

func TestWriter(t *testing.T) {
	w, err := snapshot.NewWriter(fPath, uint64(len(fileList)))
	if err != nil {
		writeTestFailed = true
		t.Fatalf("error creating writer: %s", err)
	}
	for _, f := range fileList {
		if err := w.Write(f); err != nil {
			writeTestFailed = true
			os.Remove(fPath)
			t.Fatalf("error writing file %v: %s", f, err)
		}
	}

	if err := w.Close(); err != nil {
		writeTestFailed = true
		os.Remove(fPath)
		t.Fatalf("error closing writer: %s", err)
	}
}

func TestReader(t *testing.T) {
	if writeTestFailed {
		t.Skip("TestWriter failed")
	}

	r, err := snapshot.NewReader(fPath)
	if err != nil {
		os.Remove(fPath)
		t.Fatalf("error creating reader: %s", err)
	}

	if (r.Meta.Version != internal.Version) || (r.Meta.NumberOfFiles != uint64(len(fileList))) {
		t.Errorf("metadata doesnt match:\n---> Found: %+v", r.Meta)
	}

	getFiles := make([]*common.File, 0, len(fileList))
	for r.More() {
		f, err := r.Next()
		if err != nil {
			os.Remove(fPath)
			t.Fatalf("error reading next file: %s", err)
		}
		getFiles = append(getFiles, f)
	}

	if !equivalent(fileList, getFiles) {
		os.Remove(fPath)
		t.Fatalf("files are different:\n---> ORIGINAL: %+v\n---> LIST GOT: %+v", fileList, getFiles)
	}
	os.Remove(fPath)
}

func equivalent(fl1, fl2 []*common.File) bool {
	if len(fl1) != len(fl2) {
		return false
	}

	for i := range fl1 {
		f1 := fl1[i]
		f2 := fl2[i]

		if (f1.RelPath != f2.RelPath) || (f1.Hash != f2.Hash) || (f1.Size != f2.Size) {
			return false
		}
	}
	return true
}
