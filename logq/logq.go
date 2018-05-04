package logq

import (
    "github.com/blazecrystal/beyondts-go/logq/appenders"
    "github.com/beevik/etree"
    "strings"
    "github.com/blazecrystal/beyondts-go/utils"
    "fmt"
    "runtime/debug"
    "strconv"
)

type LogQ struct {
    file       string
    scanPeriod int
    debug      bool
    as         map[string]appenders.Appender
    loggers    map[string]*Logger
}

const (
    DEBUG_PREFIX = "LogQ(debug mode) : "

    SEPT_LOGGER_NAME = "/"

    DEFAULT_SCAN_PERIOD = -1
    DEFAULT_DEBUG       = true
)

var (
    logQ *LogQ // singleton mode
    status chan string
    monitorStatus bool
)

/*
 * start process :
 * 1. load config and build logq instance
 *      1.1. if config not correctly loaded, use default config to build logq
 *      1.2. if config loaded correctly, build logq
 * 2. start all appenders
 * 3. check appenders' status
 *      3.1. if appender not started correctly
 *          3.1.1. delete appender from logq's appender map
 *          3.1.2. remove them from logger's appender ref
 *          3.1.3. if logger has no appender running, delete logger from logq's logger map
 * 4. start config monitor
 *      4.1. sleep some time
 *      4.2. reload config
 *          4.2.1. load config
 *          4.2.2. start newly added appenders
 *          4.2.3. add new loggers
 *          4.2.4. stop newly removed appenders
 *          4.2.5. remove appenders from existed loggers
 *          4.2.6. remove loggers has no running appenders
 * -----------------------
 * stop process :
 * 1. stop config monitor
 * 2. stop appenders running
 * 3. release res
 */

func Config(xmlConfigFile string) {
    doc := etree.NewDocument()
    err := doc.ReadFromFile(xmlConfigFile)
    if err != nil {
        DebugMsg("can't load config xml, use default configuration ", err)
        buildInstanceWithDefaultConfig()
    } else {
        buildInstance(xmlConfigFile, doc.SelectElement("logq"))
    }
    go startConfigMonitor()
}

func GetLogger(name string) *Logger {
    if logQ == nil {
        buildInstanceWithDefaultConfig()
    }
    return getLogger(name)
}

func Stop() {
    if monitorStatus {
        stopConfigMonitor()
    }
    for _, a := range logQ.as {
        a.Stop()
    }
}

func buildInstanceWithDefaultConfig() {
    if logQ == nil {
        logQ = &LogQ{}
    }
    logQ.file = ""
    logQ.scanPeriod = DEFAULT_SCAN_PERIOD
    logQ.debug = DEFAULT_DEBUG
    logQ.as = make(map[string]appenders.Appender, 1)
    logQ.loggers = make(map[string]*Logger, 1)
    logQ.as[appenders.DEFAULT_APPENDER] = appenders.CreateDefaultAppender()
    logQ.loggers[DEFAULT_LOGGER] = CreateDefaultLogger()
}

func buildInstance(xmlConfigFile string, emt *etree.Element) {
    if logQ == nil {
        logQ = &LogQ{}
    }
    logQ.file = xmlConfigFile
    logQ.scanPeriod = utils.Atoi(getAttrValue(emt, "scanPeriod"), DEFAULT_SCAN_PERIOD)
    logQ.debug = utils.ParseBool(getAttrValue(emt, "debug"), DEFAULT_DEBUG)
    logQ.loggers = createLoggers(emt)
    DebugMsg(len(logQ.loggers), " loggers created")
    logQ.as = createAppenders(emt.SelectElements("appender"))
    DebugMsg(len(logQ.as), " appenders created")
}

func createLoggers(emt *etree.Element) map[string]*Logger {
    loggers := make(map[string]*Logger, 1)
    loggers[DEFAULT_LOGGER] = CreateDefaultLogger()
    emts := emt.SelectElements("logger")
    for _, loggerEmt := range emts {
        l := buildLogger(emt, loggerEmt)
        loggers[l.name] = l
    }
    return loggers
}

