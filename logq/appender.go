package logq

import (
	"github.com/blazecrystal/beyondts-go/logq/appenders"
)

type Appender interface {
	GetType() string
	WriteLog(fills map[string]string)
	Run(flag chan string)
	Stop()
	Update(attrs map[string]string)
}

func newAppender(appenderType string, debugMode bool, attrs map[string]string) Appender {
	appenders.SetDebugMode(debugMode)
	switch appenderType {
	case appenders.Type_LocalFile:
		return appenders.NewLocalFileAppender(attrs)
	case appenders.Type_RemoteFile:
		return appenders.NewRemoteFileAppender(attrs)
	case appenders.Type_Syslog:
		return appenders.NewSyslogAppender(attrs)
	case appenders.Type_Database:
		return appenders.NewDatabaseAppender(attrs)
	default:
		return appenders.NewStdoutAppender(attrs)
	}
}

func createDefaultAppender() *appenders.StdoutAppender {
	return appenders.NewStdoutAppender(map[string]string{"name": DEFAULT_LOGGER, "layout": appenders.DEFAULT_LAYOUT})
}
