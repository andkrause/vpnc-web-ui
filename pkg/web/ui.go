package web

import (
	"fmt"
	"net/http"
	"strings"

	"html/template"

	"github.com/andkrause/vpnc-web-ui/pkg/vpnclient"
	log "github.com/sirupsen/logrus"
)

const TEMPLATE_LOCATION = "./templates/index.html"

type uiVpnConfig struct {
	VpnClientName        string
	VpnConfigurationName string
}

type uiData struct {
	VpnConfigurationList []uiVpnConfig
	CommandResults       []string
	OwnPublicIP          string
	ConnectionState      string
}

type UI struct {
	template            *template.Template
	vpnClientAggregator *vpnclient.VpnClientAggregator
}

func New(vpnClientAggregator *vpnclient.VpnClientAggregator) (*UI, error) {
	tmpl, err := template.ParseFiles(TEMPLATE_LOCATION)
	if err != nil {
		log.Errorf("error parsing template (%s): %s", TEMPLATE_LOCATION, err.Error())
		return nil, fmt.Errorf("error parsing template (%s): %s", TEMPLATE_LOCATION, err.Error())
	}
	return &UI{
		template:            tmpl,
		vpnClientAggregator: vpnClientAggregator,
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

	vpnConfigList, err := ui.vpnClientAggregator.ConfigurationList()
	if err != nil {
		log.Errorf("error getting configurations: %s", err)
		uiData.CommandResults =
			addCommandResult(uiData.CommandResults, fmt.Sprintf("Error getting configurations: %s", err))
	}

	for _, vpnClienConfigs := range vpnConfigList {
		configArray := make([]uiVpnConfig, len(vpnClienConfigs.AvailableVpnConfigs))
		for i, vpnClientConfigName := range vpnClienConfigs.AvailableVpnConfigs {
			configArray[i] = uiVpnConfig{
				VpnClientName:        vpnClienConfigs.VpnClientName,
				VpnConfigurationName: vpnClientConfigName,
			}
		}
		uiData.VpnConfigurationList = append(uiData.VpnConfigurationList, configArray...)
	}

	// There is stuff to do
	if r.Method == http.MethodPost {
		r.ParseForm()
		defer r.Body.Close()

		vpnconfig := r.Form.Get("vpnconfig")
		if vpnconfig == "disconnect" {

			if message, err := ui.vpnClientAggregator.Disconnect(); err != nil {
				log.Errorf("error disconnecting: %s", err)
				uiData.CommandResults =
					addCommandResult(uiData.CommandResults, fmt.Sprintf("error disconnecting: %s", err))
			} else {
				if len(message) > 0 {
					uiData.CommandResults = addCommandResult(uiData.CommandResults,
						fmt.Sprintf("success: %s", message))
				} else {
					uiData.CommandResults = addCommandResult(uiData.CommandResults,
						"successfully disconnected")
				}
			}

		} else if vpnconfig != "" {
			vpnConfigArray := strings.Split(vpnconfig, "#")

			if len(vpnConfigArray) != 2 {
				log.Errorf("vpnconfig name %q unexpected, use \"#\" as a delimiter", vpnconfig)
				uiData.CommandResults =
					addCommandResult(uiData.CommandResults,
						fmt.Sprintf("vpnconfig name %q unexpected, use \"#\" as a delimiter", vpnconfig))
			} else {

				if err := ui.vpnClientAggregator.Connect(vpnConfigArray[0], vpnConfigArray[1]); err != nil {
					log.Errorf("error connecting through vpn client %q: %s", vpnConfigArray[0], err)
					uiData.CommandResults =
						addCommandResult(uiData.CommandResults,
							fmt.Sprintf("error connecting through vpn client %q: %s", vpnConfigArray[0], err))

				}
			}
		}
	}

	vpnStatus := ui.vpnClientAggregator.Status()

	if len(vpnStatus.Message) > 0 {
		uiData.CommandResults =
			addCommandResult(uiData.CommandResults, vpnStatus.Message)
	}

	if vpnStatus.ActiveVpnConfig != "" {
		uiData.ConnectionState = fmt.Sprintf("Connected to \"%s %s\"", vpnStatus.ActiveVpnClient, vpnStatus.ActiveVpnConfig)
	} else {
		uiData.ConnectionState = "Disconnected"
	}

	uiData.OwnPublicIP = vpnStatus.CurrentPublicIp

	log.Infof("Host IP is: %s", uiData.OwnPublicIP)

	ui.template.Execute(w, uiData)
}
