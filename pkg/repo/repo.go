package repo

import (
	"errors"
	"fmt"
	"github.com/Miguel-Dorta/gkup-core/internal"
	"github.com/Miguel-Dorta/gkup-core/pkg/repo/settings"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	Sett *settings.Settings
	path string
)

func Init(repoPath string) error {
	path = repoPath

	s, err := settings.Load(filepath.Join(path))
	if err != nil {
		return fmt.Errorf("error loading settings: %w", err)
	}
	Sett = s
	return nil
}

func Create(repoPath string, s *settings.Settings) error {
	path = repoPath
	Sett = s
	Sett.Version = internal.Version

	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("error creating repository directory %s: %w", path, err)
	}

	if fs, err := ioutil.ReadDir(path); err != nil {
		return fmt.Errorf("error reading repository directory %s: %w", path, err)
	} else if len(fs) != 0 {
		return errors.New("repository directory must be empty")
	}

	if err := settings.Save(path, Sett); err != nil {
		return fmt.Errorf("error saving settings: %w", err)
	}
}
