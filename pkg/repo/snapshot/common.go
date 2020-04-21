package snapshot

type metadata struct {
	Version string `json:"version"`
}

const (
	SnapshotDirName = "snapshots"
	Extension = ".gkup"
)
