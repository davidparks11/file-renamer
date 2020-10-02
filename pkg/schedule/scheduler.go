package schedule

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/davidparks11/file-renamer/pkg/logger/loggeriface"
	"github.com/davidparks11/file-renamer/pkg/schedule/scheduleiface"
	"github.com/robfig/cron/v3"
)


var _ scheduleiface.Scheduler = &Scheduler{}

type Scheduler struct {
	cron *cron.Cron
	sigChannel chan os.Signal
	logger loggeriface.Service
}

func NewScheduleService(logger loggeriface.Service) scheduleiface.Scheduler {
	cronLogger := cron.DefaultLogger
	return &Scheduler{
		cron: cron.New(cron.WithLogger(cronLogger), cron.WithChain(cron.SkipIfStillRunning(cronLogger))),
		sigChannel: make(chan os.Signal),
		logger: logger,
	}
}

//ScheduleJob schedules a function to run on a cron job schedule
func (s *Scheduler) ScheduleJob(schedule string, process func()) error {
	_, err := s.cron.AddFunc(schedule, process)
	return err
}

func (s *Scheduler) InterrupetChannel() chan os.Signal {
	s.initInterreupt()
	return s.sigChannel
}

func (s *Scheduler) initInterreupt() {
	if s.sigChannel != nil {
		return
	}
	s.sigChannel = make (chan os.Signal, 1)
}

func (s *Scheduler) Run() {
	s.cron.Start()
	s.initInterreupt()
	signal.Notify(s.sigChannel, syscall.SIGINT, syscall.SIGTERM)
	<-s.sigChannel
	s.cron.Stop()
	if s.cron == nil {
		s.logger.Info("Cron is nil after stop")
	}
}
