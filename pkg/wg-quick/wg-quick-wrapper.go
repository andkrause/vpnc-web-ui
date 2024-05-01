package wgquick

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
)

const _CONFIG_FILE_EXTENSION_LENGTH = 5 //including separator
const _CONFIG_FILE_EXTENSION = "conf"
const _CONFIG_FILE_REGULAR_EXPRESSION = "[a-zA-Z0-9_=+.-]{1,15}[.]" + _CONFIG_FILE_EXTENSION

func New(wgQuickCommand, configFolder, wireguardNetworkInterfaceName,
	wgQuickConfigSearchDir string) *WgQuick {
	return &WgQuick{
		wgQuickCommand:                wgQuickCommand,
		configFolder:                  configFolder,
		wireguardNetworkInterfaceName: wireguardNetworkInterfaceName,
		wgQuickConfigSearchDir:        wgQuickConfigSearchDir,
		mu:                            sync.Mutex{},
		configFolderCheckRegexp:       regexp.MustCompile(_CONFIG_FILE_REGULAR_EXPRESSION),
	}
}

func (v *WgQuick) ConfigurationExists(configurationName string) bool {
	if _, err := os.Stat(filepath.Join(v.configFolder, fmt.Sprintf("%s.conf", configurationName))); err != nil {
		return false
	}
	return true
}

func (v *WgQuick) validWireguardConfigCheck(name string) bool {
	return v.configFolderCheckRegexp.MatchString(name)
}

func (v *WgQuick) readConfigFiles() ([]string, error) {
	files, err := os.ReadDir(v.configFolder)
	if err != nil {
		log.Errorf("error reading wg-quick configurations: %s", err.Error())
		return []string{}, fmt.Errorf("error reading wg-quick configurations: %s", err.Error())
	}

	configurations := []string{}
	for i := range files {
		if files[i].IsDir() {
			continue
		}

		if filename := files[i].Name(); v.validWireguardConfigCheck(filename) {
			name := filename[0 : len(filename)-_CONFIG_FILE_EXTENSION_LENGTH]
			configurations = append(configurations, name)
		}

	}

	return configurations, nil

}

func (v *WgQuick) ConfigurationList() ([]string, error) {
	return v.readConfigFiles()
}

func checkAndRemoveSymlink(fileName string) {
	if _, err := os.Lstat(fileName); err == nil {
		os.Remove(fileName)
	}
}

func (v *WgQuick) Connect(wgQuickConfig string) (string, error) {
	// Always make sure you are disconnected error can be ignored ;-)
	v.Disconnect()

	v.mu.Lock()
	defer v.mu.Unlock()

	referencedConfig := filepath.Join(v.configFolder, fmt.Sprintf("%s.%s", wgQuickConfig, _CONFIG_FILE_EXTENSION))
	targetSymLink := filepath.Join(v.wgQuickConfigSearchDir, fmt.Sprintf("%s.%s", v.wireguardNetworkInterfaceName, _CONFIG_FILE_EXTENSION))

	checkAndRemoveSymlink(targetSymLink)

	if err := os.Symlink(referencedConfig, targetSymLink); err != nil {
		return "", fmt.Errorf("error while connecting, creating link to referenced configuration failed: %s", err)
	}

	log.Info(v.wgQuickCommand, fmt.Sprintf("up %s", v.wireguardNetworkInterfaceName))

	cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("%s up %s", v.wgQuickCommand, v.wireguardNetworkInterfaceName))
	result, err := cmd.Output()

	if err != nil {

		errorMessage := err.Error()

		if exiterr, ok := err.(*exec.ExitError); ok {
			errorMessage = strings.TrimSuffix(string(exiterr.Stderr), "\n")
		}

		log.Errorf("error executing wg-quick connect command for config %q: %s",
			wgQuickConfig, errorMessage)
		return "", fmt.Errorf("error executing wg-quick connect command for config %q: %s",
			wgQuickConfig, errorMessage)
	}

	log.Infof("Connect to %s successful", wgQuickConfig)
	return strings.TrimSuffix(string(result), "\n"), nil
}

func (v *WgQuick) Disconnect() (string, error) {
	v.mu.Lock()
	defer v.mu.Unlock()

	if _, err := net.InterfaceByName(v.wireguardNetworkInterfaceName); err != nil {
		return "wireguard was not connected", nil
	}

	cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("%s down %s", v.wgQuickCommand, v.wireguardNetworkInterfaceName))

	result, err := cmd.Output()
	if err != nil {
		errorMessage := err.Error()

		if exiterr, ok := err.(*exec.ExitError); ok {
			errorMessage = strings.TrimSuffix(string(exiterr.Stderr), "\n")
		}
		log.Errorf("error executing wg-quick down command: %s", errorMessage)
		return "", fmt.Errorf("error executing wg-quick down command: %s", errorMessage)
	}
	log.Info("Disconnect successful")

	return strings.TrimSuffix(string(result), "\n"), nil
}
