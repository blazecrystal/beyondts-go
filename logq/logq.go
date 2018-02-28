package logq

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"github.com/blazecrystal/beyondts-go/properties"
	"github.com/blazecrystal/beyondts-go/utils"
)

const (
	DEFAULT_LOGGER         = "default"
	DEFAULT_APPENDER       = "default"
	DEFAULT_CONFIG_REFRESH = -1
	DEFAULT_CONFIG_DEBUG   = false
	prefix_logger          = "LogQ.logger."
	prefix_appender        = "LogQ.appender."
	debug_prefix           = "LogQ(debug mode) :"
	KEY_SYSTEM_BUILDIN     = "LogQ.system.appender"
)

// use LogQ in singleton mode, you can just use "logq.LogQ" after loading its properties
var (
	logQ              *LogQ
	appenderFlag      = make(chan string)
	rwlock            sync.RWMutex
	configMonitorFlag byte
	appendersStatus   = make(map[string]string) // appenders running stauts, "+" for running, "-" for stopped
)

type LogQ struct {
	refresh    int // default -1 for no refreshing
	loggers    map[string]*Logger
	appenders  map[string]Appender
	configFile string // config file is a properties file
	debugMode  bool
}

func GetLogger(name string) *Logger {
	if logQ == nil {
		defaultConfig()
		runDaemon()
	}
	return getLogger(name)
}

func getLogger(name string) *Logger {
	logger := logQ.loggers[name]
	if logger == nil {
		// not found current logger, found upper level
		indexOfSept := strings.LastIndex(name, "/")
		if indexOfSept > 0 {
			name = name[0:indexOfSept]
			debugMsg("logger named \"", name, "\" not found, continue to find its super logger named \"", name, "\"")
			logger = getLogger(name)
			debugMsg("logger named \"", logger.name, "\" found")
		} else {
			debugMsg("logger named \"", name, "\" and its super loggers not found, using \"default\" logger")
			logger = logQ.loggers[DEFAULT_LOGGER]
		}
	}
	return logger;
}

func GetAppender(name string) Appender {
	if logQ == nil {
		defaultConfig()
	}
	appender := logQ.appenders[name]
	if appender == nil {
		debugMsg("appender named \"", name, "\" not found, using \"default\" appender")
		appender = logQ.appenders[DEFAULT_APPENDER]
	}
	return appender
}

func Go(configFile string) {
	config(configFile)
	runDaemon()
}

func End() {
	stopDaemon()
}

func Refresh() bool {
	prop, err := properties.LoadPropertiesFromFile(logQ.configFile)
	if err != nil {
		debugMsg("failed to read config file :", logQ.configFile)
		return false
	}
	refreshLoggers(prop)
	refreshAppenders(prop)
	return true
}

func runDaemon() {
	// run appenders
	for _, appender := range logQ.appenders {
		go appender.Run(appenderFlag)
	}
	checkAppendersStatus(true)

	// run config monitor
	if logQ.refresh > 0 {
		go runConfigMonitor()
	}
	debugMsg("LogQ running")
}

// we should sure that "+" and "-" should be both written
func checkAppendersStatus(running bool) {
	for {
		flag := <-appenderFlag
		if flag[0:1] == "+" {
			debugMsg("appender named \"", flag[1:], "\" running")
		} else if flag[0:1] == "-" {
			debugMsg("appender named \"", flag[1:], "\" stopped")
		}
		if running {
			appendersStatus[flag[1:]] = flag[0:1]
			if len(appendersStatus) == len(logQ.appenders) {
				debugMsg(strconv.Itoa(calcAppenders("+")), " appenders running, ", strconv.Itoa(calcAppenders("-")), " appenders stopped")
				break
			}
		} else {
			delete(appendersStatus, flag[1:])
			if len(appendersStatus) == 0 {
				debugMsg("all appenders stopped")
				break
			}
		}
	}
}

func calcAppenders(flag string) int {
	count := 0
	for _, v := range appendersStatus {
		if v == flag {
			count++
		}
	}
	return count
}

func stopDaemon() {
	// stop config monitor
	if logQ.refresh > 0 {
		go stopConfigMonitor()
	}

	// stop appenders
	for _, appender := range logQ.appenders {
		go appender.Stop()
	}
	checkAppendersStatus(false)
	debugMsg("LogQ stopped")
}

func runConfigMonitor() {
	for {
		rwlock.RLock()
		if configMonitorFlag == 1 {
			break
		} else {
			rwlock.RUnlock()
			utils.Sleep(logQ.refresh)
			Refresh()
		}
	}
	rwlock.RUnlock()
}

func stopConfigMonitor() {
	rwlock.Lock()
	configMonitorFlag = 1
	rwlock.Unlock()
}

// we will create a default logger named "default" in loggers
// this logger "default" will be used when no logger found by given name and for all loggers baddly defined
// and also a default appender named "default" of type "stdout" will be created in appenders
// this appender "default" will be used when no appender found by given name
// attention: you can set "default" log & "default" appender in properties manually
func config(configFile string) {
	// use LogQ in singleton mode
	if logQ == nil {
		configFile = utils.ToLocalFilePath(configFile)
		prop, err := properties.LoadPropertiesFromFile(configFile)
		if err != nil {
			defaultConfig()
			GetLogger(DEFAULT_LOGGER).Warn(utils.Concat("can't read config file : ", configFile, ", use default settings !"))
		} else {
			logQ = &LogQ{refresh: -1, loggers: make(map[string]*Logger), appenders: make(map[string]Appender), configFile: configFile, debugMode: false}
			// load settings from properties
			logQ.debugMode = prop.GetBool("LogQ.config.debug", DEFAULT_CONFIG_DEBUG)
			logQ.refresh = prop.GetInt("LogQ.config.refresh", DEFAULT_CONFIG_REFRESH)
			readLoggers(prop)
			readAppenders(prop)
			//redirectSystemLog(prop)
		}
	}
}

