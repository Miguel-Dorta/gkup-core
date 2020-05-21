package snapshot

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Miguel-Dorta/gkup-core/internal"
	"github.com/Miguel-Dorta/gkup-core/pkg/common"
	"os"
)

// Writer represents the writer of a snapshot
type Writer struct {
	e *json.Encoder
	b *bufio.Writer
	f *os.File
	filesToWrite, filesWritten uint64
}

var ErrNumberOfFilesDoesNotMatch = errors.New("numbers of files written doesn't match with the numberOfFiles provided when initializing writer")

// NewWriter takes the path of a new snapshot and the number of files that it will contain, and return a new writer.
func NewWriter(path string, numberOfFiles uint64) (*Writer, error) {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return nil, fmt.Errorf("error opening snapshot: %w", err)
	}

	b := bufio.NewWriterSize(f, 128*1024)
	w := &Writer{
		e: json.NewEncoder(b),
		b: b,
		f: f,
		filesToWrite: numberOfFiles,
		filesWritten: 0,
	}

	if err := w.e.Encode(Metadata{
		Version:       internal.Version,
		NumberOfFiles: numberOfFiles,
	}); err != nil {
		return nil, fmt.Errorf("error writing metadata: %w", err)
	}
	return w, nil
}

// Write writes a new file
func (w *Writer) Write(f *common.File) error {
	err := w.e.Encode(f)
	if err != nil {
		return err
	}
	w.filesWritten++
	return nil
}

// Close flushes the buffers and closes the writer.
// If the numberOfFiles provided when initializing doesn't match the number of files wrote,
// it returns ErrNumberOfFilesDoesNotMatch AFTER closing the object.
func (w *Writer) Close() error {
	if err := w.b.Flush(); err != nil {
		return fmt.Errorf("error flushing buffer: %w", err)
	}
	if err := w.f.Close(); err != nil {
		return fmt.Errorf("error closing file: %w", err)
	}
	if w.filesWritten != w.filesToWrite {
		return ErrNumberOfFilesDoesNotMatch
	}
	return nil
}
