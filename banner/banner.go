package banner
// banner is a command for the Rebuild
// Agent that presents a quick over view
// of the system the agent is running on

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
	helpers.HandleError(logger, "Retrieving Initrd version", err)

	kernelVer, err := helpers.GetKernelVersion()
	helpers.HandleError(logger, "Retrieving Kernel version", err)

	systemType, err := helpers.GetDMIValue("product_name")
	helpers.HandleError(logger, "Retrieving DMI header value", err)

	systemSerial, err := helpers.GetDMIValue("product_serial")
	helpers.HandleError(logger, "Retrieving DMI header value", err)

	systemVendor, err := helpers.GetDMIValue("chassis_vendor")
	helpers.HandleError(logger, "Retrieving DMI header value", err)

	systemUuid, err := helpers.GetDMIValue("product_uuid")
	helpers.HandleError(logger, "Retrieving DMI header value", err)

	banner.Print("rebuild agent")

	fmt.Println(
		"\n",
		color.GreenString("Agent v%s", helpers.AgentVersion),
		"-",
		color.YellowString("Initrd v%s", initrdVer),
		"-",
		color.MagentaString("Kernel v%s", kernelVer),
		"\n\n\n",
		fmt.Sprintf("%-15s: %s\n", "System Type", systemType),
		fmt.Sprintf("%-15s: %s\n", "System Vendor", systemVendor),
		fmt.Sprintf("%-15s: %s\n", "System Serial", systemSerial),
		fmt.Sprintf("%-15s: %s\n", "System UUID", systemUuid),
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
