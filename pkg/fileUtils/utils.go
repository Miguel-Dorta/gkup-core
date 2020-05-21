package fileUtils

import "os"

// CopyFile copies a file from "from" to "to"
func CopyFile(from, to string) error {
	return copyFile(from, to)
}

// FileExists returns if a specific path exists
func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
