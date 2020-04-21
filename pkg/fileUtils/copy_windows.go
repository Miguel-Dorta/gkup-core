package fileUtils

import (
	"fmt"
	"io"
	"os"
)

var copyBuf = make([]byte, 128*1024)

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

	if _, err := io.CopyBuffer(toFile, fromFile, copyBuf); err != nil {
		return fmt.Errorf("copy error: %w", err)
	}

	if err := toFile.Close(); err != nil {
		return fmt.Errorf("error closing file %s: %w", to, err)
	}
	return nil
}
