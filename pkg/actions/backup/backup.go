package backup

import (
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
		output.NewPartialStep("Add file " + f.Path)
		if err := repo.AddFile(f); err != nil {
			output.PrintErrorf("error adding file %s to repository: %w", f.Path, err)
			os.Exit(1)
		}
	}

	output.NewGlobalStep("Save snapshot", 0)
	snap, err := repo.NewSnapshot(snapshotName, startTime)
	if err != nil {
		output.PrintError(err)
		os.Exit(1)
	}
	for _, f := range fList {
		if err := snap.Write(f); err != nil {
			output.PrintErrorf("error adding file %s to snapshot: %w", f.Path, err)
			os.Exit(1)
		}
	}
}
