package appenders

import (
	"bufio"
	"bytes"
	"os"
	"strconv"
	"strings"
	"time"
	"github.com/blazecrystal/beyondts-go/utils"
)

const (
	default_ugo = 0600
	prefix_env  = "${"
	suffix_env  = "}"
)

type LocalFileAppender struct {
	name, layout, file, dailyRolling, threshold, zip, keep string
	log                                                    chan string
}

func NewLocalFileAppender(attrs map[string]string) *LocalFileAppender {
	tmp := &LocalFileAppender{name: attrs["name"], layout: attrs["layout"], file: utils.ToLocalFilePath(attrs["file"]), dailyRolling: attrs["dailyRolling"], threshold: attrs["threshold"], zip: attrs["zip"], keep: attrs["keep"], log: make(chan string)}
	return tmp
}

func (a *LocalFileAppender) GetType() string {
	return Type_LocalFile
}

func (a *LocalFileAppender) WriteLog(fills map[string]string) {
	a.log <- fillLayout(a.layout, fills)
}

func (a *LocalFileAppender) Stop() {
	a.log <- string(CMD_END)
}

func (a *LocalFileAppender) Run(flag chan string) {
	// open file & build writer
	logFile, err := openLogFile(a.file)
	defer logFile.Close()
	if err != nil {
		debugMsg("can't open log file : ", a.file, ", this appender will be skipped\n", err)
		flag <- utils.Concat("-", a.name)
		return
	}
	writer := bufio.NewWriter(logFile)
	flag <- utils.Concat("+", a.name)
	for {
		tmp := <-a.log
		if bytes.Equal([]byte(tmp), CMD_END) {
			break
		}
		// write tmp to local file
		// redirect to new file & backup old one
		if a.backup(logFile) {
			logFile, err = openLogFile(a.file)
			defer logFile.Close()
			if err != nil {
				debugMsg("can't open log file : ", a.file, ", this appender will be skipped\n", err)
				flag <- utils.Concat("-", a.name)
				return
			}
			writer = bufio.NewWriter(logFile)
		}
		writer.WriteString(tmp)
		writer.WriteRune('\n')
		writer.Flush()
	}

	flag <- utils.Concat("-", a.name)
}

// if current file not opened at the end of the func
func (a *LocalFileAppender) backup(currentFile *os.File) bool {
	rst := false
	if strings.EqualFold(a.dailyRolling, "true") {
		rst = a.dailyBackup(currentFile)
	} else if len(a.threshold) > 0 {
		rst = a.thresholdBackup(currentFile)
	}
	// if not backup daily or by threshold, we do nothing
	return rst
}

func (a *LocalFileAppender) thresholdBackup(currentFile *os.File) bool {
	stat, err := currentFile.Stat()
	if err != nil {
		debugMsg("can't open log file stat", err)
		return false
	}
	maxSize, err := thresholdToFileSize(a.threshold)
	if err != nil {
		debugMsg("illegal threshold format, threshold will never effect")
		return false
	}
	if stat.Size() < maxSize {
		return false
	} else {
		oldPath := currentFile.Name()
		oldName := oldPath[strings.LastIndex(oldPath, string(os.PathSeparator))+1:]
		maxBackup, err := strconv.Atoi(a.keep)
		if err != nil {
			debugMsg("parameter \"keep\" is not a valid number", err)
			return false
		}
		backupLogs, err := utils.ListDir(utils.GetParentDir(currentFile), utils.Concat(oldName, ".*"), false, false)
		if err != nil {
			debugMsg("can't list backup log files", err)
			return false
		}
		if backupLogs != nil {
			if len(backupLogs) > maxBackup-1 { // current file hasn't been backuped
				backupLogs = deleteOldest(backupLogs)
			}
			backupLogs, err = utils.SortFilesByModTime(backupLogs, false)
			if err != nil {
				debugMsg("some backup log file may have not exist", err)
				return false
			}
			backupCount := len(backupLogs)
			for i, backupLog := range backupLogs {
				if strings.HasSuffix(strings.ToLower(backupLog), ".zip") {
					utils.Rename2(backupLog, utils.Concat(oldName, ".", strconv.Itoa(backupCount-i), ".zip"))
				} else {
					utils.Rename2(backupLog, utils.Concat(oldName, ".", strconv.Itoa(backupCount-i)))
				}
			}
		}
		currentFile.Close()
		needZip := a.zip == "true"
		if needZip {
			zipFile(oldPath, "0")
			os.Remove(oldPath)
		} else {
			utils.Rename2(oldPath, utils.Concat(oldName, ".0"))
		}
		return true
	}
	return false
}

