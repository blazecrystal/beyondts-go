package appenders

import (
	"fmt"
	"strings"
)

const (
	// %t% %p% %lf% %sf% %n% %sl% %c%
	Fills_Timestamp     = "%t%"
	Fills_Level         = "%p%"
	Fills_ShortFileName = "%sf%"
	Fills_LongFileName  = "%lf%"
	Fills_LineNum       = "%n%"
	Fills_LoggerName    = "%ln%"
	Fills_Content       = "%c%"
	RANDOM_STRING_ID    = "%rid%"

	// types of apppender
	Type_Stdout     = "stdout"
	Type_LocalFile  = "local_file"
	Type_RemoteFile = "remote_file"
	Type_Syslog     = "syslog"
	Type_Database   = "database"

	// default layout
	DEFAULT_LAYOUT = "%t% %p% %ln% %sf%(%n%) %c%"

	// prefix in debug mode
	debug_prefix = "LogQ(debug mode) :"
)

var (
	CMD_END   = []byte{1, 1, 9}
	debugMode bool
)

func fillLayout(layout string, fills map[string]string) string {
	tmp := layout
	// replace content first, so that, if content contains parameter flags, they will be replaced later
	tmp = strings.Replace(tmp, Fills_Content, fills[Fills_Content], -1)
	tmp = strings.Replace(tmp, Fills_Timestamp, fills[Fills_Timestamp], -1)
	tmp = strings.Replace(tmp, Fills_Level, fills[Fills_Level], -1)
	tmp = strings.Replace(tmp, Fills_LoggerName, fills[Fills_LoggerName], -1)
	tmp = strings.Replace(tmp, Fills_LongFileName, fills[Fills_LongFileName], -1)
	tmp = strings.Replace(tmp, Fills_ShortFileName, fills[Fills_ShortFileName], -1)
	tmp = strings.Replace(tmp, Fills_LineNum, fills[Fills_LineNum], -1)
	return tmp
}

func debugMsg(msgs ...interface{}) {
	if debugMode {
		var tmp []interface{}
		tmp = append(tmp, debug_prefix)
		tmp = append(tmp, msgs...)
		fmt.Println(tmp...)
	}
}

func SetDebugMode(debug bool) {
	debugMode = debug
}
