package backup

import (
	"github.com/Miguel-Dorta/gkup-core/pkg/input"
	"github.com/Miguel-Dorta/gkup-core/pkg/output"
	"github.com/Miguel-Dorta/gkup-core/pkg/repo"
	"os"
	"time"
)

func Run(snapshotName string, paths ...string) {
	startTime := time.Now().UTC()
	output.Setup(4)

	output.NewGlobalStep("List files", 0)
	fList, err := list(paths...)
	if err != nil {
		output.PrintError(err)
		os.Exit(1)
	}
	filesToProcess := uint64(len(fList))

	output.NewGlobalStep("Hash files", filesToProcess)
	if err := HashFiles(fList, repo.Sett.HashAlgorithm); err != nil {
		output.PrintError(err)
		os.Exit(1)
	}

	output.NewGlobalStep("Add files to repository", filesToProcess)
	for _, f := range fList {
		// Control channels
		select {
		case <-input.Pause:
			select {
			case <-input.Stop:
				output.PrintError(output.ErrProcessStopped)
				os.Exit(1)
			case <-input.Resume:
			}
		case <-input.Stop:
			output.PrintError(output.ErrProcessStopped)
			os.Exit(1)
		default:
		}

		output.NewPartialStep("Add file " + f.AbsPath)
		if err := repo.AddFile(f); err != nil {
			output.PrintErrorf("error adding file %s to repository: %s", f.AbsPath, err)
			os.Exit(1)
		}
	}

	output.NewGlobalStep("Save snapshot", 0)
	snap, err := repo.NewSnapshot(snapshotName, startTime, filesToProcess)
	if err != nil {
		output.PrintError(err)
		os.Exit(1)
	}
	for _, f := range fList {
		if err := snap.Write(f); err != nil {
			output.PrintErrorf("error adding file %s to snapshot: %s", f.AbsPath, err)
			os.Exit(1)
		}
	}
	if err := snap.Close(); err != nil {
		output.PrintErrorf("error closing snapshot file: %s", err)
		os.Exit(1)
	}
}
