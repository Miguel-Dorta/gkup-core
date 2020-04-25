// +build !windows

package fileUtils

import (
	"fmt"
	"golang.org/x/sys/unix"
	"os"
)

func copyFile(from, to string) error {
	fromFile, err := os.Open(from)
	if err != nil {
		return fmt.Errorf("error opening origin file %s: %w", from, err)
	}
	defer fromFile.Close()

	toFile, err := os.OpenFile(to, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return fmt.Errorf("error creating destination file %s: %w", to, err)
	}
	defer toFile.Close()

	var size int64; {
		stat, err := os.Stat(from)
		if err != nil {
			return fmt.Errorf("error getting file info from file %s: %w", from, err)
		}
		size = stat.Size()
	}

	for offset := int64(0); offset < size; {
		if _, err := unix.Sendfile(int(toFile.Fd()), int(fromFile.Fd()), &offset, int(size-offset)); err != nil {
			return fmt.Errorf("copy error: %w", err)
		}
	}

	if err := toFile.Close(); err != nil {
		return fmt.Errorf("error closing file %s: %w", to, err)
	}
	return nil
}
