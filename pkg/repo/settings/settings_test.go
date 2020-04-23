package settings

import (
	"github.com/Miguel-Dorta/gkup-core/pkg/hash"
	"os"
	"path/filepath"
	"testing"
)

func Test(t *testing.T) {
	s1 := &Settings{
		Version:       "v9.8.7-alpha+test",
		HashAlgorithm: hash.MD5,
	}
	if err := Save(os.TempDir(), s1); err != nil {
		t.Fatalf("error saving settings: %s", err)
	}

	s2, err := Load(os.TempDir())
	if err != nil {
		os.Remove(filepath.Join(os.TempDir(), filename))
		t.Fatalf("error loading settings: %s", err)
	}
	if (s1.Version != s2.Version) || (s1.HashAlgorithm != s2.HashAlgorithm) {
		os.Remove(filepath.Join(os.TempDir(), filename))
		t.Fatalf("data doesn't match:\n--> SAVED:  %+v\n--> LOADED: %+v", s1, s2)
	}
	os.Remove(filepath.Join(os.TempDir(), filename))
}
