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

	cpus, err := cpu.Info()
	helpers.HandleError(logger, "Collecting system proccessor information", err)

	cpuSockets := []*protocol.CpuSocket{}

	sock := &protocol.CpuSocket{}
	for _, cpu := range cpus {
		if cpu.PhysicalID == "" {
			cpu.PhysicalID = "0"
		}

		if helpers.StringToInt32(cpu.PhysicalID) != sock.SocketNumber {
			cpuSockets = append(cpuSockets, sock)
		}

		sock = &protocol.CpuSocket{
			SocketNumber: 	helpers.StringToInt32(cpu.PhysicalID),
			TotalCores: 	cpu.Cores,
			VendorId: 	cpu.VendorID,
			FamilyId: 	helpers.StringToInt32(cpu.Family),
			ModelId: 	helpers.StringToInt32(cpu.Model),
			ModelName: 	cpu.ModelName,
			Mhz:  		cpu.Mhz,
			CahceSize: 	cpu.CacheSize,
			Flags: 		cpu.Flags,
		}
	}

	cpuSockets = append(cpuSockets, sock)

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