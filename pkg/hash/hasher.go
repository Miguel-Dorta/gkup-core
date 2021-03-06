package hash

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"github.com/Miguel-Dorta/gkup-core/pkg/common"
	"golang.org/x/crypto/sha3"
	"hash"
	"io"
	"os"
)

// Hasher is a type for optimizing file hashing
type Hasher struct {
	h   hash.Hash
	buf []byte
}

const (
	// Valid hashes
	MD5      = "md5"
	SHA1     = "sha1"
	SHA256   = "sha256"
	SHA512   = "sha512"
	SHA3_256 = "sha3-256"
	SHA3_512 = "sha3-512"

	minimumBufferSize = 512
)

// NewHasher returns a new hasher with the hash algorithm and buffer size specified.
// It panics if the hash algorithm is invalid.
func NewHasher(hashAlgorithm string, bufSize int) *Hasher {
	var h hash.Hash
	switch hashAlgorithm {
	case MD5:
		h = md5.New()
	case SHA1:
		h = sha1.New()
	case SHA256:
		h = sha256.New()
	case SHA512:
		h = sha512.New()
	case SHA3_256:
		h = sha3.New256()
	case SHA3_512:
		h = sha3.New512()
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

// HashFile will take a file, hash it, and update its hash.
// It requires that the AbsPath of the file provided is not empty.
func (h *Hasher) HashFile(f *common.File) error {
	if err := f.CheckAbsPath(); err != nil {
		return err
	}

	osFile, err := os.Open(f.AbsPath)
	if err != nil {
		return fmt.Errorf("error opening file %s: %w", f.AbsPath, err)
	}
	defer osFile.Close()

	h.h.Reset()
	if _, err := io.CopyBuffer(h.h, osFile, h.buf); err != nil {
		return fmt.Errorf("error hashing file %s: %w", f.AbsPath, err)
	}
	f.Hash = hex.EncodeToString(h.h.Sum(nil))
	return nil
}
