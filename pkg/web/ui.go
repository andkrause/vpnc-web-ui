package web

import (
	"fmt"
	"net/http"

	"html/template"

	"github.com/andkrause/vpnc-web-ui/pkg/vpnc"
	log "github.com/sirupsen/logrus"
)

const TEMPLATE_LOCATION = "./templates/index.html"

type uiData struct {
	ConfigurationList []string
	CommandResults    []string
	OwnPublicIP       string
	ConnectionState   string
}

type UI struct {
	template *template.Template
	vpnc     *vpnc.VPNC
}

func New(vpnc *vpnc.VPNC, ipEchoUrl string) (*UI, error) {
	tmpl, err := template.ParseFiles(TEMPLATE_LOCATION)
	if err != nil {
		log.Errorf("error parsing template (%s): %s", TEMPLATE_LOCATION, err.Error())
		return nil, fmt.Errorf("error parsing template (%s): %s", TEMPLATE_LOCATION, err.Error())
	}
	return &UI{
		template: tmpl,
		vpnc:     vpnc,
	}, nil
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

			if err := ui.vpnc.Disconnect(); err != nil {
				log.Errorf("error disconnecting: %s", err)
				uiData.CommandResults =
					addCommandResult(uiData.CommandResults, fmt.Sprintf("error disconnecting: %s", err))
			}

		} else if vpnconfig != "" {

			if err := ui.vpnc.Connect(vpnconfig); err != nil {
				log.Errorf("error connecting: %s", err)
				uiData.CommandResults =
					addCommandResult(uiData.CommandResults, fmt.Sprintf("error connecting: %s", err))

			}
		}
	}

	vpncStatus := ui.vpnc.Status()

	uiData.CommandResults =
		addCommandResult(uiData.CommandResults, vpncStatus.Message)

	if vpncStatus.ActiveVpnConfig != "" {
		uiData.ConnectionState = fmt.Sprintf("Connected to %q", vpncStatus.ActiveVpnConfig)
	} else {
		uiData.ConnectionState = "Disconnected"
	}

	uiData.OwnPublicIP = vpncStatus.CurrentPublicIp

	log.Infof("Host IP is: %s", uiData.OwnPublicIP)

	ui.template.Execute(w, uiData)
}