func buildLogger(emt, loggerEmt *etree.Element) *Logger {
    l := NewLogger(getAttrValue(loggerEmt, "name"), strings.ToLower(getAttrValue(loggerEmt, "level")))
    refEmts := loggerEmt.SelectElements("appender-ref")
    if len(refEmts) > 0 {
        l.as = make([]string, 0)
        for _, refEmt := range refEmts {
            aName := getAttrValue(refEmt, "ref")
            if !utils.ExistInSlice(l.as, aName) && emt.FindElement("./appender[@name='" + aName + "']") != nil {
                l.as = append(l.as, aName)
            }
        }
    }
    return l
}

func getAllReferedAppenders() []string {
    ref := make([]string, 0)
    for _, logger := range logQ.loggers {
        for _, a := range logger.as {
            if !utils.ExistInSlice(ref, a) {
                ref = append(ref, a)
            }
        }
    }
    return ref
}

func createAppenders(emts []*etree.Element) map[string]appenders.Appender {
    as := make(map[string]appenders.Appender, 0)
    ref := getAllReferedAppenders()
    var a appenders.Appender
    var err error
    for _, emt := range emts {
        if !utils.ExistInSlice(ref, getAttrValue(emt, "name")) {
            // if appender not refered, don't create it
            break
        }
        switch strings.ToLower(getAttrValue(emt, "type")) {
        case appenders.TYPE_LOCAL_FILE:
            rollingEmt := emt.SelectElement("rolling")
            a, err = appenders.NewLocalFileAppender(getAttrValue(emt, "name"),
                getText(emt.SelectElement("layout"), appenders.DEFAULT_LAYOUT), emt.SelectElement("file").Text(),
                    strings.ToLower(getAttrValue(rollingEmt, "type")), strings.ToLower(getAttrValue(rollingEmt, "threshold")),
                        getBoolAttrValue(rollingEmt, "zip"), getIntAttrValue(rollingEmt, "keep"))
        case appenders.TYPE_REMOTE_FILE:
        case appenders.TYPE_SYSLOG:
            serverEmt := emt.SelectElement("server")
            a, err = appenders.NewSyslogAppender(getAttrValue(emt, "name"), emt.SelectElement("layout").Text(),
                getAttrValue(serverEmt, "host"), getIntAttrValue(serverEmt, "port"))
        case appenders.TYPE_DATABASE:
            conEmt := emt.SelectElement("connection")
            sqlEmt := emt.SelectElement("sql")
            paramEmts := sqlEmt.SelectElements("param")
            params := make([]string, len(paramEmts))
            for i, paramEmt := range paramEmts {
                index := getIntAttrValue(paramEmt, "index") - 1
                if index < 1 {
                    index = i
                }
                params[index] = paramEmt.Text()
            }
            a, err = appenders.NewDatabaseAppender(getAttrValue(emt, "name"), getAttrValue(conEmt, "driver"),
                getAttrValue(conEmt, "url"), getAttrValue(sqlEmt, "sqlstr"),
                    getIntAttrValue(sqlEmt, "maxRidLen"), params)
        case appenders.TYPE_STDOUT:
            a = appenders.NewStdoutAppender(getAttrValue(emt, "name"),
                getText(emt.SelectElement("layout"), appenders.DEFAULT_LAYOUT))
        }
        if err == nil {
            as[a.GetName()] = a
        }
    }
    return as
}

func startConfigMonitor() {
    if logQ.scanPeriod > 0 {
        monitorStatus = true
        for monitorStatus {
            utils.Sleep(logQ.scanPeriod)
            if monitorStatus {
                refreshAll()
                DebugMsg("logq refreshed!")
            }
        }
    }
}

func stopConfigMonitor() {
    monitorStatus = false
}

