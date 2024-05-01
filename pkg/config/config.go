package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

func ParseConfigFile(location string) (*Config, error) {

	configBytes, err := os.ReadFile(location)
	if err != nil {
		log.Errorf("Error reading config file (from %s): %s", location, err.Error())
		return nil, fmt.Errorf("error reading config file (from %s): %s", location, err.Error())
	}

	var configParsed Config
	if err := json.Unmarshal(configBytes, &configParsed); err != nil {
		log.Errorf("Error parsing config file (from %s): %s", location, err.Error())
		return nil, fmt.Errorf("error parsing config file (from %s): %s", location, err.Error())
	}

	// set default for wireguard network interface name (if wireguard is present)

	if configParsed.IsWireguardConfigured() &&
		len(configParsed.Wireguard.WireguardNetworkInterfaceName) == 0 {

		configParsed.Wireguard.WireguardNetworkInterfaceName = "wg0"
	}

	if configParsed.IsWireguardConfigured() &&
		len(configParsed.Wireguard.WgQuickConfigSearchDir) == 0 {

		configParsed.Wireguard.WgQuickConfigSearchDir = "/etc/wireguard/"
	}

	if configParsed.IsVpncConfigured() &&
		len(configParsed.Vpnc.VpncNetworkInterfaceName) == 0 {

		configParsed.Wireguard.WgQuickConfigSearchDir = "tun0"
	}

	return &configParsed, nil

}

func (c *Config) LogConfig() {

	if c.IsWireguardConfigured() {
		fmt.Printf("Wireguard Command Location: %s\n"+
			"Wireguard Config Folder (*.conf): %s\n"+
			"Wireguard Network Inteface Name: %s\n"+
			"wg-quick config search dir (depends on OS): %s\n",
			c.Wireguard.WgQuickCommand, c.Wireguard.ConfigFolder,
			c.Wireguard.WireguardNetworkInterfaceName,
			c.Wireguard.WgQuickConfigSearchDir)
	} else {
		fmt.Println("Wireguard is not configured")
	}

	if c.IsVpncConfigured() {
		fmt.Printf("VPNC Connect Command: %s\n"+
			"VPNC Disconnect Command: %s\n"+
			"VPNC Config Folder (*.conf): %s\n"+
			"VPNC Network Inteface Name: %s\n",
			c.Vpnc.ConnectCommand, c.Vpnc.DisconnectCommand, c.Vpnc.ConfigFolder,
			c.Vpnc.VpncNetworkInterfaceName)
	} else {
		fmt.Println("VPNC is not configured")
	}

	fmt.Printf("Wait time after connect (waiting for background job to start): %d seconds\n"+
		"Web UI Port: %d\n"+
		"IP Echo URL: %s\n"+
		"Max age of IP: %s\n",
		c.WaitTimeAfterConnect, c.ServerPort, c.IPEchoURL, c.MaxAgePublicIp)
}

func (c *Config) GetMaxAgePublicIpDuration() time.Duration {
	duration, err := time.ParseDuration(c.MaxAgePublicIp)
	if err != nil {
		return 10 * time.Minute
	}
	return duration
}

// validations, better be defensive if there is no testing ;-)

func (c *Config) IsWireguardConfigured() bool {
	return c.Wireguard != nil
}

func (c *Config) IsVpncConfigured() bool {
	return c.Vpnc != nil
}

func (c *Config) Validate() error {
	// either wireguard or vpnc need to be configured
	if !(c.IsWireguardConfigured() || c.IsVpncConfigured()) {
		return fmt.Errorf("error: neither vpnc nor wireguard are configured")
	}

	if c.IsVpncConfigured() {
		if err := c.Vpnc.Validate(); err != nil {
			return fmt.Errorf("vpnc configured incorrectly: %s", err)
		}
	}

	if c.IsWireguardConfigured() {
		if err := c.Wireguard.Validate(); err != nil {
			return fmt.Errorf("wireguard configured incorrectly: %s", err)
		}
	}

	if c.ServerPort == 0 {
		return fmt.Errorf("error: serverPort needs to be configured")
	}
	if len(c.IPEchoURL) == 0 {
		return fmt.Errorf("error: ipEchoURL needs to be configured")
	}
	if _, err := time.ParseDuration(c.MaxAgePublicIp); err != nil {
		return fmt.Errorf("error: maxAgePublicIp wrongly configured: %s", err)
	}

	return nil
}

func configDirectoryValid(configDirectory string) error {

	if len(configDirectory) == 0 {
		return fmt.Errorf("directory can't be an empty string")
	}

	if confDir, err := os.Stat(configDirectory); err != nil || !confDir.IsDir() {
		return fmt.Errorf("error: %s is not a directory or does not exist: %s",
			configDirectory, err)
	}

	return nil
}

func (c *VpncClientConfig) Validate() error {
	if len(c.ConnectCommand) == 0 {
		return fmt.Errorf("error: connectCommand needs to be configured")
	}
	if len(c.DisconnectCommand) == 0 {
		return fmt.Errorf("error: disconnectCommand needs to be configured")
	}
	if err := configDirectoryValid(c.ConfigFolder); err != nil {
		return fmt.Errorf("error: configFolder needs to be configured correctly: %s", err)
	}
	return nil
}

func (c *WgQuickClientConfig) Validate() error {

	if len(c.WgQuickCommand) == 0 {
		return fmt.Errorf("error: wgQuickCommand needs to be configured")
	}
	if err := configDirectoryValid(c.ConfigFolder); err != nil {
		return fmt.Errorf("error: configFolder needs to be configured correctly: %s", err)
	}

	return nil
}
