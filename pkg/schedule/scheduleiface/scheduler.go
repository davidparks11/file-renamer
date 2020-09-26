package scheduleiface

//Scheduler contains methods to facilitate 
type Scheduler interface {
	ScheduleJob(schedule string, process func()) error
	Run()
}