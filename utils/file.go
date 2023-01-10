package utils

import (
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

func Copy(src, dest string) {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return
	}
	_, err = os.Stat(dest)
	if os.IsNotExist(err) {
		if srcInfo.IsDir() {
			_ = os.MkdirAll(dest, 0666)
		}
	}
	if srcInfo.IsDir() {
		log.Infoln("复制文件")
		dir, err := os.ReadDir(src)
		if err != nil {
			return
		}
		for _, entry := range dir {
			if entry.IsDir() {
				Copy(path.Join(src, entry.Name()), path.Join(dest, entry.Name()))
				continue
			} else {
				file, _ := os.Open(path.Join(src, entry.Name()))
				newFile, _ := os.OpenFile(path.Join(dest, entry.Name()), os.O_RDWR|os.O_CREATE, 0666)
				_, err := io.Copy(newFile, file)
				if err != nil {
					return
				}
			}
		}
	} else {
		_ = os.MkdirAll(filepath.Dir(dest), 0666)
		file, _ := os.Open(src)
		newFile, _ := os.OpenFile(dest, os.O_RDWR|os.O_CREATE, 0666)
		_, err := io.Copy(newFile, file)
		if err != nil {
			return
		}
	}
}

func DownloadFile(url, filePath string) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(response.Body)
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}
	return nil
}

func FileExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	} else {
		return true
	}
}
