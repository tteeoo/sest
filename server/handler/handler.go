package handler

import (
	"github.com/tteeoo/sest/server/util"
	"io/ioutil"
	"os"
	"strings"
)

var contDir = os.Getenv("SEST_SERVER_CONT_DIR")

var Whitelist = []string{}
var Blacklist = []string{}
var whitelistFile = os.Getenv("SEST_SERVER_WHITELIST")
var blacklistFile = os.Getenv("SEST_SERVER_BLACKLIST")

func init() {

	// Load host whitelist/blacklist
	if len(whitelistFile) > 0 {
		b, err := ioutil.ReadFile(whitelistFile)
		if err != nil {
			util.Logger.Fatal("sest-server: error: ", err)
		}

		Whitelist = strings.Split(string(b), "\n")
		util.Logger.Println("sest-server: whitelist ("+whitelistFile+") loaded, only accepting hosts:", Whitelist)
	} else if len(blacklistFile) > 0 {
		b, err := ioutil.ReadFile(blacklistFile)
		if err != nil {
			util.Logger.Fatal("sest-server: error: ", err)
		}

		Blacklist = strings.Split(string(b), "\n")
		util.Logger.Println("sest-server: blacklist ("+blacklistFile+") loaded, denying hosts:", Blacklist)
	}

	// Get container directory
	if len(contDir) == 0 {
		contDir = os.Getenv("HOME") + "/.sest"
	}

	// Make contDir if it does not exist
	if _, err := os.Stat(contDir); os.IsNotExist(err) {
		util.Logger.Println("sest-server: making directory at ", contDir)
		os.Mkdir(contDir, 0700)
	}
}
