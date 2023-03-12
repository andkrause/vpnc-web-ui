package vpnc

import (
	"sync"
	"time"
)

type VPNC struct {
	connectCommand       string
	disconnectCommand    string
	configFolder         string
	waitTimeAfterConnect int
	lastUpdatePublicIp   time.Time
	currentPublicIp      string
	activeVpnConfig      string
	message              string
	ipEchoUrl            string
	maxAgePublicIp       time.Duration
	mu                   sync.Mutex
}

type VpnStatus struct {
	CurrentPublicIp string
	ActiveVpnConfig string
	Message         string
}
