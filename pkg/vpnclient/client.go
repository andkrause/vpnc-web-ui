package vpnclient

import (
	"fmt"
	"io"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

const _MAX_IP_RETRIES = 6

func New(waitTimeAfterConnect int,
	ipEchoUrl string, maxAgePublicIp time.Duration, clients ...*VpnClientContainer) *VpnClientAggregator {

	ip, err := getPublicIp(ipEchoUrl)
	var ipTime time.Time
	if err != nil {
		ipTime = time.Unix(0, 0)
	} else {
		ipTime = time.Now()
	}
	aggregator := VpnClientAggregator{
		waitTimeAfterConnect: waitTimeAfterConnect,
		activeVpnClient:      "",
		activeVpnConfig:      "",
		message:              "",
		currentPublicIp:      ip,
		lastUpdatePublicIp:   ipTime,
		ipEchoUrl:            ipEchoUrl,
		maxAgePublicIp:       maxAgePublicIp,
		clients:              map[string]VpnClient{},
	}
	for _, client := range clients {
		aggregator.clients[client.Name] = client.Client
	}

	return &aggregator
}

func (v *VpnClientAggregator) Status() *VpnStatus {

	// check if IP needs to be updated
	if v.lastUpdatePublicIp.Add(v.maxAgePublicIp).Before(time.Now()) {
		if ip, err := getPublicIp(v.ipEchoUrl); err == nil {
			v.currentPublicIp, v.lastUpdatePublicIp = ip, time.Now()
		}

	}

	return &VpnStatus{
		CurrentPublicIp: v.currentPublicIp,
		ActiveVpnClient: v.activeVpnClient,
		ActiveVpnConfig: v.activeVpnConfig,
		Message:         v.message,
	}
}

func (v *VpnClientAggregator) ConfigurationExists(clientName string, configurationName string) bool {

	if client, ok := v.clients[clientName]; ok {
		return client.ConfigurationExists(configurationName)
	}

	return false
}

func (v *VpnClientAggregator) ConfigurationList() ([]VpnConfigurations, error) {
	configs := make([]VpnConfigurations, len(v.clients))

	i := 0

	var err error
	for clientName, client := range v.clients {
		configs[i].VpnClientName = clientName
		configs[i].AvailableVpnConfigs, err = client.ConfigurationList()
		if err != nil {
			log.Errorf("error reading vpn configs for client %s: %s", clientName, err.Error())
			return nil, fmt.Errorf("error reading vpn configs for client %s: %s", clientName, err.Error())
		}
		i++
	}

	return configs, nil

}

func (v *VpnClientAggregator) Connect(clientName string, configurationName string) error {
	// Always make sure you are disconnected error can be ignored ;-)
	v.Disconnect()

	v.mu.Lock()
	defer v.mu.Unlock()

	client, ok := v.clients[clientName]
	if !ok {
		log.Errorf("vpn client %s not available", clientName)
		return fmt.Errorf("vpn client %s not available", clientName)
	}

	message, err := client.Connect(configurationName)
	if err != nil {
		log.Errorf("error conncting to vpn with client %s: %s", clientName, err.Error())
		return fmt.Errorf("error conncting to vpn with client %s: %s", clientName, err.Error())
	}

	v.activeVpnClient = clientName
	v.activeVpnConfig = configurationName
	v.message = message

	//Invalidate previous IP
	v.resetIp()

	time.Sleep(time.Duration(v.waitTimeAfterConnect) * time.Second)
	return nil
}

func (v *VpnClientAggregator) Disconnect() (string, error) {
	v.mu.Lock()
	defer v.mu.Unlock()

	var returnMessage string
	var returnError error

	//always disconnect everything
	for clientName, client := range v.clients {
		message, err := client.Disconnect()

		// don't ignore errors and messages when config is active
		if clientName == v.activeVpnClient {

			if err != nil {
				returnError = fmt.Errorf("error disconnecting currently active configuration from client %s: %s",
					v.activeVpnClient, err)
				log.Error(returnError)
			}
			returnMessage = message

			v.activeVpnClient = ""
			v.activeVpnConfig = ""
		}
	}

	//Invalidate previous IP
	v.resetIp()

	return returnMessage, returnError
}

func (v *VpnClientAggregator) resetIp() {
	v.lastUpdatePublicIp, v.currentPublicIp = time.Unix(0, 0), ""
}

func getPublicIp(ipEchoUrl string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, ipEchoUrl, nil)
	if err != nil {
		log.Errorf("error getting IP Echo: %s", err)
		return "", fmt.Errorf("error getting IP Echo: %s", err)
	}

	var errorBackup error = nil

	for i := 0; i <= _MAX_IP_RETRIES; i++ {
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
