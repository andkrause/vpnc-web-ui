package vpnc

import "time"

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
}

type VpnStatus struct {
	CurrentPublicIp string `json:"currentPublicIp"`
	ActiveVpnConfig string `json:"activeVpnConfig"`
	Message         string `json:"message"`
}
