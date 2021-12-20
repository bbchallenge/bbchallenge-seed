package bbchallenge

import (
	"io/ioutil"
	"os"
)

func InitAppendFile(logFileName string, outputDirectory string) *os.File {
	ioutil.WriteFile(outputDirectory+logFileName, []byte(""), 0644)
	logFile, _ := os.OpenFile(outputDirectory+logFileName, os.O_APPEND|os.O_WRONLY, 0644)
	return logFile
}
