package appenders

import (
    "strconv"
    "strings"
    "time"
    "runtime"
    "os"
    "bytes"
    "github.com/blazecrystal/beyondts-go/utils"
    "errors"
)

const (
    DEFAULT_UGO = 0600
    PREFIX_ENV  = "${"
    SUFFIX_ENV  = "}"

    ROLLING_DAILY           = "daily"
    ROLLING_SIZE            = "size"
    DEFAULT_KEEP            = 5
    DEFAULT_SIZE_THRESHOLD  = "1m"
    DEFAULT_DAILY_THRESHOLD = "1d"
)

type LocalFileAppender struct {
    name, layout, file, rollingType, threshold string
    zip                                        bool
    keep                                       int
    logFile                                    *os.File
}

func NewLocalFileAppender(name, layout, file, rollingType, threshold string, zip bool, keep int) (*LocalFileAppender, error) {
    if len(file) < 1 {
        return nil, errors.New("file not assigned : " + file)
    }
    file = utils.ToLocalFilePath(file)
    logFile := openLogFile(file)
    if logFile == nil {
        return nil, errors.New("file not found and can't be created : " + file)
    }
    if rollingType != ROLLING_DAILY && rollingType != ROLLING_SIZE {
        rollingType = ROLLING_DAILY
    }
    if keep < 1 {
        keep = DEFAULT_KEEP
    }
    if len(layout) < 1 {
        layout = DEFAULT_LAYOUT
    }
    // set all illegal threshold to ""
    unit := threshold[len(threshold)-1:]
    period, err := strconv.Atoi(threshold[: len(threshold)-1])
    if err != nil || period < 1 || (rollingType == ROLLING_DAILY && (unit != "h" && unit != "d" && unit != "m" && unit != "y")) ||
        (rollingType == ROLLING_SIZE && (unit != "b" && unit != "k" && unit != "m" && unit != "g")) {
        threshold = ""
    }
    if len(threshold) < 2 {
        switch rollingType {
        case ROLLING_DAILY:
            threshold = DEFAULT_DAILY_THRESHOLD
        case ROLLING_SIZE:
            threshold = DEFAULT_SIZE_THRESHOLD
        }
    }
    return &LocalFileAppender{name: name, layout: layout, file: file, rollingType: rollingType, threshold: threshold, zip: zip, keep: keep, logFile: logFile}, nil
}

func openLogFile(path string) *os.File {
    file, err := os.OpenFile(getFilepath(path), os.O_RDWR|os.O_CREATE|os.O_APPEND, DEFAULT_UGO)
    if err != nil {
        file.Close()
        return nil
    }
    return file
}

func getFilepath(path string) string {
    buf := bytes.Buffer{}
    tmp := path
    var prefixIndex, suffixIndex int
    for prefixIndex, suffixIndex = strings.Index(tmp, PREFIX_ENV), strings.Index(tmp, SUFFIX_ENV); prefixIndex > -1 && suffixIndex-prefixIndex > 2; prefixIndex, suffixIndex = strings.Index(tmp, PREFIX_ENV), strings.Index(tmp, SUFFIX_ENV) {
        buf.WriteString(tmp[:prefixIndex])
        buf.WriteString(os.Getenv(tmp[prefixIndex+2: suffixIndex]))
        tmp = tmp[suffixIndex+1: ]
    }
    buf.WriteString(tmp[suffixIndex+1: ])
    return buf.String()
}

func (a *LocalFileAppender) GetName() string {
    return a.name
}

func (a *LocalFileAppender) GetType() string {
    return TYPE_LOCAL_FILE
}

func (a *LocalFileAppender) GetId() string {
    return TYPE_LOCAL_FILE + "#" + a.file + "#" + a.layout + "#" + a.rollingType + "#" + a.threshold + "#" + strconv.FormatBool(a.zip) + "#" + strconv.Itoa(a.keep)
}

func (a *LocalFileAppender) Stop() {
    if a.logFile != nil {
        a.logFile.Close()
    }
}

func (a *LocalFileAppender) Write(loggerName, level, content string) {
    // rolling
    a.rolling()
    // write
    a.write(a.buildLog(loggerName, level, content))
}

