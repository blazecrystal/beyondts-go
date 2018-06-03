package appenders

type Appender interface {
    GetName() string
    GetType() string
    GetId() string
    Write(loggerName, level, content string)
    Stop() // used to release resources, eg. file, connection
}

var (
    CMD_STOP = []byte{1, 1, 9}
)

const (
    DEFAULT_APPENDER   = "default"
    DEFAULT_LAYOUT = "%t% %p% %ln% %sf%(%n%) %c%"

    TYPE_STDOUT      = "stdout"
    TYPE_LOCAL_FILE  = "local_file"
    TYPE_REMOTE_FILE = "remote_file"
    TYPE_SYSLOG      = "syslog"
    TYPE_DATABASE    = "database"

    // %t% %p% %lf% %sf% %n% %sl% %c%
    TIMESTAMP        = "%t%"
    LEVEL            = "%p%"
    SHORT_FILE_NAME  = "%sf%"
    LONG_FILE_NAME   = "%lf%"
    LINE_NUMBER      = "%n%"
    LOGGER_NAME      = "%ln%"
    CONTENT          = "%c%"
    RANDOM_STRING_ID = "%rid%"

    UNKNOWN = "unknown"
)

func CreateDefaultAppender() Appender {
    return NewStdoutAppender(DEFAULT_APPENDER, DEFAULT_LAYOUT)
}
