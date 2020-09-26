package fileactions

import (
	"github.com/davidparks11/file-renamer/pkg/fileactions/fileactionsiface"
	"github.com/davidparks11/file-renamer/pkg/logger/loggeriface"
)

var _ fileactionsiface.Process = &Process{}

type Process struct {
	logger loggeriface.Service
}

func NewProcess(logger loggeriface.Service) fileactionsiface.Process {
	return &Process{
		logger: logger,
	}
}

func (p *Process) Run() {
	panic("Implement me")
}