func (a *LocalFileAppender) write(log string) {
    a.logFile.WriteString(log + "\n")
}

func (a *LocalFileAppender) rolling() {
    var suffix string
    switch a.rollingType {
    case ROLLING_DAILY:
        suffix = a.dailyBackupSuffix()
    case ROLLING_SIZE:
        suffix = a.sizeBackupSuffix()
    }
    if len(suffix) > 0 {
        // need rolling
        // delete oldest backup
        a.deleteOldest()
        a.backup(suffix)
        logFile := openLogFile(a.file)
        if logFile == nil {
            return
        }
        a.logFile = logFile
    }
}

func (a *LocalFileAppender) deleteOldest() {
    filePath := a.logFile.Name()
    fileName := filePath[strings.LastIndex(filePath, string(os.PathSeparator))+1:]
    parentDir := utils.GetParentDir(a.logFile)
    oldFiles, err := utils.ListDir(parentDir, fileName+".*", false, false)
    if err != nil {
        return
    }
    if len(oldFiles) > 0 && a.rollingType == ROLLING_SIZE {
        oldFiles, _ = utils.SortFilesByModTime(oldFiles, false)
        renameAllOldFiles(filePath, oldFiles)
    }
    if len(oldFiles) > a.keep {
        os.Remove(oldFiles[len(oldFiles)-1])
    }
}

func renameAllOldFiles(filePath string, oldFiles []string) {
    // oldFiles already sorted, from 1 to n
    for i, oldFile := range oldFiles {
        if strings.HasSuffix(oldFile, ".zip") {
            os.Rename(oldFile, filePath+"."+strconv.Itoa(i+2)+".zip")
        } else {
            os.Rename(oldFile, filePath+"."+strconv.Itoa(i+2))
        }
    }
}

// always return "1", if exist old files, should rename
func (a *LocalFileAppender) sizeBackupSuffix() string {
    stat, err := a.logFile.Stat()
    if err != nil {
        return ""
    }
    length := len(a.threshold)
    size, _ := strconv.Atoi(a.threshold[: length-1]) // never appear error, has set when appender is created
    switch a.threshold[length-1:] {
    case "k":
        size = size * 1024
    case "m":
        size = size * 1024 * 1024
    case "g":
        size = size * 1024 * 1024 * 1024
    }
    if stat.Size() >= int64(size) {
        return "1"
    }
    return ""
}

// return suffix of backup file
// if no need to backup, suffix will be ""
func (a *LocalFileAppender) dailyBackupSuffix() string {
    stat, err := a.logFile.Stat()
    if err != nil {
        return ""
    }
    length := len(a.threshold)
    period, _ := strconv.Atoi(a.threshold[: length-1]) // never appear error, has set when appender is created
    var fileTime, now, format string
    switch a.threshold[length-1:] {
    case "h":
        format = "2006010215"
        fileTime = stat.ModTime().Add(time.Duration(period) * time.Hour).Format(format)
    case "d":
        format = "20060102"
        fileTime = stat.ModTime().AddDate(0, 0, period).Format(format)
    case "m":
        format = "200601"
        fileTime = stat.ModTime().AddDate(0, period, 0).Format(format)
    case "y":
        format := "2006"
        fileTime = stat.ModTime().AddDate(period, 0, 0).Format(format)
    }
    now = time.Now().Format(format)
    if strings.Compare(now, fileTime) >= 0 {
        return fileTime
    }
    return ""
}

func (a *LocalFileAppender) backup(suffix string) {
    newFilePath := a.logFile.Name() + "." + suffix
    a.logFile.Close()
    if a.zip {
        newFilePath += ".zip"
        a.zipFile(newFilePath)
        os.Remove(a.logFile.Name())
    } else {
        os.Rename(a.logFile.Name(), newFilePath)
    }
}

func (a *LocalFileAppender) zipFile(zipFilePath string) {
    utils.ZipFile(zipFilePath, a.logFile.Name())
}

func (a *LocalFileAppender) buildLog(loggerName, level, content string) string {
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
