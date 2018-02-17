package utils

import (
	"archive/zip"
	"errors"
	"io"
	"os"
)

func ZipFile(zipPath string, srcPaths ...string) error {
	if srcPaths == nil || len(srcPaths) < 1 {
		return errors.New("no src file path, at least 1 src file path should be given")
	}
	zipFile, err := os.OpenFile(zipPath, os.O_RDWR|os.O_CREATE, 0600)
	defer zipFile.Close()
	if err != nil {
		return err
	}
	zipWriter := zip.NewWriter(zipFile)
	defer func() {
		zipWriter.Flush()
		zipWriter.Close()
	}()
	for _, srcPath := range srcPaths {
		writer, err := zipWriter.Create(GetFileName2(srcPath))
		if err != nil {
			return err
		}
		srcFile, err := os.OpenFile(srcPath, os.O_RDONLY, 0600)
		defer srcFile.Close()
		if err != nil {
			return err
		}
		PipeTo(srcFile, writer, io.EOF)
	}
	return nil
}
