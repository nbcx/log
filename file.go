package log

import (
	"io"
	"os"
	"path/filepath"
)

func FileWrite(path string) io.Writer {
	_ = CreateDirIfNotExists(filepath.Dir(path))
	file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	return file
}

// CreateDirIfNotExists 目录不存在时创建目录
func CreateDirIfNotExists(path string) error {
	if !Exists(path) {
		return os.MkdirAll(path, os.ModePerm)
	}
	return nil
}

// Exists 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
