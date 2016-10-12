package banner

import (
	"github.com/RebuildTools/rebuild-agent/helpers"
	"github.com/Sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"github.com/CrowdSurge/banner"
	"fmt"
	"github.com/fatih/color"
	"time"
)

// Run is called by the main of
// the agent to run this command.
func Run(logger *logrus.Logger) {
	injectExitHandler(logger)

	initrdVer, err := helpers.GetInitrdVersion()

	if err != nil { logger.WithError(err).Fatal("Failed to get the Rebuild Initrd version from file") }

	banner.Print("rebuild agent")

	fmt.Println()

	fmt.Println(
		color.GreenString("Agent v%s", helpers.AgentVersion),
		"-",
		color.YellowString("Initrd v%s", initrdVer),
	)

	for { time.Sleep(60 * time.Second) }
}

// injectExitHandler sets up a signal
// handler that is called when a user
// tries to close the application with
// a SigTerm or Interrupt
func injectExitHandler(logger *logrus.Logger) {
	exitHandle := make(chan os.Signal, 1)
	signal.Notify(exitHandle, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-exitHandle
		logger.Debug("The banner cannot be closed, if you really need to, please use SIGKILL")
	}()
}
