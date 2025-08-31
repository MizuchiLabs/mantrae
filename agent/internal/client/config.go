package client

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mizuchilabs/mantrae/pkg/util"
)

type Config struct {
	Token               string
	ServerURL           string
	ProfileID           int64
	AgentID             string
	ActiveIP            string
	HealthCheckInterval time.Duration
	UpdateInterval      time.Duration
	ConnectionTimeout   time.Duration
	HealthTimeout       time.Duration
}

func Load() (*Config, error) {
	token := os.Getenv("TOKEN")
	host := os.Getenv("HOST")
	if token == "" || host == "" {
		return nil, errors.New("TOKEN and HOST must be set")
	}

	profileID, agentID, err := parseToken(token)
	if err != nil {
		return nil, err
	}

	return &Config{
		Token:               token,
		ServerURL:           util.CleanURL(host),
		ProfileID:           profileID,
		AgentID:             agentID,
		ActiveIP:            "",
		HealthCheckInterval: 15 * time.Second,
		UpdateInterval:      10 * time.Second,
		ConnectionTimeout:   10 * time.Second,
		HealthTimeout:       5 * time.Second,
	}, nil
}

func parseToken(token string) (int64, string, error) {
	parts := strings.Split(token, ".")
	if len(parts) < 2 {
		return 0, "", errors.New("invalid token format")
	}

	profileID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return 0, "", errors.New("invalid profile ID in token")
	}

	return profileID, parts[1], nil
}
