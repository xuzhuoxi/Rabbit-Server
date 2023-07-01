// Package utils
// Create on 2023/7/1
// @author xuzhuoxi
package utils

import (
	"github.com/xuzhuoxi/infra-go/filex"
	"github.com/xuzhuoxi/infra-go/osxu"
)

func FixFilePath(filePath string) string {
	if filex.IsFile(filePath) {
		return filePath
	}
	return filex.Combine(osxu.GetRunningDir(), filePath)
}

func FixDirPath(dirPath string) string {
	if filex.IsDir(dirPath) {
		return dirPath
	}
	return filex.Combine(osxu.GetRunningDir(), dirPath)
}

func FixPath(path string) string {
	if filex.IsExist(path) {
		return path
	}
	return filex.Combine(osxu.GetRunningDir(), path)
}
