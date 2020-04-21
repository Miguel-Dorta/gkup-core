package backup

import (
	"fmt"
	"github.com/Miguel-Dorta/gkup-core/pkg/common"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	DefaultTotalListSize = 1000
	DefaultListSize = 10
)

func list(paths ...string) ([]*common.File, error) {
	fileList := make([]*common.File, 0, DefaultTotalListSize)

	for _, path := range paths {
		pathStat, err := os.Stat(path)
		if err != nil {
			return nil, fmt.Errorf("error getting information of file %s: %w", path, err)
		}

		if pathStat.IsDir() {
			pathList, err := listDir(path)
			if err != nil {
				return nil, err
			}
			fileList = append(fileList, pathList...)
		} else if pathStat.Mode().IsRegular() {
			fileList = append(fileList, &common.File{
				Path: path,
				Size: pathStat.Size(),
			})
		} else {
			return nil, fmt.Errorf("unsupported type in file %s", path)
		}
	}

	return fileList, nil
}

func listDir(path string) ([]*common.File, error) {
	fileList := make([]*common.File, 0, DefaultListSize)

	fs, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("error reading dir %s: %w", path, err)
	}

	for _, f := range fs {
		fPath := filepath.Join(path, f.Name())
		if f.IsDir() {
			childList, err := listDir(fPath)
			if err != nil {
				return nil, err
			}
			fileList = append(fileList, childList...)
		} else if f.Mode().IsRegular() {
			fileList = append(fileList, &common.File{
				Path: fPath,
				Size: f.Size(),
			})
		} else {
			return nil, fmt.Errorf("unsupported type in file %s", fPath)
		}
	}

	return fileList, nil
}
