package backup

import (
	"github.com/Miguel-Dorta/gkup-core/pkg/common"
	"github.com/Miguel-Dorta/gkup-core/pkg/hash"
	"github.com/Miguel-Dorta/gkup-core/pkg/output"
	"runtime"
	"sync"
)

func HashFiles(fs []*common.File, hashAlgorithm string) error {
	var commonErr error
	sfl := &common.SafeFileList{List: fs}
	wg := new(sync.WaitGroup)
	wg.Add(runtime.NumCPU())

	for i:=0; i<runtime.NumCPU(); i++ {
		go func() {
			defer wg.Done()
			h := hash.NewHasher(hashAlgorithm, 128*1024)

			for {
				if commonErr != nil {
					return
				}

				f := sfl.Next()
				if f == nil {
					break
				}

				output.NewPartialStep("Hash file " + f.AbsPath)
				if err := h.HashFile(f); err != nil {
					commonErr = err
					return
				}
			}
		}()
	}
	wg.Wait()
	return commonErr
}
