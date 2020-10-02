package fileactions

import (
	"github.com/davidparks11/file-renamer/pkg/fileactions/fileactionsiface"
	"github.com/davidparks11/file-renamer/pkg/fileretriever/fileretrieveriface"
	"github.com/davidparks11/file-renamer/pkg/logger/loggeriface"
)

var _ fileactionsiface.Process = &Process{}

type Process struct {
	logger loggeriface.Service
	fileRetriever fileretrieveriface.FileRetriever
	name string
}

func NewProcess(logger loggeriface.Service, fileRetriever fileretrieveriface.FileRetriever) fileactionsiface.Process {
	return &Process{
		logger: logger,
		fileRetriever: fileRetriever,
		name: "File-Renamer",
	}
}

func (p *Process) Run() {
	p.fileRetriever.GetFileInfo()
}
