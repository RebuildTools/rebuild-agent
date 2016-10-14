package helpers
// helpers is a package of functions
// designed to make the development
// process for the Rebuild Agent a
// little bit easier by providing
// commonly used functions and variables
// in a common package

import (
	"os"
	"path"
	"fmt"
	"io/ioutil"
	"golang.org/x/sys/unix"
	"strings"
	"github.com/Sirupsen/logrus"
	"strconv"
)

// AgentVersion defines the build of
// the agent along with a short
// git hash of the point in the
// repo when the agent as built
var AgentVersion string = "UNSET"

// initrdReleaseFile defines the path
// to the file containing the Rebuild
// Initrd release version string
const initrdReleaseFile = "/lib/rebuild/initrd-release"

// dmiPath defines the location
// of file based DMI headers and
// values, these files are pseudo files
// presented by the Linux Kernel with
// values pulled from the DMI
const dmiPath = "/sys/class/dmi/id"

// HandleError provides a simple
// function to check if a error is
// empty and if not to log the error
// and the current state when checking
// the error. This function is purely
// to save space in code.
func HandleError(logger *logrus.Logger, state string, err error) {
	if err != nil {
		logger.WithError(err).WithField("state", state).Fatal("An error was caught!")
	}
}

// GetInitrdVersion reads a file from
// the system that contains the version
// of the Rebuild Initrd this agent is
// running from
func GetInitrdVersion() (string, error) {
	ver, err := ioutil.ReadFile(initrdReleaseFile)

	return strings.TrimSpace(string(ver)), err
}

// GetKernelVersion queries uname
// and returns the "Release" field
// as a string
func GetKernelVersion() (string, error){
	uname := unix.Utsname{}
	err := unix.Uname(&uname)

	if err != nil {
		return "", nil
	}

	releaseBytes := make([]byte, len(uname.Release))
	for i, v := range uname.Release {
		releaseBytes[i] = byte(v)
	}

	return string(releaseBytes), nil
}

// GetDMIValue reads a DMI header
// value from a Linux Kernel pseudo file
func GetDMIValue(header string) (string, error) {
	dmiHeaderFile := path.Join(dmiPath, header)

	if _, err := os.Stat(dmiHeaderFile); os.IsNotExist(err) {
		return "", fmt.Errorf("No such DMI header [%s] could be found", header)
	}

	headerValue, err := ioutil.ReadFile(dmiHeaderFile)

	return strings.TrimSpace(string(headerValue)), err
}

// stringToInt32 converts a string
// into a 32bit integer
func StringToInt32(value string) int32 {
	val, err := strconv.ParseInt(value, 10, 32)

	if err != nil { panic(err) }

	return int32(val)
}

// stringToInt64 converts a string
// into a 64bit integer
func StringToInt64(value string) int64 {
	val, err := strconv.ParseInt(value, 10, 64)

	if err != nil { panic(err) }

	return int64(val)
}