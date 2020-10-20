package main

import (
	"log"
	"os"
	"strconv"

	"github.com/davidparks11/file-renamer/pkg/config"
	"github.com/davidparks11/file-renamer/pkg/fileactions"
	"github.com/davidparks11/file-renamer/pkg/fileretriever"
	"github.com/davidparks11/file-renamer/pkg/logger"
	"github.com/davidparks11/file-renamer/pkg/logger/loggeriface"
	"github.com/davidparks11/file-renamer/pkg/schedule"
	"github.com/davidparks11/file-renamer/pkg/schedule/scheduleiface"
)

var logService loggeriface.Service
var cfg *config.Config
var scheduler scheduleiface.Scheduler

func main() {

	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal("Failed to load configuration from config.json\n", err.Error())
	}

	var logLevel int
	//command line override
	if len(os.Args) == 0 {
		logLevel, _ = strconv.Atoi(os.Args[0])
	} else {
		//get log level from config
		logLevel = logger.ParseLogLevel(cfg.LogLevel)
	}

	//Set up log service
	logService = logger.NewLogService(logLevel, cfg.LogLocation, cfg.LogToConsole)
	logService.Info("Program Start")
	defer func() {
		if logService != nil {
			logService.Info("System interrupt exiting program")
			logService.Stop()
		}
	}()

	//Set up drive service
	ft := fileretriever.NewFileRetriever(logService, cfg)

	//Set up file actions
	fileRenamer := fileactions.NewProcess(logService, ft, cfg)

	//Set up scheduler 
	scheduler = schedule.NewScheduleService(logService)

	if cfg.RunAtLaunch {
		fileRenamer.Run()
	}
	
	for _, schedule := range cfg.CronSchedules {
		scheduler.ScheduleJob(schedule, func() {
			fileRenamer.Run()
		})
	}
	logService.Info("Starting scheduler")
	scheduler.Run()
}
