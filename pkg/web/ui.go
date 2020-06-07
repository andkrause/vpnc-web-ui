package web

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"html/template"

	"github.com/andkrause/vpnc-web-ui/pkg/vpnc"
	log "github.com/sirupsen/logrus"
)

const TEMPLATE_LOCATION = "./templates/index.html"
const MAX_IP_RETRIES = 4

type uiData struct {
	ConfigurationList []string
	CommandResults    []string
	OwnPublicIP       string
	ConnectionState   string
}

type UI struct {
	template  *template.Template
	vpnc      *vpnc.VPNC
	ipEchoUrl string
}

func New(vpnc *vpnc.VPNC, ipEchoUrl string) (*UI, error) {
	tmpl, err := template.ParseFiles(TEMPLATE_LOCATION)
	if err != nil {
		log.Errorf("error parsing template (%s): %s", TEMPLATE_LOCATION, err.Error())
		return nil, fmt.Errorf("error parsing template (%s): %s", TEMPLATE_LOCATION, err.Error())
	}
	return &UI{
		template:  tmpl,
		vpnc:      vpnc,
		ipEchoUrl: ipEchoUrl,
	}, nil
}

func (ui *UI) getIPEcho() (string, error) {

	req, err := http.NewRequest(http.MethodGet, ui.ipEchoUrl, nil)
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
		respBytes, err := ioutil.ReadAll(resp.Body)
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

func addCommandResult(commandResults []string, commandResult string) []string {
	if commandResults == nil {
		commandResults = []string{}
	}
	return append(commandResults, commandResult)
}

func (ui *UI) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	uiData := &uiData{}
	var err error
	uiData.ConfigurationList, err = ui.vpnc.ConfigurationList()
	if err != nil {
		log.Errorf("error getting configurations: %s", err)
		uiData.CommandResults =
			addCommandResult(uiData.CommandResults, fmt.Sprintf("Error getting configurations: %s", err))
	}

	// There is stuff to do
	if r.Method == http.MethodPost {
		r.ParseForm()
		defer r.Body.Close()

		vpnconfig := r.Form.Get("vpnconfig")
		if vpnconfig == "disconnect" {

			if commandResult, err := ui.vpnc.Disconnect(); err != nil {
				log.Errorf("error disconnecting: %s", err)
				uiData.CommandResults =
					addCommandResult(uiData.CommandResults, fmt.Sprintf("error disconnecting: %s", err))

			} else {
				uiData.CommandResults =
					addCommandResult(uiData.CommandResults, commandResult)
			}

		} else if vpnconfig != "" {

			if commandResult, err := ui.vpnc.Connect(vpnconfig); err != nil {
				log.Errorf("error connecting: %s", err)
				uiData.CommandResults =
					addCommandResult(uiData.CommandResults, fmt.Sprintf("error connecting: %s", err))

			} else {
				uiData.CommandResults =
					addCommandResult(uiData.CommandResults, commandResult)
			}
		}
	}

	if isConnected := ui.vpnc.CheckConnection(); isConnected {
		uiData.ConnectionState = "Connected"
	} else {
		uiData.ConnectionState = "Disconnected"
	}

	uiData.OwnPublicIP, err = ui.getIPEcho()
	if err != nil {
		log.Errorf("error determining own public ip: %s", err)
		uiData.CommandResults =
			addCommandResult(uiData.CommandResults, fmt.Sprintf("error determining own public ip: %s", err))

	}

	log.Infof("Host IP is: %s", uiData.OwnPublicIP)

	ui.template.Execute(w, uiData)
}
