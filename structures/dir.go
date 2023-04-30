package structures

// Define a struct to represent a directory.
type Directory struct {
	Path    string       `json:"path"`
	Subdirs []*Directory `json:"subdirs"`
	Files   []string     `json:"files"`
}

type Stats struct {
	Files       int
	Directories int
	TotalSize   int64
}
