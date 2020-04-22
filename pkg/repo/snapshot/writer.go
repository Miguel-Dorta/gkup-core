package snapshot

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/Miguel-Dorta/gkup-core/internal"
	"github.com/Miguel-Dorta/gkup-core/pkg/common"
	"os"
)

type Writer struct {
	e *json.Encoder
	b *bufio.Writer
	f *os.File
}

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
	}

	if err := w.e.Encode(Metadata{
		Version:       internal.Version,
		NumberOfFiles: numberOfFiles,
	}); err != nil {
		return nil, fmt.Errorf("error writing metadata: %w", err)
	}
	return w, nil
}

func (w *Writer) Write(f *common.File) error {
	return w.e.Encode(f)
}

func (w *Writer) Close() error {
	if err := w.b.Flush(); err != nil {
		return fmt.Errorf("error flushing buffer: %w", err)
	}
	if err := w.f.Close(); err != nil {
		return fmt.Errorf("error closing file: %w", err)
	}
	return nil
}
