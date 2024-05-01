package vpnclient

import (
	"sync"
	"time"
)

type VpnClient interface {
	ConfigurationExists(configurationName string) bool
	ConfigurationList() ([]string, error)
	Connect(vpncConfig string) (string, error)
	Disconnect() (string, error)
}

type VpnClientAggregator struct {
	clients              map[string]VpnClient
	waitTimeAfterConnect int
	lastUpdatePublicIp   time.Time
	currentPublicIp      string
	activeVpnClient      string
	activeVpnConfig      string
	message              string
	ipEchoUrl            string
	maxAgePublicIp       time.Duration
	mu                   sync.Mutex
}

type VpnClientContainer struct {
	Client VpnClient
	Name   string
}

type VpnConfigurations struct {
	VpnClientName       string
	AvailableVpnConfigs []string
}

type VpnStatus struct {
	CurrentPublicIp string
	ActiveVpnClient string
	ActiveVpnConfig string
	Message         string
}
