package main

import (
	"flag"
	"fmt"
	"github.com/Miguel-Dorta/gkup-core/internal"
	"github.com/Miguel-Dorta/gkup-core/pkg/actions/backup"
	"github.com/Miguel-Dorta/gkup-core/pkg/actions/create"
	"github.com/Miguel-Dorta/gkup-core/pkg/actions/list"
	"github.com/Miguel-Dorta/gkup-core/pkg/actions/restore"
	"github.com/Miguel-Dorta/gkup-core/pkg/output"
	"github.com/Miguel-Dorta/gkup-core/pkg/repo"
	"os"
	"time"
)

var (
	action, repoPath, snapName, hashAlgorithm, destination string
	snapTime                                               int64
)

func init() {
	var verbose bool
	flag.StringVar(&action, "action", "", "[ALL] Action to do: BACKUP, CREATE, LIST, RESTORE, VERSION")
	flag.StringVar(&repoPath, "repo", "", "[ALL] Repository path")
	flag.BoolVar(&verbose, "v", false, "[ALL] Verbose")
	flag.StringVar(&snapName, "snapshot-name", "", "[BACKUP,RESTORE] Snapshot name")
	flag.StringVar(&hashAlgorithm, "hash-algorithm", "sha256", "[CREATE] Hash algorithm: md5, sha1, sha256, sha512, sha3-256, sha3-512")
	flag.Int64Var(&snapTime, "snapshot-time", 0, "[RESTORE] Snapshot timestamp")
	flag.StringVar(&destination, "restore-destination", "", "[RESTORE] Destination to restore snapshot")
	flag.Parse()

	output.Verbose = verbose
}

func main() {
	switch action {
	case "BACKUP":
		repoInit()
		backup.Run(snapName, flag.Args()...)
	case "CREATE":
		create.Run(repoPath, hashAlgorithm)
	case "LIST":
		repoInit()
		list.Run()
	case "RESTORE":
		repoInit()
		restore.Run(snapName, time.Unix(snapTime, 0), destination)
	case "VERSION":
		fmt.Println(internal.Version)
	default:
		output.PrintErrorf("invalid action: %s", action)
	}
}

func repoInit() {
	if err := repo.Init(repoPath); err != nil {
		output.PrintErrorf("error initializing repository: %s", err)
		os.Exit(1)
	}
}
