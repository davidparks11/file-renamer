package fileactions

import (
	"github.com/davidparks11/file-renamer/pkg/fileactions/fileactionsiface"
	"github.com/davidparks11/file-renamer/pkg/gdrive/gdriveiface"
	"github.com/davidparks11/file-renamer/pkg/logger/loggeriface"
)

var _ fileactionsiface.Process = &Process{}

type Process struct {
	logger loggeriface.Service
	drive gdriveiface.Drive
	name string
}

func NewProcess(logger loggeriface.Service, drive gdriveiface.Drive) fileactionsiface.Process {
	return &Process{
		logger: logger,
		drive: drive,
		name: "File-Renamer",
	}
}

func (p *Process) Run() {
	
}
