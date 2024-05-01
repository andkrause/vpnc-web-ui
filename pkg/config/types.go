package config

type Config struct {
	Wireguard            *WgQuickClientConfig `json:"wireguard"`
	Vpnc                 *VpncClientConfig    `json:"vpnc"`
	WaitTimeAfterConnect int                  `json:"waitTimeAfterConnect"`
	ServerPort           int                  `json:"serverPort"`
	IPEchoURL            string               `json:"ipEchoURL"`
	MaxAgePublicIp       string               `json:"maxAgePublicIp"`
}

type VpncClientConfig struct {
	ConnectCommand           string `json:"connectCommand"`
	DisconnectCommand        string `json:"disconnectCommand"`
	ConfigFolder             string `json:"configFolder"`
	VpncNetworkInterfaceName string `json:"vpncNetworkInterfaceName"`
}

type WgQuickClientConfig struct {
	WgQuickCommand                string `json:"wgQuickCommand"`
	ConfigFolder                  string `json:"configFolder"`
	WireguardNetworkInterfaceName string `json:"wireguardNetworkInterfaceName"`
	WgQuickConfigSearchDir        string `json:"wgQuickConfigSearchDir"`
}
