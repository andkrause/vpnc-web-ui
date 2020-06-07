package vpnc

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
)

type VPNC struct {
	connectCommand       string
	disconnectCommand    string
	configFolder         string
	pidFile              string
	waitTimeAfterConnect int
}

func New(connectCommand string, disconnectCommand string,
	configFolder string, pidFile string, waitTimeAfterConnect int) *VPNC {
	return &VPNC{
		connectCommand:       connectCommand,
		disconnectCommand:    disconnectCommand,
		configFolder:         configFolder,
		pidFile:              pidFile,
		waitTimeAfterConnect: waitTimeAfterConnect,
	}
}

func (v *VPNC) ConfigurationList() ([]string, error) {
	files, err := ioutil.ReadDir(v.configFolder)
	if err != nil {
		log.Errorf("error reading vpnc configurations: %s", err.Error())
		return []string{}, fmt.Errorf("error reading vpnc configurations: %s", err.Error())
	}

	configurations := []string{}

	for i, _ := range files {
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

func (v *VPNC) CheckConnection() bool {
	_, err := os.Stat(v.pidFile)

	if os.IsNotExist(err) {
		return false
	}

	return true
}

func (v *VPNC) Connect(vpncConfig string) (string, error) {

	// Always make sure you are disconnected error can be ignored ;-)
	v.Disconnect()

	cmd := exec.Command(v.connectCommand, vpncConfig)

	result, err := cmd.Output()
	if err != nil {
		log.Errorf("error executing vpnc connect command (%s %s): %s",
			v.connectCommand, vpncConfig, err.Error())
		return "", fmt.Errorf("error executing vpnc connect command (%s %s): %s",
			v.connectCommand, vpncConfig, err.Error())
	}
	log.Infof("Connect to %s successful", vpncConfig)

	time.Sleep(time.Duration(v.waitTimeAfterConnect) * time.Second)
	return string(result), nil
}

func (v *VPNC) Disconnect() (string, error) {

	cmd := exec.Command(v.disconnectCommand)
	result, err := cmd.Output()
	if err != nil {
		log.Errorf("error executing vpnc disconnect command (%s): %s",
			v.connectCommand, err.Error())
		return "", fmt.Errorf("error executing vpnc connect command (%s): %s",
			v.connectCommand, err.Error())
	}
	log.Info("Disconnect successful")
	return string(result), nil
}
