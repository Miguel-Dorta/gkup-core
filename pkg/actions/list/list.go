package list

import (
	"encoding/json"
	"fmt"
	"github.com/Miguel-Dorta/gkup-core/pkg/output"
	"github.com/Miguel-Dorta/gkup-core/pkg/repo"
	"os"
)

func Run() {
	snaps, err := repo.ListSnapshots()
	if err != nil {
		output.PrintErrorf("error listing snapshots: %s", err)
		os.Exit(1)
	}

	data, err := json.Marshal(snaps)
	if err != nil {
		output.PrintErrorf("error serializing snapshots: %s", err)
		os.Exit(1)
	}
	fmt.Println(string(data))
}
