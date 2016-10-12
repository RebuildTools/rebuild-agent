package main
// rebuild-agent is a lightweight
// agent run on a Linux live system
// to collect information about the
// host machine, this information
// is then returned to the Rebuild Core
//
// Author: Liam Haworth <liam@haworth.id.au>
//

import (
	"github.com/RebuildTools/rebuild-agent/banner"

	"os"
	"strconv"
	"github.com/Sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

// logger is used as a primary
// logging method for this agent
var logger = logrus.New()

// supportedCommands defines a map of
// commands to their function counter
// parts
var supportedCommands = map[string]func(*logrus.Logger){}

// init is called before the main
// function to initialize the agent
// with everything it might need before
// running
func init() {
	logger.Formatter = new(prefixed.TextFormatter)
	var logLevel int64

	if os.Getenv("LOGLEVEL") == "" {
		logLevel = 4
	} else {
		logLevel, _ = strconv.ParseInt(os.Getenv("LOGLEVEL"), 10, 0)
	}

	if logLevel >= 5 {
		logger.Level = logrus.DebugLevel
	} else if logLevel >= 4 {
		logger.Level = logrus.InfoLevel
	} else {
		logger.Level = logrus.WarnLevel
	}

	// Add supported commands
	supportedCommands["banner"] = banner.Run
}

func main() {
	if len(os.Args) < 2 {
		logger.Fatal("No command specified!")
	}

	if supportedCommands[os.Args[1]] == nil {
		logger.WithField("command", os.Args[1]).Fatal("Unsupported command specified")
	} else {
		supportedCommands[os.Args[1]](logger)
	}
}
