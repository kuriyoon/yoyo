package yoyoSystem

import (
	"io/ioutil"
	"os"
)


/*
CreateDirIfNotExist(LOG_DIR)
	CreateDirIfNotExist(READ_DIR)
 */
func CreateDirIfNotExist(dir string) (rErr error) {
	if _, rErr := os.Stat(dir); os.IsNotExist(rErr) {
		rErr= os.MkdirAll(dir, 0755)
	}
	return
}

func ReadFileTest(filePath string) (rErr error) {
	_, rErr = ioutil.ReadFile(filePath)
	return
}

func WriteFile(filePath string,Contents string) (rErr error) {
	f, rErr := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if rErr != nil {
		return
	}
	defer f.Close()
	_, rErr = f.WriteString(Contents)
	return
}