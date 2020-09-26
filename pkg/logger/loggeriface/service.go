package loggeriface

//Service contains methods to log at different levels
type Service interface { 
	Info(string)
	Error(string)
	Warn(string)
	Stop()
}