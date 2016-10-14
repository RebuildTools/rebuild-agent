package profiler
// profiler is a command for
// the Rebuild Agent that
// retrieves system information,
// builds a profile and sends the
// profile back to the Rebuild Core

import (
	"github.com/Sirupsen/logrus"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/net"
	protocol "github.com/RebuildTools/rebuild-agent-protocol/golang"
	"github.com/RebuildTools/rebuild-agent/helpers"
	"fmt"
	"encoding/json"
	"os/exec"
)

// Run is called by the main of
// the agent to run this command.
func Run(logger *logrus.Logger) {
	logger.Info("Collecting system information to build a profile")

	systemType, err := helpers.GetDMIValue("product_name")
	helpers.HandleError(logger, "Retrieving DMI header value", err)

	systemSerial, err := helpers.GetDMIValue("product_serial")
	helpers.HandleError(logger, "Retrieving DMI header value", err)

	systemVendor, err := helpers.GetDMIValue("chassis_vendor")
	helpers.HandleError(logger, "Retrieving DMI header value", err)

	systemUuid, err := helpers.GetDMIValue("product_uuid")
	helpers.HandleError(logger, "Retrieving DMI header value", err)

	systemProfile := protocol.System{
		SystemSerial:		systemSerial,
		SystemUuid:		systemUuid,
		SystemVendor:		systemVendor,
		SystemProductName: 	systemType,
	}

	profileSystemMemory(&systemProfile, logger)
	profileSystemProcessors(&systemProfile, logger)
	profileNetworkInterfaces(&systemProfile, logger)
	profileSystemStorage(&systemProfile, logger)

	// TODO(lh): Marshal the profile to bytes and send it via TCP to the core
	out, _ := json.Marshal(systemProfile)
	fmt.Printf("%s\n", out)
}

// profileSystemMemory queries the system
// for its total RAM in bytes
func profileSystemMemory(sysProfile *protocol.System, logger *logrus.Logger) {
	logger.Debug("Profilling System Memory")

	virtMem, err := mem.VirtualMemory()
	helpers.HandleError(logger, "Collecting system memory information", err)

	sysProfile.SystemMemoryTotal = int64(virtMem.Total)
}

// profileSystemProcessors queries the system
// for every available CPU and profiles each one
func profileSystemProcessors(sysProfile *protocol.System, logger *logrus.Logger) {
	logger.Debug("Profilling System Processors")

	sysCpus, err := cpu.Info()
	helpers.HandleError(logger, "Collecting system proccessor information", err)

	filteredCpus := map[int32]*cpu.InfoStat{}

	for _, sysCpu := range sysCpus {
		if sysCpu.PhysicalID == "" {
			sysCpu.PhysicalID = "0"
		}

		if filteredCpus[helpers.StringToInt32(sysCpu.PhysicalID)] != nil {
			continue
		}

		filteredCpus[helpers.StringToInt32(sysCpu.PhysicalID)] = &sysCpu
	}

	cpuSockets := []*protocol.CpuSocket{}
	for sysCpuSocket, sysCpu := range filteredCpus {
		cpuSockets = append(cpuSockets, &protocol.CpuSocket{
			SocketNumber: 	sysCpuSocket,
			TotalCores: 	sysCpu.Cores,
			VendorId: 	sysCpu.VendorID,
			FamilyId: 	helpers.StringToInt32(sysCpu.Family),
			ModelId: 	helpers.StringToInt32(sysCpu.Model),
			ModelName: 	sysCpu.ModelName,
			Mhz:  		sysCpu.Mhz,
			CahceSize: 	sysCpu.CacheSize,
			Flags: 		sysCpu.Flags,
		})
	}

	sysProfile.SystemSockets = cpuSockets
}

// stringSliceContains is a simple "contains"
// check of string slices
func stringSliceContains(slice []string, element string) bool {
	for _, sliceElement := range slice {
		if sliceElement == element {
			return true
		}
	}

	return false
}

// profileNetworkInterfaces queries the system
// for network interfaces nad profiles them, excluding
// the initrd loopback interface
func profileNetworkInterfaces(sysProfile *protocol.System, logger *logrus.Logger) {
	logger.Debug("Profilling System Network Interfaces")

	netInts, err := net.Interfaces()
	helpers.HandleError(logger, "Collecting system network interfaces information", err)

	interfaces := []*protocol.NetworkInterface{}

	for _, netInt := range netInts {
		if stringSliceContains(netInt.Flags, "loopback") {
			continue
		}

		interfaces = append(interfaces, &protocol.NetworkInterface{
			Name: 		netInt.Name,
			HardwareAddr: 	netInt.HardwareAddr,
		})
	}

	sysProfile.SystemNetworkInterfaces = interfaces
}

// profileSystemStorage calls the Linux tool
// "lsblk" to list the block devices know to
// the system, it takes this information and
// converts it to message for Rebuild and adds
// it to the system profile
func profileSystemStorage(sysProfile *protocol.System, logger *logrus.Logger) {
	logger.Debug("Profilling System Storage")

	lsblkPath, err := exec.LookPath("lsblk")
	helpers.HandleError(logger, "Searching for lsblk executable", err)

	lsblkOutput, err := exec.Command(lsblkPath, "--output", "NAME,SIZE,TYPE,MODEL,SERIAL,VENDOR,REV", "-J", "-d", "-b").Output()
	helpers.HandleError(logger, "Collecting block device information", err)

	unmarshaledBlockDevices := struct {
		Devices	[]struct{
			Name	string `json:"name"`
			Size	string `json:"size"`
			Type	string `json:"type"`
			Model	string `json:"model"`
			Serial	string `json:"serial"`
			Vendor	string `json:"vendor"`
			Rev	string `json:"rev"`
		} `json:"blockdevices"`
	}{}

	err = json.Unmarshal(lsblkOutput, &unmarshaledBlockDevices)
	helpers.HandleError(logger, "Unmarshalling block device information", err)

	blockDevices := []*protocol.BlockDevice{}
	for _, device := range unmarshaledBlockDevices.Devices {
		blockDevices = append(blockDevices, &protocol.BlockDevice{
			Name:		device.Name,
			Size:		helpers.StringToInt64(device.Size),
			Type:		device.Type,
			Model:		device.Model,
			Serial:		device.Serial,
			Vendor:		device.Vendor,
			Revision:	device.Rev,
		})
	}

	sysProfile.SystemStorage = blockDevices
}