package collector

import (
	"fmt"
	"log/slog"
	"math/rand"
	"os"

	"github.com/mizuchilabs/mantrae/pkg/util"
)

type Machine struct {
	Hostname  string
	PrivateIP string
	PublicIPs util.IPAddresses
}

var (
	adjectives = []string{"fuzzy", "brave", "silent", "sleepy", "noisy", "happy", "angry"}
	nouns      = []string{"panda", "ninja", "otter", "dragon", "unicorn", "robot", "android"}
)

// GetMachineInfo retrieves information about the local machine
func GetMachineInfo() *Machine {
	var result Machine
	var err error

	result.Hostname, err = os.Hostname()
	if err != nil {
		result.Hostname = randomName()
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

func randomName() string {
	adj := adjectives[rand.Intn(len(adjectives))]
	noun := nouns[rand.Intn(len(nouns))]
	return fmt.Sprintf("%s-%s", adj, noun)
}
