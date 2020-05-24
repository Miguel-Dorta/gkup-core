package backup

import (
	"github.com/Miguel-Dorta/gkup-core/pkg/common"
	"github.com/Miguel-Dorta/gkup-core/pkg/hash"
	"github.com/Miguel-Dorta/gkup-core/pkg/input"
	"github.com/Miguel-Dorta/gkup-core/pkg/output"
	"runtime"
	"sync"
)

// HashFiles hash a list of files concurrently using the hash algorithm provided
func HashFiles(fs []*common.File, hashAlgorithm string) error {
	var commonErr error
	sfl := &common.SafeFileList{List: fs}

	// End workers channel handler
	wg := new(sync.WaitGroup)
	wg.Add(runtime.NumCPU())
	done := make(chan bool)
	go func() {
		wg.Wait()
		done<-true
	}()

	// Make control channels
	pauseChans, pause := createMultipleChans(runtime.NumCPU())
	resumeChans, resume := createMultipleChans(runtime.NumCPU())
	stopChans, stop := createMultipleChans(runtime.NumCPU())

	for i:=0; i<runtime.NumCPU(); i++ {
		go func(pauseChan, resumeChan, stopChan chan bool) {
			defer wg.Done()
			h := hash.NewHasher(hashAlgorithm, 128*1024)

			for {
				select {
				case <-pauseChan:
					select {
					case <-stopChan:
						return
					case <-resumeChan:
					}
				case <-stopChan:
					return
				default:
				}

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
		}(pauseChans[i], resumeChans[i], stopChans[i])
	}

	for {
		select {
		case <-input.Pause:
			pause()
		case <-input.Resume:
			resume()
		case <-input.Stop:
			stop()
			<-done
			return output.ErrProcessStopped
		case <-done:
			return commonErr
		}
	}
}

func createMultipleChans(i int) ([]chan bool, func()) {
	chans := make([]chan bool, i)
	for i := range chans {
		chans[i] = make(chan bool, 1)
	}
	return chans, func() {
		for i := range chans {
			chans[i] <- true
		}
	}
}
