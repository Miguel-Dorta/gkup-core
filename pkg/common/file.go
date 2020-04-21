package common

import "sync"

type File struct {
	Path string `json:"path"`
	Hash string `json:"hash"`
	Size int64  `json:"size"`
}

type SafeFileList struct {
	List []*File
	i    int
	m    sync.Mutex
}

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
