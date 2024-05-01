package vpnc

import (
	"fmt"
	"net"
	"strings"
	"time"

	"os"
	"os/exec"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

const _MAX_DISCONNECT_RETRIES = 4

func New(connectCommand, disconnectCommand,
	configFolder, vpncNetworkInterfaceName string) *VPNC {

	return &VPNC{
		connectCommand:           connectCommand,
		disconnectCommand:        disconnectCommand,
		vpncNetworkInterfaceName: vpncNetworkInterfaceName,
		configFolder:             configFolder,
	}
}

func (v *VPNC) ConfigurationExists(configurationName string) bool {
	if _, err := os.Stat(filepath.Join(v.configFolder, fmt.Sprintf("%s.conf", configurationName))); err != nil {
		return false
	}
	return true
}

func (v *VPNC) ConfigurationList() ([]string, error) {
	files, err := os.ReadDir(v.configFolder)
	if err != nil {
		log.Errorf("error reading vpnc configurations: %s", err.Error())
		return []string{}, fmt.Errorf("error reading vpnc configurations: %s", err.Error())
	}

	configurations := []string{}

	for i := range files {
		if files[i].IsDir() {
			continue
		}
		filename := files[i].Name()

		if filename == "default.conf" {
			continue
		}

		extension := filepath.Ext(filename)
		if extension != ".conf" {
			continue
		}

		name := filename[0 : len(filename)-len(extension)]
		configurations = append(configurations, name)
	}

	return configurations, nil
}

func (v *VPNC) Connect(vpncConfig string) (string, error) {
	v.mu.Lock()
	defer v.mu.Unlock()
	if v.isVpncConnected() {
		return "", fmt.Errorf("error executing vpnc connect command interface %s is in use",
			v.vpncNetworkInterfaceName)
	}
	cmd := exec.Command(v.connectCommand, filepath.Join(v.configFolder, vpncConfig))

	result, err := cmd.Output()
	if err != nil {
		log.Errorf("error executing vpnc connect command (%s %s): %s",
			v.connectCommand, vpncConfig, err.Error())
		return "", fmt.Errorf("error executing vpnc connect command (%s %s): %s",
			v.connectCommand, vpncConfig, err.Error())
	}
	log.Infof("Connect to %s successful", vpncConfig)

	return strings.TrimSuffix(string(result), "\n"), nil
}

func (v *VPNC) isVpncConnected() bool {
	// just an existence check
	if _, err := net.InterfaceByName(v.vpncNetworkInterfaceName); err == nil {
		return true
	}
	return false
}

func (v *VPNC) Disconnect() (string, error) {

	v.mu.Lock()
	defer v.mu.Unlock()

	cmd := exec.Command(v.disconnectCommand)
	result, err := cmd.Output()
	if err != nil {
		log.Errorf("error executing vpnc disconnect command (%s): %s",
			v.disconnectCommand, err.Error())
		return "", fmt.Errorf("error executing vpnc connect command (%s): %s",
			v.disconnectCommand, err.Error())
	}
	log.Info("Disconnect successful")

	for i := 0; i <= _MAX_DISCONNECT_RETRIES; i++ {
		//exponential backoff
		time.Sleep(time.Duration(i^2) * time.Second)
		if !v.isVpncConnected() {
			break
		}

	}

	return strings.TrimSuffix(string(result), "\n"), nil
}
