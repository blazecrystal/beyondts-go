package logq

import (
    "github.com/blazecrystal/beyondts-go/logq/appenders"
    "strings"
    "fmt"
    "github.com/blazecrystal/beyondts-go/utils"
    "runtime/debug"
)

var (
    LEVELS = []string{"debug", "info", "warn", "error", "fatal"}
)

const (
    DEFAULT_LOGGER = "default"
    DEFAULT_LEVEL  = "debug"
)

type Logger struct {
    name, level string
    as          []string
}

func NewLogger(name, level string) *Logger {
    return &Logger{name: name, level: level}
}

func CreateDefaultLogger() *Logger {
    return &Logger{name: DEFAULT_LOGGER, level: DEFAULT_LEVEL, as: []string{appenders.DEFAULT_APPENDER}}
}

func (l1 *Logger) NeedUpdate(l2 *Logger) bool {
    return l1.getId() != l2.getId()
}

func (l1 *Logger) UpdateTo(l2 *Logger) {
    l1.level = l2.level
    l1.as = l2.as
}

func (l *Logger) getId() string {
    id := l.name + "#" + l.level
    for _, a := range l.as {
        id += a + "#"
    }
    return id
}

func (l *Logger) Debug(contents ...interface{}) {
    // content template, arg1, arg2, ..., argn, eg. Debug("abc{}def{}ghi", 11, 22)
    if len(contents) > 0 && getLevelIndex(l.level) <= 0 {
        for _, a := range l.as {
            appender := getAppender(a)
            if appender == nil {
                break
            }
            appender.Write(l.name, "debug", buildContent(contents...))
        }
    }
}

func (l *Logger) Info(contents ...interface{}) {
    if len(contents) > 0 && getLevelIndex(l.level) <= 1 {
        for _, a := range l.as {
            appender := getAppender(a)
            if appender == nil {
                break
            }
            appender.Write(l.name, "info", buildContent(contents...))
        }
    }
}

func (l *Logger) Warn(contents ...interface{}) {
    if len(contents) > 0 && getLevelIndex(l.level) <= 2 {
        for _, a := range l.as {
            appender := getAppender(a)
            if appender == nil {
                break
            }
            appender.Write(l.name, "warn", buildContent(contents...))
        }
    }
}

func (l *Logger) Error(err error, contents ...interface{}) {
    if len(contents) > 0 && getLevelIndex(l.level) <= 3 {
        for _, a := range l.as {
            appender := getAppender(a)
            if appender == nil {
                break
            }
            appender.Write(l.name, "error", buildContentWithError(err, contents...))
        }
    }
}

func (l *Logger) Fatal(err error, contents ...interface{}) {
    if len(contents) > 0 && getLevelIndex(l.level) <= 4 {
        for _, a := range l.as {
            appender := getAppender(a)
            if appender == nil {
                break
            }
            appender.Write(l.name, "fatal", buildContentWithError(err, contents...))
        }
    }
}

func getLevelIndex(level string) int {
    for index, tmp := range LEVELS {
        if strings.EqualFold(tmp, level) {
            return index
        }
    }
    return 99
}

func buildContent(v ...interface{}) string {
    switch v[0].(type) {
    case string:
        if strings.Contains(v[0].(string), "{}") {
            x := strings.Split(v[0].(string), "{}")
            var r string
            for i := 0; i < len(x)-1; i++ {
                if i+1 < len(v) {
                    r = fmt.Sprint(r, x[i], v[i+1])
                } else {
                    r = fmt.Sprint(r, x[i])
                }
            }
            r = fmt.Sprint(r, x[len(x)-1])
            for i := 0; i < len(v)-len(x); i++ {
                r = fmt.Sprint(r, v[len(x)+i])
            }
            return r
        }
        return utils.Concat(v...)
    default:
        return utils.Concat(v...)
    }
}

func buildContentWithError(err error, v ...interface{}) string {
    return utils.Concat2("\n", buildContent(v...), err.Error(), "\n", string(debug.Stack()))
}

func (l *Logger) Print(v ...interface{}) {
    // layout, arg1, arg2, ..., argn
    l.Debug(v...)
}