func redirectSystemLog(prop *properties.Properties) {
	// TODO
}

func defaultConfig() {
	if logQ == nil {
		logQ = &LogQ{refresh: -1, loggers: make(map[string]*Logger), appenders: make(map[string]Appender), configFile: "", debugMode: false}
		logQ.loggers[DEFAULT_LOGGER] = createDefaultLogger()
		logQ.appenders[DEFAULT_APPENDER] = createDefaultAppender()
		debugMsg("warning! using default configurations !")
	}
}

func readLoggers(prop *properties.Properties) {
	logQ.loggers[DEFAULT_LOGGER] = createDefaultLogger()
	for _, key := range prop.Keys() {
		if strings.HasPrefix(key, prefix_logger) {
			name := key[len(prefix_logger):]
			tmp := septValues(prop.Get(key))
			logQ.loggers[name] = newLogger(name, tmp[0], logQ.debugMode, tmp[1:]...)
			debugMsg("found logger named \"", name, "\"")
		}
	}
}
func septValues(value string) []string {
	vs := strings.Split(value, ",")
	for k, v := range vs {
		vs[k] = strings.TrimSpace(v)
	}
	return vs
}

func refreshLoggers(prop *properties.Properties) {
	var newLoggerNames []string
	newLoggerNames = append(newLoggerNames, DEFAULT_LOGGER)
	for _, key := range prop.Keys() {
		if strings.HasPrefix(key, prefix_logger) {
			name := key[len(prefix_logger):]
			newLoggerNames = append(newLoggerNames, name)
			tmp := septValues(prop.Get(key))
			if logQ.loggers[name] == nil {
				// logger not exist
				logQ.loggers[name] = newLogger(name, tmp[0], logQ.debugMode, tmp[1:]...)
				debugMsg("found new logger named \"", name, "\"")
			} else {
				// this named logger already exist
				logQ.loggers[name].update(name, tmp[0], tmp[1:]...)
				debugMsg("logger named \"", name, "\" updated")
			}
		}
	}
	// remove not configed loggers
	for current, _ := range logQ.loggers {
		if !utils.ExistInSliece(newLoggerNames, current) {
			delete(logQ.loggers, current)
		}
	}
}

func readAppenders(prop *properties.Properties) {
	logQ.appenders[DEFAULT_APPENDER] = createDefaultAppender()
	for _, logger := range logQ.loggers {
		if logger.appenders != nil && len(logger.appenders) > 0 {
			for _, appenderName := range logger.appenders {
				_, exist := logQ.appenders[appenderName]
				if !exist {
					prefix := utils.Concat("", prefix_appender, appenderName, ".")
					tmp := make(map[string]string)
					tmp["name"] = appenderName
					for _, key := range prop.Keys() {
						if strings.HasPrefix(key, prefix) {
							tmp[key[len(prefix):]] = prop.Get(key)
						}
					}
					logQ.appenders[appenderName] = newAppender(tmp["type"], logQ.debugMode, tmp)
					debugMsg("found appender named \"", appenderName, "\"")
				}
			}
		}
	}
	/*sysAppenderName := strings.TrimSpace(prop.Get(KEY_SYSTEM_BUILDIN))
	  if sysAppenderName != "" {
	      _, exist := logQ.appenders[sysAppenderName]
	      if !exist {
	          prefix := utils.Concat(prefix_appender, sysAppenderName, ".")
	          tmp := make(map[string]string)
	          tmp["name"] = sysAppenderName
	          for _, key := range prop.Keys() {
	              if strings.HasPrefix(key, prefix) {
	                  tmp[key[len(prefix):]] = prop.Get(key)
	              }
	          }
	          logQ.appenders[sysAppenderName] = NewAppender(tmp["type"], tmp)
	          debug("found appender named \"", sysAppenderName, "\"")
	      }
	  }*/
}

func refreshAppenders(prop *properties.Properties) {
	var newAppenderNames []string
	newAppenderNames = append(newAppenderNames, DEFAULT_APPENDER)
	for _, logger := range logQ.loggers {
		if logger.appenders != nil && len(logger.appenders) > 0 {
			for _, appenderName := range logger.appenders {
				if !utils.ExistInSliece(newAppenderNames, appenderName) {
					// have not updated
					newAppenderNames = append(newAppenderNames, appenderName)
					prefix := utils.Concat(prefix_appender, appenderName, ".")
					tmp := make(map[string]string)
					tmp["name"] = appenderName
					for _, key := range prop.Keys() {
						if strings.HasPrefix(key, prefix) {
							tmp[key[len(prefix):]] = prop.Get(key)
						}
					}
					_, exist := logQ.appenders[appenderName]
					if exist && strings.EqualFold(logQ.appenders[appenderName].GetType(), tmp["type"]) {
						logQ.appenders[appenderName].Update(tmp)
						debugMsg("appender named \"", appenderName, "\" updated")
					} else {
						logQ.appenders[appenderName] = newAppender(tmp["type"], logQ.debugMode, tmp)
						logQ.appenders[appenderName].Run(appenderFlag)
						debugMsg("found new appender named \"", appenderName, "\" which is newly defined or type changed")
					}
				}
			}
		}
	}
	for appenderName, _ := range logQ.appenders {
		if !utils.ExistInSliece(newAppenderNames, appenderName) {
			delete(logQ.appenders, appenderName)
		}
	}
}

func debugMsg(msgs ...interface{}) {
	if logQ.debugMode {
		fmt.Print(debug_prefix)
		fmt.Print(msgs...)
		fmt.Println()
	}
}
