package main

import (
	"bytes"
	"io/ioutil"
	"encoding/json"
	"bufio"
	"fmt"
	"github.com/tteeoo/sest/lib"
	"io"
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

func readPwd() (string, error) {
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
	reader := bufio.NewReader(os.Stdin)
	password, err := reader.ReadString('\n')
	exec.Command("stty", "-F", "/dev/tty", "echo").Run()
	print("\n")
	if err != nil {
		return "", err
	}
	return password, nil
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
		fmt.Println("sest: error: invalid arguments, run 'sest --help' for usage")
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
			fmt.Println("sest: error: provide a name for the new container")
			os.Exit(1)
		}

		for _, char := range []string{".", "/", "@", "\\", "&", "*"} {
			if strings.Contains(args[1], char) {
				fmt.Println("sest: error: invalid character '" + char + "' in container name")
				os.Exit(1)
			}
		}

		if _, err := os.Stat(contDir + "/" + args[1] + ".cont.json"); !os.IsNotExist(err) {
			fmt.Println("sest: error: a container with that name already exists")
			os.Exit(1)
		}

		fmt.Print("sest: new container password: ")
		password, err := readPwd()
		if err != nil {
			fmt.Println("sest: error:", err)
			os.Exit(1)
		}

		c, err := lib.NewContainer(args[1], contDir, password)
		if err != nil {
			fmt.Println("sest: error:", err)
			os.Exit(1)
		}

		err = c.Write()
		if err != nil {
			fmt.Println("sest: error:", err)
			os.Exit(1)
		}

		os.Exit(0)

	// Import a container from json
	case "imp":
		if len(args) < 2 {
			fmt.Println("sest: error: provide a name for the new container")
			os.Exit(1)
		} else if len(args) < 3 {
			fmt.Println("sest: error: provide a file path to import from")
			os.Exit(1)
		}

		for _, char := range []string{".", "/", "@", "\\", "&", "*"} {
			if strings.Contains(args[1], char) {
				fmt.Println("sest: error: invalid character '" + char + "' in container name")
				os.Exit(1)
			}
		}

		if _, err := os.Stat(contDir + "/" + args[1] + ".cont.json"); !os.IsNotExist(err) {
			fmt.Println("sest: error: a container with that name already exists")
			os.Exit(1)
		}

		fmt.Print("sest: new container password: ")
		password, err := readPwd()
		if err != nil {
			fmt.Println("sest: error:", err)
			os.Exit(1)
		}

		c, err := lib.NewContainer(args[1], contDir, password)
		if err != nil {
			fmt.Println("sest: error:", err)
			os.Exit(1)
		}

		b, err := ioutil.ReadFile(args[2])
		if err != nil {
			fmt.Println("sest: error:", err)
			os.Exit(1)
		}

		var data map[string]string
		err = json.Unmarshal(b, &data)
		if err != nil {
			fmt.Println("sest: error:", err)
			os.Exit(1)
		}

		err = c.SetData(data, password)
		if err != nil {
			fmt.Println("sest: error:", err)
			os.Exit(1)
		}

		err = c.Write()
		if err != nil {
			fmt.Println("sest: error:", err)
			os.Exit(1)
		}

		os.Exit(0)

	// Deletes a container
	case "del":
		if len(args) < 2 {
			fmt.Println("sest: error: provide a name for the container to delete")
			os.Exit(1)
		}

		if _, err := os.Stat(contDir + "/" + args[1] + ".cont.json"); os.IsNotExist(err) {
			fmt.Println("sest: error: a container with that name does not exist")
			os.Exit(1)
		}

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("sest: press enter to delete the container '" + args[1] + "'")
		_, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("sest: error:", err)
			os.Exit(1)
		}

		err = os.Remove(contDir + "/" + args[1] + ".cont.json")
		if err != nil {
			fmt.Println("sest: error:", err)
			os.Exit(1)
		}
		os.Exit(0)

	// Deletes a key from a container
	case "rm":
		if len(args) < 2 {
			fmt.Println("sest: error: provide a name for the container to remove from")
			os.Exit(1)
		} else if len(args) < 3 {
			fmt.Println("sest: error: provide a key to remove from the container")
			os.Exit(1)
		}

		if _, err := os.Stat(contDir + "/" + args[1] + ".cont.json"); os.IsNotExist(err) {
			fmt.Println("sest: error: a container with that name does not exist")
			os.Exit(1)
		}


		fmt.Print("sest: container password: ")
		password, err := readPwd()
		if err != nil {
			fmt.Println("sest: error:", err)
			os.Exit(1)
		}

		c, err := lib.OpenContainer(args[1], contDir)
		if err != nil {
			fmt.Println("sest: error:", err)
			os.Exit(1)
		}

		data, err := c.GetData(password)
		if err != nil {
			fmt.Println("sest: error:", err)
			os.Exit(1)
		}

		if _, ok := data[args[2]]; ok {
			delete(data, args[2])
			err = c.SetData(data, password)
			if err != nil {
				fmt.Println("sest: error:", err)
				os.Exit(1)
			}

			err = c.Write()
			if err != nil {
				fmt.Println("sest: error:", err)
				os.Exit(1)
			}
			os.Exit(0)
		}
		fmt.Println("sest: error: the key '" + args[2] + "' does not exist in that container")
		os.Exit(1)

	// Stores a key-value pair in a container
	case "in":
		if len(args) < 2 {
			fmt.Println("sest: error: provide a name for the container to store in")
			os.Exit(1)
		} else if len(args) < 3 {
			fmt.Println("sest: error: provide a key to store in the container")
			os.Exit(1)
		}

		if _, err := os.Stat(contDir + "/" + args[1] + ".cont.json"); os.IsNotExist(err) {
			fmt.Println("sest: error: a container with that name does not exist")
			os.Exit(1)
		}

		fmt.Print("sest: container password: ")
		password, err := readPwd()
		if err != nil {
			fmt.Println("sest: error:", err)
			os.Exit(1)
		}

		c, err := lib.OpenContainer(args[1], contDir)
		if err != nil {
			fmt.Println("sest: error:", err)
			os.Exit(1)
		}

		data, err := c.GetData(password)
		if err != nil {
			fmt.Println("sest: error:", err)
			os.Exit(1)
		}

		if _, ok := data[args[2]]; ok {
			fmt.Println("sest: note that the key '" + args[2] + "' already exists in that container, if you continue its value will be changed")
		}

		fmt.Print("sest: new key value: ")
		value, err := readPwd()
		if err != nil {
			fmt.Println("sest: error:", err)
			os.Exit(1)
		}

		data[args[2]] = value[0 : len(value)-1]
		err = c.SetData(data, password)
		if err != nil {
			fmt.Println("sest: error:", err)
			os.Exit(1)
		}
		err = c.Write()
		if err != nil {
			fmt.Println("sest: error:", err)
			os.Exit(1)
		}
		os.Exit(0)

	// Gets the value from a key that is inside a container
	case "out", "cp":
		if len(args) < 2 {
			fmt.Println("sest: error: provide a name for the container to read from")
			os.Exit(1)
		} else if len(args) < 3 {
			fmt.Println("sest: error: provide a key to read from the container")
			os.Exit(1)
		}

		if _, err := os.Stat(contDir + "/" + args[1] + ".cont.json"); os.IsNotExist(err) {
			fmt.Println("sest: error: a container with that name does not exist")
			os.Exit(1)
		}

		fmt.Print("sest: container password: ")
		password, err := readPwd()
		if err != nil {
			fmt.Println("sest: error:", err)
			os.Exit(1)
		}

		c, err := lib.OpenContainer(args[1], contDir)

		data, err := c.GetData(password)
		if err != nil {
			fmt.Println("sest: error:", err)
			os.Exit(1)
		}

		if _, ok := data[args[2]]; ok {

			// copy to clipboard (needs xclip)
			if args[0] == "cp" {
				cmd := exec.Command("xclip")
				stdin, err := cmd.StdinPipe()
				if err != nil {
					fmt.Println("sest: error:", err)
					os.Exit(1)
				}

				go func() {
					defer stdin.Close()
					io.WriteString(stdin, data[args[2]])
				}()

				err = cmd.Run()
				os.Exit(0)
				if err != nil {
					fmt.Println("sest: error:", err)
					os.Exit(1)
				}
			}

			fmt.Println(data[args[2]])
			os.Exit(0)
		}
		fmt.Println("sest: error: the key '" + args[2] + "' does not exist in that container")
		os.Exit(1)

	// Lists all the keys in a container
	case "ln":
		if len(args) < 2 {
			fmt.Println("sest: error: provide a name for the container to read from")
			os.Exit(1)
		}

		if _, err := os.Stat(contDir + "/" + args[1] + ".cont.json"); os.IsNotExist(err) {
			fmt.Println("sest: error: a container with that name does not exist")
			os.Exit(1)
		}

		fmt.Print("sest: container password: ")
		password, err := readPwd()
		if err != nil {
			fmt.Println("sest: error:", err)
			os.Exit(1)
		}

		c, err := lib.OpenContainer(args[1], contDir)

		data, err := c.GetData(password)
		if err != nil {
			fmt.Println("sest: error:", err)
			os.Exit(1)
		}

		for key := range data {
			fmt.Println(key)
		}

		os.Exit(0)

	// Changes a container's password
	case "chp":
		if len(args) < 2 {
			fmt.Println("sest: error: provide a name for the container to change to password of")
			os.Exit(1)
		}

		if _, err := os.Stat(contDir + "/" + args[1] + ".cont.json"); os.IsNotExist(err) {
			fmt.Println("sest: error: a container with that name does not exist")
			os.Exit(1)
		}

		fmt.Print("sest: container password: ")
		password, err := readPwd()
		if err != nil {
			fmt.Println("sest: error:", err)
			os.Exit(1)
		}

		c, err := lib.OpenContainer(args[1], contDir)
		if err != nil {
			fmt.Println("sest: error:", err)
			os.Exit(1)
		}

		fmt.Print("sest: new container password: ")
		newPassword, err := readPwd()
		if err != nil {
			fmt.Println("sest: error:", err)
			os.Exit(1)
		}

		newC, err := c.ChangePassword(password, newPassword)
		if err != nil {
			fmt.Println("sest: error:", err)
			os.Exit(1)
		}

		err = newC.Write()
		if err != nil {
			fmt.Println("sest: error:", err)
			os.Exit(1)
		}
		os.Exit(0)

	// Exports container as json
	case "exp":
		if len(args) < 2 {
			fmt.Println("sest: error: provide a name of the container to export")
			os.Exit(1)
		} else if len(args) < 3 {
			fmt.Println("sest: error: provide a file path to export to")
			os.Exit(1)
		}

		if _, err := os.Stat(contDir + "/" + args[1] + ".cont.json"); os.IsNotExist(err) {
			fmt.Println("sest: error: a container with that name does not exist")
			os.Exit(1)
		}

		fmt.Print("sest: container password: ")
		password, err := readPwd()
		if err != nil {
			fmt.Println("sest: error:", err)
			os.Exit(1)
		}

		c, err := lib.OpenContainer(args[1], contDir)
		if err != nil {
			fmt.Println("sest: error:", err)
			os.Exit(1)
		}

		data, err := c.GetData(password)
		if err != nil {
			fmt.Println("sest: error:", err)
			os.Exit(1)
		}

		buffer := &bytes.Buffer{}
		encoder := json.NewEncoder(buffer)
		encoder.SetEscapeHTML(false)
		encoder.SetIndent("", " ")
		err = encoder.Encode(data)
		if err != nil {
			fmt.Println("sest: error:", err)
			os.Exit(1)
		}

		err = ioutil.WriteFile(args[2], buffer.Bytes(), 0600)
		if err != nil {
			fmt.Println("sest: error:", err)
			os.Exit(1)
		}

		os.Exit(0)

	case "-V", "--version":
		fmt.Println("sest version 0.1.7")

	case "-h", "--help":
		fmt.Println("sest: secure strings\n\n" +
			"usage: sest [-h | --help]\n" +
			"            [-V | --version]\n" +
			"            [<command> [arguments]]\n\n" +
			"commands:\n" +
			"    ls                        lists all containers\n" +
			"    mk  <container>           makes a new container\n" +
			"    ln  <container>           lists all keys in a container\n" +
			"    chp <container>           changes a container's password\n" +
			"    del <container>           deletes a container; asks for confirmation\n" +
			"    in  <container> <key>     stores a new key-value pair in a container or changes an existing key\n" +
			"    cp  <container> <key>     copies the value of a key from a container to the clipboard (requires xclip)\n" +
			"    rm  <container> <key>     removes a key-value pair from a container\n" +
			"    out <container> <key>     prints out the value of a key from a container\n" +
			"    exp <container> <path>    export a container to a json file\n" +
			"    imp <container> <path>    import a container from a json file\n\n" +
			"licensed under the BSD 2-clause license\n" +
			"set the environment variable 'SEST_DIR' to the directory where you want containers to be stored, it defaults to '$HOME/.sest'")

	default:
		fmt.Println("sest: error: invalid arguments, run 'sest --help' for usage")
		os.Exit(1)
	}

	os.Exit(0)
}
