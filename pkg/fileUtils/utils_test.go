package fileUtils_test

import (
	"bytes"
	"github.com/Miguel-Dorta/gkup-core/pkg/fileUtils"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
	"time"
)

func TestFileExists(t *testing.T) {
	exists, err := fileUtils.FileExists("/")
	if err != nil {
		t.Fatalf("error checking existence of /: %s", err)
	}
	if !exists {
		t.Error("it report that / doesn't exists")
	}

	exists, err = fileUtils.FileExists("/jdslvjkbjvbiukafudsghcxkjzhdhsfasffadcxz")
	if err != nil {
		t.Fatalf("error checking existence of non-existent file: %s", err)
	}
	if exists {
		t.Error("it report that a non-existent file exists")
	}

	exists, err = fileUtils.FileExists("/root/hi")
	if err == nil {
		t.Fatalf("it can get the existence of a file that it shouldn't access. Are you executing this with root permissions?")
	}
}

func TestCopyFile(t *testing.T) {
	f1, f2 := "/tmp/gkup-core_pkg_fileUtils_TestCopyFile1", "/tmp/gkup-core_pkg_fileUtils_TestCopyFile2"
	f1Data := randomData(8*1024*1024)

	if err := ioutil.WriteFile(f1, f1Data, 0644); err != nil {
		t.Fatalf("cannot write tmp file %s: %s", f1, err)
	}

	if err := fileUtils.CopyFile(f1, f2); err != nil {
		os.Remove(f1)
		t.Fatalf("error copying file from %s to %s: %s", f1, f2, err)
	}

	f2Data, err := ioutil.ReadFile(f2)
	if err != nil {
		os.Remove(f1)
		os.Remove(f2)
		t.Fatalf("error reading file %s: %s", f2, err)
	}

	if !bytes.Equal(f1Data, f2Data) {
		os.Remove(f1)
		os.Remove(f2)
		t.Fatal("files don't match")
	}
	os.Remove(f1)
	os.Remove(f2)
}

func randomData(size int) []byte {
	b := make([]byte, size)
	rand.Seed(time.Now().UnixNano())
	rand.Read(b)
	return b
}
