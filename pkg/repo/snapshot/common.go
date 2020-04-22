package snapshot

type Metadata struct {
	Version       string `json:"version"`
	NumberOfFiles uint64 `json:"files"`
}

const (
	SnapshotDirName = "snapshots"
	Extension       = ".gkup"
)
