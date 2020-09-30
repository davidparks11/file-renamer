package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/davidparks11/file-renamer/pkg/logger/loggeriface"
)

var _ loggeriface.Service = &Service{}

const (
	//ERROR level for logging
	ERROR = 1
	//WARN level for logging
	WARN = 2
	//INFO level for logging
	INFO = 3
	//defaultLogPath the location written if the log path is not edited in config.json
	defaultLogPath = "file_renamer_logs"
)

//Service facilitates logging by writing to a log file at three levels
type Service struct {
	Level            int
	logToConsole     bool
	path             string
	file             *os.File
	fileCreationDate *time.Time
	fatalLogger      *log.Logger
	errorLogger      *log.Logger
	warnLogger       *log.Logger
	infoLogger       *log.Logger
}

//NewLogService serves a new log Service
func NewLogService(level int, logPath string, logToConsole bool) loggeriface.Service {

	//time for file creation
	fileCreationTime := time.Now()

	//check for no specified log directory
	if logPath == "" {
		logPath = defaultLogPath
	}
	logFileName := newLogFileName(&fileCreationTime, logPath)
	fmt.Println("Creating log: " + logFileName)
	//open log file
	logFile, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal("Failed to init logger")
	}

	service := Service{
		Level:            level,
		logToConsole:     logToConsole,
		path:             logPath,
		file:             logFile,
		fileCreationDate: &fileCreationTime,
		warnLogger:       log.New(logFile, "WARNING: ", log.Ldate|log.Ltime),
		infoLogger:       log.New(logFile, "INFO: ", log.Ldate|log.Ltime),
		errorLogger:      log.New(logFile, "ERROR: ", log.Ldate|log.Ltime),
		fatalLogger:      log.New(logFile, "FATAL: ", log.Ldate|log.Ltime),
	}
	return &service
}

//ParseLogLevel takes a string with the log level names or enum number values
//and parses it into an integer
func ParseLogLevel(logLevel string) (level int) {
	var err error
	switch strings.ToUpper(logLevel) {
	case "INFO":
		level = INFO
		break
	case "WARN":
		level = WARN
		break
	case "WARNING":
		level = WARN
		break
	case "ERROR":
		level = ERROR
		break
	default:
		level, err = strconv.Atoi(logLevel)
		if err != nil {
			level = INFO
		}
	}
	return
}

//Fatal writes a fatal message to logs and exits program
func (s *Service) Fatal(msg string) {
	if s.isNewLogDay() {
		s.refresh()
	}
	if s.logToConsole {
		fmt.Println("FATAL: ", msg)
	}

	s.fatalLogger.Fatal(msg)
}

//Error writes an error message to logs
func (s *Service) Error(msg string) {
	if !s.isLoggable(ERROR) {
		return
	}
	if s.isNewLogDay() {
		s.refresh()
	}
	if s.logToConsole {
		fmt.Println("ERROR: ", msg)
	}

	s.errorLogger.Println(msg)
}

//Warn writes a warning message to logs
func (s *Service) Warn(msg string) {
	if !s.isLoggable(WARN) {
		return
	}
	if s.isNewLogDay() {
		s.refresh()
	}
	if s.logToConsole {
		fmt.Println("WARN: ", msg)
	}

	s.warnLogger.Println(msg)
}

//Info writes an info message to logs
func (s *Service) Info(msg string) {
	if !s.isLoggable(INFO) {
		return
	}
	if s.isNewLogDay() {
		s.refresh()
	}
	if s.logToConsole {
		fmt.Println("INFO: ", msg)
	}

	s.infoLogger.Println(msg)
}

func (s *Service) isLoggable(msgLevel int) bool {
	return msgLevel <= s.Level
}

func (s *Service) isNewLogDay() bool {
	return s.fileCreationDate.Day() != time.Now().Day()
}

func newLogFileName(creationTime *time.Time, path string) string {
	year, month, day := creationTime.Date()

	if _, err := os.Stat(path); err != nil {
		if err = os.MkdirAll(path, 0666); err != nil {
			log.Fatal("Log directory (" + path + ") does not exist, failed to created log directory")
		}
	}

	name := fmt.Sprintf("%d-%d-%d.log", year, month, day)
	return filepath.Join(path, name)
}

//refresh closes the current file, opens a new file, sets the file, filecreationdate, and log outputs
func (s *Service) refresh() {
	s.file.Close()

	//time for file creation
	now := time.Now()
	logFileName := newLogFileName(&now, s.path)

	//open log file
	logFile, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal("Failed to init logger")
	}

	//reassign logger file and creation time
	s.file = logFile
	s.fileCreationDate = &now

	//reassgin logger outputs
	s.warnLogger.SetOutput(logFile)
	s.errorLogger.SetOutput(logFile)
	s.infoLogger.SetOutput(logFile)
}

//Stop performs logging service clean up
func (s *Service) Stop() {
	if s.file != nil {
		s.file.Close()
	}
}
