package main

import (
	"fmt"
	"os"
	"strings"
)

var contDir string = os.Getenv("SEST_DIR")

func init() {

	if len(contDir) == 0 {
		contDir = os.Getenv("HOME") + "/.sest"
	}
}

func main() {

	// Make dir if it does not exist
	if _, err := os.Stat(contDir); os.IsNotExist(err) {
		fmt.Println("sest: making directory at", contDir)
		os.Mkdir(contDir, 0700)
	}

	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("sest: error: invalid arguments, run `sest --help` for usage")
		os.Exit(1)
	}

	switch args[0] {

	// Prints the all the available containers
	case "ls":
		file, err := os.Open(contDir)
		if err != nil {
			fmt.Println("sest: error:", err)
		}

		containers, err := file.Readdirnames(0)
		if err != nil {
			fmt.Println("sest: error:", err)
		}

		for _, name := range containers {
			split := strings.SplitN(name, ".", 2)
			if len(split) == 2 {
				if split[1] == "cont.json" {
					fmt.Println(split[0])
				}

			}
		}

	case "mk":
	case "rm":

	case "del":
	case "in":
	case "out":

	case "-V", "--version":
		fmt.Println("sest: version: 0.1.0")

	case "-h", "--help":
		fmt.Println("sest: usage:")

	default:
		fmt.Println("sest: error: invalid arguments, run `sest --help` for usage")
		os.Exit(1)
	}

	os.Exit(0)

}
