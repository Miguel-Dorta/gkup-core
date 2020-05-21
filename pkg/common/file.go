package common

import "errors"

// File is the gkup internal representation of a file. It is NOT guaranteed that all of any of its fields have data, so
// if you don't trust its source you should check for its zero value.
//
// Also, if you implement this type, you should return the specified zero value in each file you left empty.
//
// The zero values of each field are:
//   - AbsPath: ""
//   - RelPath: ""
//   - Hash:    ""
//   - Size:    Any value below 0 (0 itself excluded, it's a valid size)
type File struct {
	// AbsPath represents the actual path in the filesystem
	AbsPath string `json:"-"`

	// RelPath represents the relative path of the file in the backup structure
	RelPath string `json:"path"`

	// Hash is an hexadecimal representation of the hash of the file
	Hash    string `json:"hash"`

	// Size represent the size of the file in bytes
	Size    int64  `json:"size"`
}

var (
	// Errors
	ErrEmptyAbsPath = errors.New("empty file path")
	ErrEmptyRelPath = errors.New("empty file relative path")
	ErrEmptyHash = errors.New("empty file hash")
	ErrEmptySize = errors.New("empty file size")
)

// CheckAbsPath returns ErrEmptyAbsPath if the AbsPath property of the object provided has a zero-value
func (f *File) CheckAbsPath() error {
	if f.AbsPath == "" {
		return ErrEmptyAbsPath
	}
	return nil
}

// CheckRelPath returns ErrEmptyRelPath if the RelPath property of the object provided has a zero-value
func (f *File) CheckRelPath() error {
	if f.RelPath == "" {
		return ErrEmptyRelPath
	}
	return nil
}

// CheckHash returns ErrEmptyHash if the Hash property of the object provided has a zero-value
func (f *File) CheckHash() error {
	if f.Hash == "" {
		return ErrEmptyHash
	}
	return nil
}

// CheckSize returns ErrEmptySize if the Size property of the object provided has a zero-value
func (f *File) CheckSize() error {
	if f.Size < 0 {
		return ErrEmptySize
	}
	return nil
}
