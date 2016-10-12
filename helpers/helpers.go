package helpers

import "io/ioutil"

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

	return string(ver), err
}