package common

import "sync"

// SafeFileList represents a list of File safe for concurrent use.
type SafeFileList struct {
	List []*File
	i    int
	m    sync.Mutex
}

// Next will return the next value, or nil if no more values are found.
func (sfl *SafeFileList) Next() *File {
	sfl.m.Lock()
	defer sfl.m.Unlock()
	if sfl.i < len(sfl.List) {
		f := sfl.List[sfl.i]
		sfl.i++
		return f
	}
	return nil
}
