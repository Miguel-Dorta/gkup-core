package settings

import (
	"fmt"
	"github.com/pelletier/go-toml"
	"io/ioutil"
	"path/filepath"
)

// Settings represents the settings of a repo
type Settings struct {
	Version       string `toml:"version"`
	HashAlgorithm string `toml:"hash_algorithm"`
}

const filename = "settings.toml"

// Load will take the path of a repo and return a Settings object that contains the repo's settings
func Load(path string) (*Settings, error) {
	data, err := ioutil.ReadFile(filepath.Join(path, filename))
	if err != nil {
		return nil, fmt.Errorf("error reading settings file %s: %w", filepath.Join(path, filename), err)
	}

	var s Settings
	if err := toml.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("error parsing settings: %w", err)
	}

	return &s, nil
}

// Save will take the path of a repo and a Settings object, and write it
func Save(path string, s *Settings) error {
	data, err := toml.Marshal(s)
	if err != nil {
		return fmt.Errorf("error serializing settings: %w", err)
	}

	if err := ioutil.WriteFile(filepath.Join(path, filename), data, 0644); err != nil {
		return fmt.Errorf("error writing settings file in %s: %w", filepath.Join(path, filename), err)
	}
	return nil
}
