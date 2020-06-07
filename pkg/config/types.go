package config

type Config struct {
	VPNC struct {
		ConnectCommand    string `json:"connectCommand"`
		DisconnectCommand string `json:"disconnectCommand"`
		ConfigFolder      string `json:"configFolder"`
		PIDFile           string `json:"pidFile"`
	} `json:"vpnc"`
	WebUI struct {
		ServerPort int    `json:"serverPort"`
		IPEchoURL  string `json:"ipEchoURL"`
	} `json:"webUI"`
}
