package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
)

func ParseConfigFile(location string) (*Config, error) {

	configBytes, err := ioutil.ReadFile(location)
	if err != nil {
		log.Errorf("Error reading config file (from %s): %s", location, err.Error())
		return nil, fmt.Errorf("Error reading config file (from %s): %s", location, err.Error())
	}

	var configParsed Config
	if err := json.Unmarshal(configBytes, &configParsed); err != nil {
		log.Errorf("Error parsing config file (from %s): %s", location, err.Error())
		return nil, fmt.Errorf("Error parsing config file (from %s): %s", location, err.Error())
	}

	return &configParsed, nil

}

func (c *Config) LogConfig() {
	fmt.Printf("VPNC Connect Command: %s\n"+
		"VPNC Disconnect Command: %s\n"+
		"VPNC Config Folder (*.conf): %s\n"+
		"PID File: %s\n"+
		"Wait time after connect (waiting for background job to start): %d seconds\n"+
		"Web UI Port: %d\n"+
		"IP Echo URL: %s\n",
		c.VPNC.ConnectCommand, c.VPNC.ConnectCommand, c.VPNC.ConfigFolder, c.VPNC.PIDFile, c.VPNC.WaitTimeAfterConnect,
		c.WebUI.ServerPort, c.WebUI.IPEchoURL)
}
