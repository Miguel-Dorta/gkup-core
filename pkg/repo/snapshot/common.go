package snapshot

// Metadata represents the metadata of a snapshot.
// It's always the first JSON found.
type Metadata struct {
	Version       string `json:"version"`
	NumberOfFiles uint64 `json:"files"`
}

const Extension = ".gkup"
