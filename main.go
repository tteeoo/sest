package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
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
			print("\n")

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
		fmt.Println("sest: press enter to delete the container " + args[1])
		_, _ = reader.ReadString('\n')

		err := os.Remove(contDir + "/" + args[1] + ".cont.json")
		if err != nil {
			fmt.Println("sest: error:", err)
			os.Exit(1)
		}
		os.Exit(0)

	// Deletes a key from a container
	case "rm":
		if len(args) < 2 {
			fmt.Println("sest: error: please provide a name for the container to remove from")
			os.Exit(1)
		} else if len(args) < 3 {
			fmt.Println("sest: error: please provide a key to remove from the container")
			os.Exit(1)
		}

		if _, err := os.Stat(contDir + "/" + args[1] + ".cont.json"); os.IsNotExist(err) {
			fmt.Println("sest: error: a container with that name does not exist")
			os.Exit(1)
		}

		exec.Command("stty", "-F", "/dev/tty", "-echo").Run()

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("sest: container password: ")
		password, _ := reader.ReadString('\n')
		print("\n")

		c, err := openContainer(args[1])
		if err != nil {
			fmt.Println("sest: error:", err)
			os.Exit(1)
		}

		data, err := c.getData(password)
		if err != nil {
			fmt.Println("sest: error:", err)
			os.Exit(1)
		}

		if _, ok := data[args[2]]; ok {
			delete(data, args[2])
			err = c.setData(data, password)
			if err != nil {
				fmt.Println("sest: error:", err)
				os.Exit(1)
			}
			os.Exit(0)
		}
		fmt.Println("sest: error: the key " + args[2] + " does not exist in that container")
		os.Exit(1)

	// Stores a key-value pair in a container
	case "in":
		if len(args) < 2 {
			fmt.Println("sest: error: please provide a name for the container to store in")
			os.Exit(1)
		} else if len(args) < 3 {
			fmt.Println("sest: error: please provide a key to store in the container")
			os.Exit(1)
		}

		if _, err := os.Stat(contDir + "/" + args[1] + ".cont.json"); os.IsNotExist(err) {
			fmt.Println("sest: error: a container with that name does not exist")
			os.Exit(1)
		}

		exec.Command("stty", "-F", "/dev/tty", "-echo").Run()

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("sest: container password: ")
		password, _ := reader.ReadString('\n')
		print("\n")

		c, err := openContainer(args[1])
		if err != nil {
			fmt.Println("sest: error:", err)
			os.Exit(1)
		}

		data, err := c.getData(password)
		if err != nil {
			fmt.Println("sest: error:", err)
			os.Exit(1)
		}

		if _, ok := data[args[2]]; ok {
			fmt.Println("sest: error: the key " + args[2] + " already exists in that container")
			os.Exit(1)
		}

		exec.Command("stty", "-F", "/dev/tty", "-echo").Run()

		reader = bufio.NewReader(os.Stdin)
		fmt.Print("sest: new key value: ")
		value, _ := reader.ReadString('\n')
		print("\n")

		data[args[2]] = value
		err = c.setData(data, password)
		c.write()
		if err != nil {
			fmt.Println("sest: error:", err)
			os.Exit(1)
		}
		os.Exit(0)

	// Gets the value from a key that is inside a container
	case "out":
		if len(args) < 2 {
			fmt.Println("sest: error: please provide a name for the container to read from")
			os.Exit(1)
		} else if len(args) < 3 {
			fmt.Println("sest: error: please provide a key to read from the container")
			os.Exit(1)
		}

		if _, err := os.Stat(contDir + "/" + args[1] + ".cont.json"); os.IsNotExist(err) {
			fmt.Println("sest: error: a container with that name does not exist")
			os.Exit(1)
		}

		exec.Command("stty", "-F", "/dev/tty", "-echo").Run()

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("sest: container password: ")
		password, _ := reader.ReadString('\n')
		print("\n")

		c, err := openContainer(args[1])

		data, err := c.getData(password)
		if err != nil {
			fmt.Println("sest: error:", err)
			os.Exit(1)
		}

		if _, ok := data[args[2]]; ok {
			fmt.Print(args[2]+":", data[args[2]])
			os.Exit(0)
		}
		fmt.Println("sest: error: the key " + args[2] + " does not exist in that container")
		os.Exit(1)

	// Lists all the keys in a container
	case "ln":
		if len(args) < 2 {
			fmt.Println("sest: error: please provide a name for the container to read from")
			os.Exit(1)
		}

		if _, err := os.Stat(contDir + "/" + args[1] + ".cont.json"); os.IsNotExist(err) {
			fmt.Println("sest: error: a container with that name does not exist")
			os.Exit(1)
		}

		exec.Command("stty", "-F", "/dev/tty", "-echo").Run()

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("sest: container password: ")
		password, _ := reader.ReadString('\n')
		print("\n")

		c, err := openContainer(args[1])

		data, err := c.getData(password)
		if err != nil {
			fmt.Println("sest: error:", err)
			os.Exit(1)
		}

		for _, key := range data {
			fmt.Print(key)
		}

		os.Exit(0)

	case "-V", "--version":
		fmt.Println("sest: version: 0.1.0")

	case "-h", "--help":
		fmt.Println("sest: secure strings\n\n" +
			"usage:\n\tsest [--version | -V] | [--help | -h] | [<command> [arguments]]\n\n" +
			"commands:\n\tmk <container name>: makes a new container, will ask for a master password\n" +
			"\n\tdel <container name>: deletes a container, will ask for confirmation\n" +
			"\n\tls: lists all containers\n" +
			"\n\tin <container name> <key name>: stores a new key-value pair in a container, will ask for a master password and a value\n" +
			"\n\tout <container name> <key name>: prints out the value of a key from a container, will ask for a master password\n" +
			"\n\tln: lists all keys in a container, will ask for a master password\n" +
			"\n\trm <container name> <key name>: removes a key-value pair from a container, will ask for a master password\n\n" +
			"source hosted on GitHub at https://github.com/tteeoo/sest")

	default:
		fmt.Println("sest: error: invalid arguments, run `sest --help` for usage")
		os.Exit(1)
	}

	os.Exit(0)
}
