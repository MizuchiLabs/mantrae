package collector

import (
	"log/slog"
	"os"

	"github.com/mizuchilabs/mantrae/pkg/util"
)

type Machine struct {
	MachineID string
	Hostname  string
	PrivateIP string
	PublicIPs util.IPAddresses
}

// GetMachineInfo retrieves information about the local machine
func GetMachineInfo() *Machine {
	var result Machine
	var err error

	id, err := os.ReadFile("/etc/machine-id")
	if err != nil {
		id, err = os.ReadFile("/var/lib/dbus/machine-id")
		if err != nil {
			slog.Error("Failed to read machine ID", "error", err)
		}
	}
	if len(id) > 0 {
		result.MachineID = string(id)
	}

	result.Hostname, err = os.Hostname()
	if err != nil {
		result.Hostname = "unknown"
		slog.Error("Failed to get hostname", "error", err)
	}

	result.PrivateIP, err = util.GetHostIPv4()
	if err != nil {
		slog.Error("Failed to get local IP", "error", err)
	}

	result.PublicIPs, err = util.GetPublicIPs()
	if err != nil {
		slog.Error("Failed to get public IP", "error", err)
	}

	return &result
}
