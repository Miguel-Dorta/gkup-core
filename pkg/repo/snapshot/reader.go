package snapshot

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/Miguel-Dorta/gkup-core/pkg/common"
	"os"
)

// Reader represents a snapshot reader
type Reader struct {
	Meta      Metadata
	d         *json.Decoder
	f         *os.File
	filesRead uint64
}

// NewReader will take the path of a snapshot and return a reader of it
func NewReader(path string) (*Reader, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error opening snapshot: %w", err)
	}
	defer f.Close()

	r := &Reader{
		d:         json.NewDecoder(bufio.NewReaderSize(f, 128*1024)),
		f:         f,
		filesRead: 0,
	}

	if err := r.d.Decode(&r.Meta); err != nil {
		return nil, fmt.Errorf("error decoding metadata: %w", err)
	}

	return r, nil
}

// More returns if there are more files
func (r *Reader) More() bool {
	return r.Meta.NumberOfFiles > r.filesRead
}

// Next returns the next file.
// As a consequence of the common.File serialization, the AbsPath of the file returned will be empty.
func (r *Reader) Next() (*common.File, error) {
	var f common.File
	if err := r.d.Decode(&f); err != nil {
		return nil, fmt.Errorf("error decoding file: %w", err)
	}
	r.filesRead++
	return &f, nil
}

// Close closes the reader
func (r *Reader) Close() error {
	return r.f.Close()
}
