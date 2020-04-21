package settings

import (
	"fmt"
	"github.com/pelletier/go-toml"
	"io/ioutil"
)

type Settings struct {
	Version       string `toml:"version"`
	HashAlgorithm string `toml:"hash_algorithm"`
}

const Filename = "settings.toml"

func Load(path string) (*Settings, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading settings file %s: %w", path, err)
	}

	var s *Settings
	if err := toml.Unmarshal(data, s); err != nil {
		return nil, fmt.Errorf("error parsing settings: %w", err)
	}

	return s, nil
}

func Save(path string, s *Settings) error {
	data, err := toml.Marshal(s)
	if err != nil {
		return fmt.Errorf("error serializing settings: %w", err)
	}

	if err := ioutil.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("error writing settings file in %s: %w", path, err)
	}
	return nil
}