func thresholdToFileSize(threshold string) (int64, error) {
	nums, err := utils.SliceAtoi(strings.Split(threshold, "|"))
	if err != nil {
		return -1, err
	}
	length := len(nums)
	size := 0
	for i, v := range nums {
		size += v << (10 * uint(length-1-i))
	}
	return int64(size), nil
}

func (a *LocalFileAppender) dailyBackup(currentFile *os.File) bool {
	stat, err := currentFile.Stat()
	if err != nil {
		debugMsg("can't open log file stat", err)
		return false
	}
	fileTime := stat.ModTime().Format("20060102")
	now := time.Now().Format("20060102")
	if strings.Compare(fileTime, now) < 0 {
		oldPath := currentFile.Name()
		oldName := oldPath[strings.LastIndex(oldPath, string(os.PathSeparator))+1:]
		currentFile.Close()
		needZip := a.zip == "true"
		if needZip {
			zipFile(oldPath, fileTime)
			os.Remove(oldPath)
		} else {
			utils.Rename2(oldPath, utils.Concat(oldName, ".", fileTime))
		}
		maxBackup, err := strconv.Atoi(a.keep)
		if err != nil {
			debugMsg("parameter \"keep\" is not a valid number", err)
			return true
		}
		backupLogs, err := utils.ListDir(utils.GetParentDir2(oldPath), utils.Concat(oldName, ".*"), false, false)
		if err != nil {
			debugMsg("some backup log file may have not exist", err)
			return true
		}
		if backupLogs != nil && len(backupLogs) > maxBackup {
			deleteOldest(backupLogs)
		}
		return true
	}
	return false
}

func (a *LocalFileAppender) Update(attrs map[string]string) {
	a.name, a.layout, a.file, a.dailyRolling, a.threshold, a.zip, a.keep = attrs["name"], attrs["layout"], attrs["file"], attrs["dailyRolling"], attrs["threshold"], attrs["zip"], attrs["keep"]
}

func deleteOldest(backups []string) []string {
	if backups != nil && len(backups) > 0 {
		index, earliestFile := utils.FindEarliest(backups)
		os.Remove(earliestFile)
		return append(backups[0:index], backups[index+1:]...)
	}
	return nil
}

func zipFile(toZipFilePath, suffix string) {
	utils.ZipFile(utils.Concat(toZipFilePath, ".", suffix, ".zip"), toZipFilePath)
}

func openLogFile(path string) (*os.File, error) {
	file, err := os.OpenFile(getFilepath(path), os.O_RDWR|os.O_CREATE|os.O_APPEND, default_ugo)
	if err != nil {
		file.Close()
		return nil, err
	}
	return file, nil
}

func getFilepath(path string) string {
	buf := bytes.Buffer{}
	tmp := path
	var prefixIndex, suffixIndex int
	for prefixIndex, suffixIndex = strings.Index(tmp, prefix_env), strings.Index(tmp, suffix_env); prefixIndex > -1 && suffixIndex-prefixIndex > 2; prefixIndex, suffixIndex = strings.Index(tmp, prefix_env), strings.Index(tmp, suffix_env) {
		buf.WriteString(tmp[0:prefixIndex])
		buf.WriteString(os.Getenv(tmp[prefixIndex+2 : suffixIndex]))
		tmp = tmp[suffixIndex+1:]
	}
	buf.WriteString(tmp[suffixIndex+1:])
	return buf.String()
}
