package restore

import (
	"fmt"
	"github.com/Miguel-Dorta/gkup-core/pkg/output"
	"github.com/Miguel-Dorta/gkup-core/pkg/repo"
	"io/ioutil"
	"os"
	"time"
)

func Run(snapshotName string, snapshotTime time.Time, destination string) {
	isEmpty, err := isDirEmpty(destination)
	if err != nil {
		output.PrintErrorf("error checking if dir %s is empty: %s", destination, err)
		os.Exit(1)
	}
	if !isEmpty {
		output.PrintError("destination must be empty")
		os.Exit(1)
	}

	snap, err := repo.OpenSnapshot(snapshotName, snapshotTime)
	if err != nil {
		output.PrintError(err)
		os.Exit(1)
	}
	defer snap.Close()

	output.Setup(1)
	output.NewGlobalStep("Restore files", snap.Meta.NumberOfFiles)
	for snap.More() {
		f, err := snap.Next()
		if err != nil {
			output.PrintError("error reading next file")
			os.Exit(1)
		}

		output.NewPartialStep(fmt.Sprintf("Restore %s", f.RelPath))
		if err := repo.RestoreFile(f, destination); err != nil {
			output.PrintErrorf("error restoring file: %s", err)
			os.Exit(1)
		}
	}
}

func isDirEmpty(path string) (bool, error) {
	fs, err := ioutil.ReadDir(path)
	if err != nil {
		return false, err
	}
	return len(fs) == 0, nil
}
