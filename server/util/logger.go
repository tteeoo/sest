package util

import (
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

var Logger *log.Logger
var LogFile os.File

func init() {

	// Log to terminal and a file
	LogFile, err := os.Create("./log/go-website-" + strconv.Itoa(int(time.Now().Unix())) + ".log")
	if err != nil {
		log.Fatal(err)
	}

	LogFile.Sync()

	Logger = log.New(io.MultiWriter(LogFile, os.Stdout), "", log.Ldate|log.Ltime|log.Lshortfile)
}
