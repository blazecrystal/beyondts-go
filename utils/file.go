package utils

import (
	"os"
	"strings"
	"time"
)

// newName is only filename part not contains any dirs, this func will only rename file and never move it to another dir
func Rename(oldFile *os.File, newName string) (string, error) {
	path := oldFile.Name()[:strings.LastIndex(oldFile.Name(), string(os.PathSeparator))+1]
	newPath := Concat(path, newName)
	err := os.Rename(oldFile.Name(), newPath)
	return newPath, err
}

// newName is only filename part not contains any dirs, this func will only rename file and never move it to another dir
// oldFilePath must contain whole dirs
func Rename2(oldFilePath, newName string) (string, error) {
	oldFilePath = ToLocalFilePath(oldFilePath)
	if oldFilePath[len(oldFilePath)-1] == os.PathSeparator {
		oldFilePath = oldFilePath[0 : len(oldFilePath)-1]
	}
	path := oldFilePath[:strings.LastIndex(oldFilePath, string(os.PathSeparator))+1]
	newPath := Concat(path, newName)
	err := os.Rename(oldFilePath, newPath)
	return newPath, err
}

// newDir is only dir path, not contains filename, this func will only move file to a new dir and never change its name
func MoveTo(oldFile *os.File, newDir string) error {
	name := oldFile.Name()[strings.LastIndex(oldFile.Name(), string(os.PathSeparator))+1:]
	return os.Rename(oldFile.Name(), Concat(ToLocalDirPath(newDir), name))
}

// just change \ and / to xosutils.PathSeparator
func ToLocalFilePath(path string) string {
	tmp := strings.TrimSpace(path)
	sep := string(os.PathSeparator)
	tmp = strings.Replace(tmp, "\\", sep, -1)
	tmp = strings.Replace(tmp, "/", sep, -1)
	return tmp
}

// do ToLocalFilePath and append additional file.PathSeparator if path is not end with file.PathSeparator after changing
func ToLocalDirPath(path string) string {
	tmp := ToLocalFilePath(path)
	if tmp[len(tmp)-1] != os.PathSeparator {
		tmp = Concat(tmp, string(os.PathSeparator))
	}
	return tmp
}

func GetParentDir(file *os.File) string {
	return GetParentDir2(file.Name())
}

func GetParentDir2(path string) string {
	tmp := ToLocalFilePath(path)
	if path[len(path)-1] == os.PathSeparator {
		tmp = tmp[:len(tmp)-1]
	}
	lastPathSept := strings.LastIndex(tmp, string(os.PathSeparator))
	return tmp[:lastPathSept+1]
}

func GetFileName(file *os.File) string {
	return GetFileName2(file.Name())
}

func GetFileName2(path string) string {
	pathSept := string(os.PathSeparator)
	tmp := path
	lastPathSept := strings.LastIndex(tmp, pathSept)
	if lastPathSept == len(tmp)-1 {
		tmp = tmp[:len(tmp)-1]
		lastPathSept = strings.LastIndex(tmp, pathSept)
	}
	return tmp[lastPathSept+1:]
}

func ListDir(dir, keywords string, includeSubDir, caseSensitive bool) ([]string, error) {
	var list []string
	tmp := ToLocalDirPath(dir)
	toListDir, err := os.OpenFile(tmp, os.O_RDONLY, 0600)
	defer toListDir.Close()
	if err != nil {
		return nil, err
	}
	subs, err := toListDir.Readdir(-1)
	if err != nil {
		return nil, err
	}
	for _, sub := range subs {
		if matches(sub, keywords, includeSubDir, caseSensitive) {
			list = append(list, Concat(tmp, sub.Name()))
		}
	}
	return list, nil
}

func matches(info os.FileInfo, keywords string, includeSubDir, caseSensitive bool) bool {
	name := GetFileName2(info.Name())
	tmp := strings.TrimSpace(keywords)
	if !includeSubDir && info.IsDir() {
		return false
	}
	if !caseSensitive {
		name = strings.ToLower(name)
		tmp = strings.ToLower(tmp)
	}
	if strings.Contains(keywords, "*") {
		parts := strings.Split(tmp, "*")
		for _, part := range parts {
			if part != "" {
				index := strings.Index(name, part)
				if index < 0 {
					return false
				}
				name = name[index+len(part):]
			}
		}
		return true
	} else {
		return strings.Contains(name, tmp)
	}
	return false
}

func FindEarliest(files []string) (int, string) {
	earliestFile := ""
	earliestModTime := time.Now().UnixNano()
	index := -1
	for i, file := range files {
		tmp, err := os.Stat(file)
		if err != nil {
			continue
		}
		if tmp.ModTime().UnixNano() <= earliestModTime {
			earliestFile = file
			earliestModTime = tmp.ModTime().UnixNano()
			index = i
		}
	}
	return index, earliestFile
}

func SortFilesByModTime(files []string, desc bool) ([]string, error) {
	for i := 0; i < len(files); i++ {
		statI, err := os.Stat(files[i])
		if err != nil {
			return nil, err
		}
		for j := i + 1; j < len(files); j++ {
			statJ, err := os.Stat(files[j])
			if err != nil {
				return nil, err
			}
			if (!desc && statI.ModTime().UnixNano() > statJ.ModTime().UnixNano()) || (desc && statI.ModTime().UnixNano() < statJ.ModTime().UnixNano()) {
				files[i], files[j] = files[j], files[i]
			}
		}
	}
	return files, nil
}
