package hash

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/Miguel-Dorta/gkup-core/pkg/common"
	"hash"
	"io"
	"os"
)

type Hasher struct {
	h hash.Hash
	buf []byte
}

const (
	SHA256 = "sha256"

	minimumBufferSize = 512
)

func NewHasher(hashAlgorithm string, bufSize int) *Hasher {
	var h hash.Hash
	switch hashAlgorithm {
	case SHA256:
		h = sha256.New()
	default:
		panic("invalid hash algorithm: " + hashAlgorithm)
	}

	if bufSize < minimumBufferSize {
		bufSize = minimumBufferSize
	}

	return &Hasher{
		h:   h,
		buf: make([]byte, bufSize),
	}
}

func (h *Hasher) HashFile(f *common.File) error {
	osFile, err := os.Open(f.Path)
	if err != nil {
		return fmt.Errorf("error opening file %s: %w", f.Path, err)
	}

	h.h.Reset()
	if _, err := io.CopyBuffer(h.h, osFile, h.buf); err != nil {
		return fmt.Errorf("error hashing file %s: %w", f.Path, err)
	}
	f.Hash = hex.EncodeToString(h.h.Sum(nil))
	return nil
}
