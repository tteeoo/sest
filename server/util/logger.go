package util

import (
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

var logDir = os.Getenv("SEST_SERVER_LOG_DIR")
var LogFile os.File
var Logger *log.Logger

func init() {

	if len(logDir) == 0 {
		logDir = "./log"
	}

	// Log to terminal and a file
	LogFile, err := os.Create(logDir + "/go-website-" + strconv.Itoa(int(time.Now().Unix())) + ".log")
	if err != nil {
		log.Fatal(err)
	}

	LogFile.Sync()

	Logger = log.New(io.MultiWriter(LogFile, os.Stdout), "", log.Ldate|log.Ltime)
}
