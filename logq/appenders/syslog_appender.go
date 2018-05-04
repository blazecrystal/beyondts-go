package appenders

import (
    "net"
    "errors"
    "strconv"
    "strings"
    "time"
    "runtime"
    "github.com/blazecrystal/beyondts-go/utils"
    "os"
)

const (
    NET_UDP      = "udp"
    DEFAULT_PORT = 514
)

type SyslogAppender struct {
    name, layout, host string
    port               int
    con                *net.UDPConn
}

func NewSyslogAppender(name, layout, host string, port int) (*SyslogAppender, error) {
    if len(layout) < 1 {
        layout = DEFAULT_LAYOUT
    }
    if len(host) < 1 {
        return nil, errors.New("host not assigned")
    }
    if port < 1 {
        port = DEFAULT_PORT
    }
    servAddr, err := net.ResolveUDPAddr(NET_UDP, host+":"+strconv.Itoa(port))
    if err != nil {
        return nil, err
    }
    con, err := net.DialUDP(NET_UDP, nil, servAddr)
    if err != nil {
        return nil, err
    }
    return &SyslogAppender{name: name, layout: layout, host: host, port: port, con: con}, nil
}

func (s *SyslogAppender) GetName() string {
    return s.name
}

func (s *SyslogAppender) GetType() string {
    return TYPE_SYSLOG
}

func (s *SyslogAppender) GetId() string {
    return TYPE_SYSLOG + "#" + s.layout + "#" + s.host + "#" + strconv.Itoa(s.port)
}

func (s *SyslogAppender) Stop() {
    if s.con != nil {
        s.con.Close()
    }
}

func (s *SyslogAppender) Write(loggerName, level, content string) {
    s.con.Write([]byte(s.buildLog(loggerName, level, content)))
}

func (a *SyslogAppender) buildLog(loggerName, level, content string) string {
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
