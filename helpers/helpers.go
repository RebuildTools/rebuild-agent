package helpers

import (
	"io/ioutil"
	"golang.org/x/sys/unix"
	"strings"
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