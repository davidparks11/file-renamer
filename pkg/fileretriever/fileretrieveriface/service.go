package fileretrieveriface

import "os"

type FileRetriever interface {
	GetFileInfo() (*[]os.FileInfo, error)
	UpdateFiles() error
}
