package utils

import (
	"fmt"
	"strings"
	"sync"
)

var (
	fileInfos = make(map[string]FileInfo, 0)
	lock      = sync.RWMutex{}
)

type FileInfo struct {
	Token    string
	Filename string
	Filetype string
	FilePath string

	SaveFromTime bool
	SaveMinutes  int

	IsDownloaded bool

	QRCode string
}

func AddFileInfo(filename string, path string, token string, codePath string) {
	lock.Lock()
	defer lock.Unlock()

	fileType := "file"
	if strings.Contains(filename, ".") {
		fileType = filename[strings.LastIndex(filename, ".")+1:]
	}

	fileInfos[token] = FileInfo{
		Token:        token,
		Filename:     filename,
		Filetype:     fileType,
		FilePath:     path,
		SaveFromTime: false,
		SaveMinutes:  0,
		IsDownloaded: false,
		QRCode:       codePath,
	}
}

func GetFileInfo(token string) (FileInfo, error) {
	lock.RLock()
	defer lock.RUnlock()

	info, ok := fileInfos[token]
	if !ok {
		return info, fmt.Errorf("file not exist")
	}

	return info, nil
}
