package appenders

import (
    "fmt"
    "strings"
    "time"
    "runtime"
    "os"
    "strconv"
    "github.com/blazecrystal/beyondts-go/utils"
)

type StdoutAppender struct {
    name, layout string
    logs []string
}

func NewStdoutAppender(name, layout string) *StdoutAppender {
    if len(layout) < 1 {
        layout = DEFAULT_LAYOUT
    }
    return &StdoutAppender{name: name, layout: layout}
}

func (a *StdoutAppender) GetName() string {
    return a.name
}

func (a *StdoutAppender) GetType() string {
    return TYPE_STDOUT
}

func (a *StdoutAppender) GetId() string {
    return TYPE_STDOUT + "#" + a.layout
}

func (a *StdoutAppender) Stop() {
    // nothing to release
}

func (a *StdoutAppender) Write(loggerName, level, content string) {
    log := a.buildLog(loggerName, level, content)
    fmt.Println(log)
}

func (a *StdoutAppender) buildLog(loggerName, level, content string) string {
    tmp := a.layout
    tmp = strings.Replace(tmp, TIMESTAMP, time.Now().Format("2006-01-02 15:04:05.000"), -1)
    tmp = strings.Replace(tmp, LEVEL, level, -1)
    _, file, line, ok := runtime.Caller(2) // skip Caller & current stack
    if ok {
        tmp = strings.Replace(tmp, LONG_FILE_NAME, file, -1)
        file = utils.ToLocalFilePath(file)
        lastPathSeptIndex := strings.LastIndex(file, string(os.PathSeparator))
        tmp = strings.Replace(tmp, SHORT_FILE_NAME, file[lastPathSeptIndex+1: ], -1)
        tmp = strings.Replace(tmp, LINE_NUMBER, strconv.Itoa(line), -1)
    } else {
        tmp = strings.Replace(tmp, LONG_FILE_NAME, UNKNOWN, -1)
        tmp = strings.Replace(tmp, SHORT_FILE_NAME, UNKNOWN, -1)
        tmp = strings.Replace(tmp, LINE_NUMBER, UNKNOWN, -1)
    }
    tmp = strings.Replace(tmp, LOGGER_NAME, loggerName, -1)
    tmp = strings.Replace(tmp, CONTENT, content, -1)
    return tmp
}
