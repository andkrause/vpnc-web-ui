package vpnc

import (
	"fmt"
	"io"

	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
)

const MAX_IP_RETRIES = 4

func New(connectCommand string, disconnectCommand string,
	configFolder string, waitTimeAfterConnect int,
	ipEchoUrl string, maxAgePublicIp time.Duration) *VPNC {

	ip, err := getPublicIp(ipEchoUrl)
	var ipTime time.Time
	if err != nil {
		ipTime = time.Unix(0, 0)
	} else {
		ipTime = time.Now()
	}
	return &VPNC{
		connectCommand:       connectCommand,
		disconnectCommand:    disconnectCommand,
		configFolder:         configFolder,
		waitTimeAfterConnect: waitTimeAfterConnect,
		activeVpnConfig:      "",
		messages:             "",
		currentPublicIp:      ip,
		lastUpdatePublicIp:   ipTime,
		ipEchoUrl:            ipEchoUrl,
		maxAgePublicIp:       maxAgePublicIp,
	}
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

func (v *VPNC) Connect(vpncConfig string) error {

	// Always make sure you are disconnected error can be ignored ;-)
	v.Disconnect()

	cmd := exec.Command(v.connectCommand, fmt.Sprintf("%s%s", v.configFolder, vpncConfig))

	result, err := cmd.Output()
	if err != nil {
		log.Errorf("error executing vpnc connect command (%s %s): %s",
			v.connectCommand, vpncConfig, err.Error())
		return fmt.Errorf("error executing vpnc connect command (%s %s): %s",
			v.connectCommand, vpncConfig, err.Error())
	}
	log.Infof("Connect to %s successful", vpncConfig)

	v.activeVpnConfig = vpncConfig
	v.messages = string(result)
	//Invalidate previous IP
	v.resetIp()

	time.Sleep(time.Duration(v.waitTimeAfterConnect) * time.Second)
	return nil
}

func (v *VPNC) Disconnect() error {

	cmd := exec.Command(v.disconnectCommand)
	result, err := cmd.Output()
	if err != nil {
		log.Errorf("error executing vpnc disconnect command (%s): %s",
			v.connectCommand, err.Error())
		return fmt.Errorf("error executing vpnc connect command (%s): %s",
			v.connectCommand, err.Error())
	}
	log.Info("Disconnect successful")
	// Invalidate previous IP
	v.resetIp()
	v.activeVpnConfig = ""
	v.messages = string(result)

	return nil
}

func (v *VPNC) Status() *VpnStatus {

	// check if IP needs to be updated
	if v.lastUpdatePublicIp.Add(v.maxAgePublicIp).Before(time.Now()) {
		if ip, err := getPublicIp(v.ipEchoUrl); err == nil {
			v.currentPublicIp, v.lastUpdatePublicIp = ip, time.Now()
		}

	}

	return &VpnStatus{
		CurrentPublicIp: v.currentPublicIp,
		ActiveVpnConfig: v.activeVpnConfig,
		Messages:        v.messages,
	}
}

func (v *VPNC) resetIp() {
	v.lastUpdatePublicIp, v.currentPublicIp = time.Unix(0, 0), ""
}

func getPublicIp(ipEchoUrl string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, ipEchoUrl, nil)
	if err != nil {
		log.Errorf("error getting IP Echo: %s", err)
		return "", fmt.Errorf("error getting IP Echo: %s", err)
	}

	var errorBackup error = nil

	for i := 0; i <= MAX_IP_RETRIES; i++ {
		//exponential backoff
		time.Sleep(time.Duration(i^2) * time.Second)
		client := http.Client{
			Timeout: 1 * time.Second,
			Transport: &http.Transport{
				DisableKeepAlives: true,
			},
		}

		resp, err := client.Do(req)
		if err != nil {
			log.Errorf("error getting IP Echo: %s", err)
			errorBackup = fmt.Errorf("error getting IP Echo: %s", err)
			continue
		}
		respBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Errorf("error reading IP Echo: %s", err)
			errorBackup = fmt.Errorf("error reading IP Echo: %s", err)
			continue
		}
		defer resp.Body.Close()

		return string(respBytes), nil
	}

	return "", errorBackup
}
