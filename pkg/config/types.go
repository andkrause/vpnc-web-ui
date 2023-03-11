package config

type Config struct {
	ConnectCommand       string `json:"connectCommand"`
	DisconnectCommand    string `json:"disconnectCommand"`
	ConfigFolder         string `json:"configFolder"`
	WaitTimeAfterConnect int    `json:"waitTimeAfterConnect"`
	ServerPort           int    `json:"serverPort"`
	IPEchoURL            string `json:"ipEchoURL"`
	MaxAgePublicIp       string `json:"maxAgePublicIp"`
}
