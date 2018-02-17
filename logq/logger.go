package logq

import (
	"os"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
	"github.com/blazecrystal/beyondts-go/utils"
	"github.com/blazecrystal/beyondts-go/logq/appenders"
)

const (
	unknown = "unknown"
)

var (
	levels = []string{"debug", "info", "warn", "error", "fatal"}
)

type Logger struct {
	name, level string
	appenders   []string //name of appenders
	debugMode   bool
}

func createDefaultLogger() *Logger {
	return &Logger{name: DEFAULT_LOGGER, level: levels[0], appenders: []string{DEFAULT_APPENDER}}
}

func newLogger(name, level string, debugMode bool, appenderNames ...string) *Logger {
	return &Logger{name: name, level: level, debugMode: debugMode, appenders: appenderNames}
}

func (l *Logger) update(name, level string, appenderNames ...string) {
	l.name, l.level, l.appenders = name, level, appenderNames
}

func (l *Logger) Debug(content ...interface{}) {
	l.write(0, content...)
}

func (l *Logger) Info(content ...interface{}) {
	l.write(1, content...)
}

func (l *Logger) Warn(content ...interface{}) {
	l.write(2, content...)
}

func (l *Logger) Error(err error, content ...interface{}) {
	l.ErrorMore(err, false, content...)
}

func (l *Logger) ErrorOnly(err error) {
	l.ErrorMore(err, false)
}

func (l *Logger) ErrorMore(err error, ignoreStacks bool, content ...interface{}) {
	msg := utils.Concat2(" ", content...)
	if err != nil {
		msg = utils.Concat(msg, "\n", err.Error())
	}
	if !ignoreStacks {
		msg = utils.Concat(msg, "\n")
		tmp := strings.Split(string(debug.Stack()), "\n")
		for i, sentence := range tmp {
			if i != 1 && i != 2 && i != 3 && i != 4 {
				msg = utils.Concat(msg, sentence, "\n")
			}
		}
	}
	l.write(3, msg)
}

func (l *Logger) Fatal(err error, content ...interface{}) {
	l.FatalMore(err, false, content...)
}

func (l *Logger) FatalOnly(err error) {
	l.FatalMore(err, false)
}

func (l *Logger) FatalMore(err error, ignoreStacks bool, content ...interface{}) {
	msg := utils.Concat2(" ", content...)
	if err != nil {
		msg = utils.Concat(msg, "\n", err.Error())
	}
	if !ignoreStacks {
		msg = utils.Concat(msg, "\n")
		tmp := strings.Split(string(debug.Stack()), "\n")
		for i, sentence := range tmp {
			if i != 1 && i != 2 && i != 3 && i != 4 {
				msg = utils.Concat(msg, sentence, "\n")
			}
		}
	}
	l.write(4, msg)
}

func (l *Logger) write(levelIndex int, content ...interface{}) {
	if utils.IndexInSlice(levels, strings.ToLower(l.level)) < levelIndex+1 {
		fills := make(map[string]string)
		fills[appenders.Fills_Timestamp] = time.Now().Format("2006-01-02 15:04:05.000")
		fills[appenders.Fills_Level] = strings.ToUpper(levels[levelIndex])
		fills[appenders.Fills_LoggerName] = l.name
		fills[appenders.Fills_Content] = utils.Concat2(" ", content...)
		_, file, line, ok := runtime.Caller(2) // skip Caller & current stack
		if ok {
			fills[appenders.Fills_LongFileName] = utils.ToLocalFilePath(file)
			lastPathSept := strings.LastIndex(fills[appenders.Fills_LongFileName], string(os.PathSeparator))
			fills[appenders.Fills_ShortFileName] = fills[appenders.Fills_LongFileName][lastPathSept+1:]
			fills[appenders.Fills_LineNum] = strconv.Itoa(line)
		} else {
			fills[appenders.Fills_LongFileName] = unknown
			fills[appenders.Fills_ShortFileName] = unknown
			fills[appenders.Fills_LineNum] = unknown
		}
		for _, appenderName := range l.appenders {
			go GetAppender(appenderName).WriteLog(fills)
		}
	}
}

func (l *Logger) Print(v ...interface{}) {
	l.Debug(v...)
}
