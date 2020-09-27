package gdriveiface

import "os"

type Drive interface {
	GetFileInfo() (*[]os.FileInfo, error)
	UpdateFiles() error
}