func refreshAll() {
    doc := etree.NewDocument()
    err := doc.ReadFromFile(logQ.file)
    if err != nil {
        return // do nothing
    }
    emt := doc.SelectElement("logq")
    logQ.debug = utils.ParseBool(getAttrValue(emt, "debug"), true)
    tmp := utils.Atoi(getAttrValue(emt, "scanPeriod"), DEFAULT_SCAN_PERIOD)
    if tmp != logQ.scanPeriod {
        logQ.scanPeriod = tmp
        stopConfigMonitor()
        go startConfigMonitor()
    }
    refreshLoggers(emt)
    refreshAppenders(emt.SelectElements("appender"))
}

// if loggers created not in new config, delete them
// if loggers created in new config, updated them
// if there are some loggers newly added in config, create and run them
func refreshLoggers(emt *etree.Element) {
    loggers := createLoggers(emt)
    for name, logger := range logQ.loggers {
        newLogger := loggers[name]
        if newLogger == nil {
            delete(logQ.loggers, name)
            break
        }
        if logger.NeedUpdate(loggers[name]) {
            logger.UpdateTo(newLogger)
        }
    }
    for name, logger := range loggers {
        oldLogger := logQ.loggers[name]
        if oldLogger == nil {
            logQ.loggers[name] = logger
        }
    }
}

// delete all non-refered appenders
// if updated in new config, update appenders
// if new added refered appenders, add them
func refreshAppenders(emts []*etree.Element) {
    ref := getAllReferedAppenders()
    as := createAppenders(emts)
    for name, a := range logQ.as {
        if !utils.ExistInSlice(ref, name) {
            //a.Stop()
            delete(logQ.as, name)
            break
        }
        if appenderNeedUpdate(a, as[name]) {
            updateAppender(a, as[name])
        }
    }
    for name, a := range as {
        if _, ok := logQ.as[name]; !ok {
            logQ.as[name] = a
        }
    }
}

func appenderNeedUpdate(oldAppender, newAppender appenders.Appender) bool {
    return oldAppender.GetId() != newAppender.GetId()
}

func updateAppender(oldAppender, newAppender appenders.Appender) {
    oldAppender.Stop()
    delete(logQ.as, oldAppender.GetName())
    logQ.as[newAppender.GetName()] = newAppender
}

func getAttrValue(emt *etree.Element, key string) string {
    attr := emt.SelectAttr(key)
    if attr == nil {
        return ""
    }
    return attr.Value
}

func getBoolAttrValue(emt *etree.Element, key string) bool {
    tmp := getAttrValue(emt, key)
    v, err := strconv.ParseBool(tmp)
    if err != nil {
        v = false
    }
    return v
}

func getIntAttrValue(emt *etree.Element, key string) int {
    tmp := getAttrValue(emt, key)
    v, err := strconv.Atoi(tmp)
    if err != nil {
        v = 0
    }
    return v
}

func getText(emt *etree.Element, dflt string) string {
    txt := emt.Text()
    if len(txt) < 1 {
        return dflt
    }
    return txt
}

func getAppender(appenderName string) appenders.Appender {
    return logQ.as[appenderName]
}

// get logger from loggers
// if loggers is nil or len < 0, return default logger
func getLogger(loggerName string) *Logger {
    if len(loggerName) < 1 {
        return logQ.loggers[DEFAULT_LOGGER]
    }
    logger := logQ.loggers[loggerName]
    if logger != nil {
        return logger
    }
    lastSeptIndex := strings.LastIndex(loggerName, SEPT_LOGGER_NAME)
    if lastSeptIndex < 0 {
        return logQ.loggers[DEFAULT_LOGGER]
    }
    loggerName = loggerName[0: lastSeptIndex]
    return getLogger(loggerName)
}

func DebugMsg(msgs ...interface{}) {
    if logQ.debug {
        fmt.Print(DEBUG_PREFIX)
        for _, msg := range msgs {
            switch msg.(type) {
            case error:
                err := msg.(error)
                tmp := err.Error() + "\n" + string(debug.Stack()) + "\n"
                fmt.Print(tmp)
            default:
                fmt.Print(msg)
            }
        }
        fmt.Println()
    }
}
