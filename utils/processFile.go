package utils

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
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

func AddFileInfo(filename string, path string, token string, codePath string, saveTime int) {
	lock.Lock()
	defer lock.Unlock()

	fileType := "file"
	if strings.Contains(filename, ".") {
		fileType = filename[strings.LastIndex(filename, ".")+1:]
	}

	fileInfo := FileInfo{
		Token:        token,
		Filename:     filename,
		Filetype:     fileType,
		FilePath:     path,
		SaveMinutes:  saveTime,
		IsDownloaded: false,
		QRCode:       codePath,
	}
	if saveTime > 0 {
		fileInfo.SaveFromTime = true
	}

	fileInfos[token] = fileInfo
}

func GetFileInfo(token string) (FileInfo, error) {
	lock.RLock()
	defer lock.RUnlock()

	info, ok := fileInfos[token]
	if !ok {
		return info, fmt.Errorf("file not existed or already expired")
	}

	return info, nil
}

func DelFile(token string) {
	lock.Lock()
	defer lock.Unlock()

	fileInfo, ok := fileInfos[token]
	if !ok { // already deleted maybe
		return
	}

	err := os.Remove(fileInfo.FilePath)
	if err != nil {
		fmt.Printf("error in delete file %s because %v\n", fileInfo.Filename, err)
	}

	err = os.Remove(fileInfo.QRCode)
	if err != nil {
		fmt.Printf("error in delete file %s's qrcode because %v\n", fileInfo.Filename, err)
	}

	delete(fileInfos, token)
	fmt.Printf("delete file %s\n", fileInfo.Filename)

}

func DelFileAfter(token string, after int) {
	lock.RLock()
	defer lock.RUnlock()

	if f, ok := fileInfos[token]; ok {
		timer := time.NewTimer(time.Duration(after) * time.Minute)
		fmt.Printf("delete file %s after %d minutes\n", f.Filename, after)
		go func() {
			<-timer.C
			DelFile(token)
		}()
	}
}
