package logging

import (
	"log"
	"os"
	"path/filepath"

	"github.com/artnikel/vacancystats/internal/constants"
)

type Logger struct {
	Info  *log.Logger
	Error *log.Logger
}

func NewLogger(dir string) (*Logger, error) {
	err := os.MkdirAll(dir, constants.DirPerm)
	if err != nil {
		return nil, err
	}

	logPath := filepath.Join(dir, "app.log")
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, constants.FilePerm)
	if err != nil {
		return nil, err
	}

	return &Logger{
		Info:  log.New(logFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		Error: log.New(logFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
	}, nil
}
