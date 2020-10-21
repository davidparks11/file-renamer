package fileactions

import (
	"fmt"
	"strings"
	"time"

	"github.com/davidparks11/file-renamer/pkg/config"
	"github.com/davidparks11/file-renamer/pkg/fileactions/fileactionsiface"
	"github.com/davidparks11/file-renamer/pkg/fileretriever/fileretrieveriface"
	"github.com/davidparks11/file-renamer/pkg/logger/loggeriface"
)

var _ fileactionsiface.Process = &Renamer{}

//Renamer is a process used to rename a set of files
type Renamer struct {
	logger loggeriface.Service
	fileRetriever fileretrieveriface.FileRetriever
	name string
	config *config.Config
	processedFiles map[string]bool
}

//NewProcess returns a Renamer that uniquely names each file based 
//on its configured persistent words and creation date
func NewProcess(logger loggeriface.Service, fileRetriever fileretrieveriface.FileRetriever, config *config.Config) fileactionsiface.Process {
	return &Renamer{
		logger: logger,
		fileRetriever: fileRetriever,
		name: "File-Renamer",
		config: config,
		processedFiles: nil,
	}
}

//allows control of time for testing
var now = func() time.Time {
	return time.Now()
}

//Run is called on a repeated schedule by a scheduler
func (r *Renamer) Run() error {
	r.logger.Info(fmt.Sprintf("~~~~ %s started ~~~~", r.name))
	files, err := r.fileRetriever.GetFileInfo()
	if err != nil {
		return err
	}

	//get all processed files. Runs each time in case of deletions
	r.processedFiles = r.fileRetriever.GetProcessedFiles()

	for _, file := range files {
		file.Name, err = r.generateNewName(file.Name, file.CreatedDate)
		if err != nil {
			//skip file if error is encountered
			r.logger.Error(fmt.Sprintf("Error generating new file name %s - %s", file.ID, err.Error()))
			continue
		}
		err = r.fileRetriever.UpdateFile(file)
		if err != nil {
			r.logger.Error(fmt.Sprintf("Error updating file %s:%s - %s", file.Name, file.ID, err.Error()))
		}
		r.logger.Info("Updated file name to " + file.Name)
		r.processedFiles[file.Name] = true
	}
	r.logger.Info(fmt.Sprintf("~~~~ %s ended ~~~~", r.name))
	return nil
}

//FIRSTNAME_GAMEMODE_YYYY_MMDD_#.

func (r *Renamer) generateNewName(name string, createdDate string) (string, error) {
	var newName string
	var suffix string

	//extrac suffix
	suffixIndex := strings.LastIndex(name, ".")
	if suffixIndex == -1 {
		suffix = ""
	} else {
		suffix = name[suffixIndex:]
	}
		
	//add any persistent words to the the new name
	for _, word := range r.config.PersistentWords {
		if strings.Contains(strings.ToLower(name), strings.ToLower(word)) {
			newName += word + r.config.NameDelimiter
		}
	}

	parsedDate, err := r.parseRFC3339(createdDate)
	if err != nil {
		return "", err
	}
	newName += parsedDate + r.config.NameDelimiter

	var dupCheck string
	for dupFileCount := 0; true; dupFileCount++ {
		dupCheck = fmt.Sprintf("%s%d%s", newName, dupFileCount, suffix)
		if r.processedFiles[dupCheck] == false {
			break
		}
	}
	return dupCheck, nil
}

//a time format of YYYY_MMDD
const timeFormat = "2006_0102"

func (r *Renamer) parseRFC3339(date string) (string, error) {
	parsed, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return "", err
	}
	return parsed.Format(timeFormat), nil
}