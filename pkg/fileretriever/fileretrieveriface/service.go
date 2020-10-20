package fileretrieveriface

type RenameInfo struct {
	ID 	 string
	Name string
	CreatedDate string
}

type FileRetriever interface {
	GetFileInfo() ([]*RenameInfo, error)
	IsUniqueName(name string) bool
	GetProcessedFiles(date string) map[string]bool
	UpdateFile(*RenameInfo) error
}
