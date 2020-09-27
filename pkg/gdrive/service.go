package gdrive

import (
	"os"

	"github.com/davidparks11/file-renamer/pkg/gdrive/gdriveiface"
	"github.com/davidparks11/file-renamer/pkg/logger/loggeriface"
)


var _ gdriveiface.Drive = &Drive{}

type Drive struct {
	logger loggeriface.Service
}

func NewDriveService(logger loggeriface.Service) gdriveiface.Drive {
	return &Drive{
		logger: logger,
	}
}

func (d *Drive) GetFileInfo() (*[]os.FileInfo, error) {

	return nil, nil
}

func (d *Drive) UpdateFiles() error {
	return nil
}
