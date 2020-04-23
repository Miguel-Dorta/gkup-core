package snapshot

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/Miguel-Dorta/gkup-core/pkg/common"
	"os"
)

type Reader struct {
	Meta Metadata
	d    *json.Decoder
	f    *os.File
}

func NewReader(path string) (*Reader, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error opening snapshot: %w", err)
	}
	defer f.Close()

	r := &Reader{
		d: json.NewDecoder(bufio.NewReaderSize(f, 128*1024)),
		f: f,
	}

	if err := r.d.Decode(&r.Meta); err != nil {
		return nil, fmt.Errorf("error decoding metadata: %w", err)
	}

	return r, nil
}

func (r *Reader) More() bool {
	return r.d.More()
}

func (r *Reader) Next() (*common.File, error) {
	var f common.File
	if err := r.d.Decode(&f); err != nil {
		return nil, fmt.Errorf("error decoding file: %w", err)
	}
	return &f, nil
}

func (r *Reader) Close() error {
	return r.f.Close()
}
