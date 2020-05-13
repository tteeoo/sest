package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"reflect"
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

	// Get args
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
			os.Exit(1)
		}

		containers, err := file.Readdirnames(0)
		if err != nil {
			fmt.Println("sest: error:", err)
			os.Exit(1)
		}

		for _, name := range containers {
			split := strings.SplitN(name, ".", 2)
			if len(split) == 2 {
				if split[1] == "cont.json" {
					fmt.Println(split[0])
				}

			}
		}

	// Makes a new container
	case "mk":
		if len(args) < 2 {
			fmt.Println("sest: error: please provide a name for the new container")
			os.Exit(1)
		}

		for _, char := range []string{".", "/", "@", "\\", "&", "*"} {
			if strings.Contains(args[1], char) {
				fmt.Println("sest: error: invalid character (" + char + ") in container name")
				os.Exit(1)
			}
		}

		if _, err := os.Stat(contDir + "/" + args[1] + ".cont.json"); os.IsNotExist(err) {
			exec.Command("stty", "-F", "/dev/tty", "-echo").Run()

			reader := bufio.NewReader(os.Stdin)
			fmt.Print("sest: new container password: ")
			password, _ := reader.ReadString('\n')

			cont, err := newContainer(args[1], password)
			if err != nil {
				fmt.Println("sest: error:", err)
				os.Exit(1)
			}

			err = cont.write()
			if err != nil {
				fmt.Println("sest: error:", err)
				os.Exit(1)
			}

			os.Exit(0)
		}

		fmt.Println("sest: error: a container with that name already exists")
		os.Exit(1)

	// Deletes a container
	case "del":
		if len(args) < 2 {
			fmt.Println("sest: error: please provide a name for the container to delete")
			os.Exit(1)
		}

		if _, err := os.Stat(contDir + "/" + args[1] + ".cont.json"); os.IsNotExist(err) {
			fmt.Println("sest: error: a container with that name does not exist")
			os.Exit(1)
		}

		exec.Command("stty", "-F", "/dev/tty", "-echo").Run()

		reader := bufio.NewReader(os.Stdin)
		fmt.Println("sest: you are about to delete the container " + args[1] + ", type the container's password to confirm")
		fmt.Print("container password: ")
		password, _ := reader.ReadString('\n')

		c, err := openContainer(args[1])
		if err != nil {
			fmt.Println("sest: error:", err)
			os.Exit(1)
		}

		validHash, err := bDecode(c.Master[0])
		if err != nil {
			fmt.Println("sest: error:", err)
			os.Exit(1)
		}

		salt, err := bDecode(c.Master[1])
		if err != nil {
			fmt.Println("sest: error:", err)
			os.Exit(1)
		}

		newHash := a2Hash(password, salt)

		if reflect.DeepEqual(newHash, validHash) {
			err = os.Remove(c.getPath())
			if err != nil {
				fmt.Println("sest: error:", err)
				os.Exit(1)
			}
			os.Exit(0)
		}

		fmt.Println("sest: error: invalid password for container", c.Name)
		os.Exit(1)

	case "rm":
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
