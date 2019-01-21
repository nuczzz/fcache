// +build windows

package fcache

import (
	"fmt"
	"os"
	"syscall"
)

// FileTime file's ctime, atime, mtime with nanoseconds
type FileTime struct {
	CreateTime int64
	AccessTime int64
	ModifyTime int64
}

func GetFileTime(file string) (*FileTime, error) {
	fileInfo, err := os.Stat(file)
	if err != nil {
		return nil, err
	}
	if stat, ok := fileInfo.Sys().(*syscall.Win32FileAttributeData); ok {
		return &FileTime{
			CreateTime: stat.CreationTime.Nanoseconds(),
			AccessTime: stat.LastAccessTime.Nanoseconds(),
			ModifyTime: stat.LastWriteTime.Nanoseconds(),
		}, nil
	}
	return nil, fmt.Errorf("not support file info in current platform")
}
