package hash_test

import (
	"fmt"
	"github.com/Miguel-Dorta/gkup-core/pkg/common"
	"github.com/Miguel-Dorta/gkup-core/pkg/hash"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestHasher_HashFile(t *testing.T) {
	tmpFile := filepath.Join(os.TempDir(), "gkup-core_pkg_hash_TestHasher-HashFile_" + strconv.FormatInt(time.Now().UnixNano(), 10))
	if err := createRandomFile(tmpFile); err != nil {
		t.Fatalf("cannot create tmpFile %s: %s", tmpFile, err)
	}

	testHasher(hash.MD5, tmpFile, t)
	testHasher(hash.SHA1, tmpFile, t)
	testHasher(hash.SHA256, tmpFile, t)
	testHasher(hash.SHA512, tmpFile, t)
	testHasher(hash.SHA3_256, tmpFile, t)
	testHasher(hash.SHA3_512, tmpFile, t)

	os.Remove(tmpFile)
}

func testHasher(hashAlgorithm, path string, t *testing.T) {
	f := &common.File{AbsPath: path}
	h := hash.NewHasher(hashAlgorithm, 128*1024)

	if err := h.HashFile(f); err != nil {
		t.Errorf("error hashing file %s with hashing algorithm %s: %s", path, hashAlgorithm, err)
		return
	}

	opensslHash, err := getOpensslHash(hashAlgorithm, path)
	if err != nil {
		t.Error(err)
		return
	}

	if f.Hash != opensslHash {
		t.Error("hashes don't match")
	}
}

func getOpensslHash(hashAlgorithm, path string) (string, error) {
	stdout, stderr := new(strings.Builder), new(strings.Builder)
	cmd := exec.Command("openssl", "dgst", "-"+hashAlgorithm, "-r", path)
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	if err := cmd.Run(); err != nil {
		return stderr.String(), fmt.Errorf("error executing openssl command with hash algorithm %s in file %s: %s", hashAlgorithm, path, err)
	}
	if stderr.Len() != 0 {
		return "", fmt.Errorf("openssl printed this on stderr: %s", stderr.String())
	}

	result := new(strings.Builder)
	for _, r := range stdout.String() {
		if !isHex(r) {
			break
		}
		result.WriteRune(r)
	}
	return result.String(), nil
}

func createRandomFile(path string) error {
	data := make([]byte, 4*1024*1024)
	rand.Seed(time.Now().UnixNano())
	rand.Read(data)
	return ioutil.WriteFile(path, data, 0644)
}

func isHex(r rune) bool {
	return (r >= '0' && r <='9') || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F')
}
