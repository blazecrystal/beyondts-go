package utils

import "testing"

func TestZipFile(t *testing.T) {
	srcs := []string{"D:\\workspaces\\workspace-go\\beyondts\\logs\\test.log"}
	dst := "d:\\test.zip"
	ZipFile(dst, srcs...)
}
