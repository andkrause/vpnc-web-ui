package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/andkrause/vpnc-web-ui/gen/vpnapi"
	"github.com/andkrause/vpnc-web-ui/pkg/api"
	"github.com/andkrause/vpnc-web-ui/pkg/config"
	"github.com/andkrause/vpnc-web-ui/pkg/vpnc"
	"github.com/andkrause/vpnc-web-ui/pkg/vpnclient"
	wgquick "github.com/andkrause/vpnc-web-ui/pkg/wg-quick"
	log "github.com/sirupsen/logrus"
)

func parseInputs() (configFilePath string) {
	flag.StringVar(&configFilePath, "configFilePath", "conf/config.json", "Path to vpnc-web-ui config file (json format)")
	flag.Parse()
	return
}

func main() {

	configFilePath := parseInputs()

	serverConfig, err := config.ParseConfigFile(configFilePath)
	if err != nil {
		log.Fatal(err.Error())
	}

	if err := serverConfig.Validate(); err != nil {
		log.Fatalf("Invalid configuration in %q: %s", configFilePath, err)
	}

	vpnContainer := []*vpnclient.VpnClientContainer{}

	if serverConfig.IsVpncConfigured() {

		vpncClient := &vpnclient.VpnClientContainer{
			Name: "vpnc",
			Client: vpnc.New(serverConfig.Vpnc.ConnectCommand, serverConfig.Vpnc.DisconnectCommand,
				serverConfig.Vpnc.ConfigFolder, serverConfig.Vpnc.VpncNetworkInterfaceName),
		}

		vpnContainer = append(vpnContainer, vpncClient)
	}

	if serverConfig.IsWireguardConfigured() {
		wireguardClient := &vpnclient.VpnClientContainer{
			Name: "wireguard",
			Client: wgquick.New(serverConfig.Wireguard.WgQuickCommand, serverConfig.Wireguard.ConfigFolder,
				serverConfig.Wireguard.WireguardNetworkInterfaceName,
				serverConfig.Wireguard.WgQuickConfigSearchDir),
		}

		vpnContainer = append(vpnContainer, wireguardClient)

	}

	vpnAggregator := vpnclient.New(serverConfig.WaitTimeAfterConnect,
		serverConfig.IPEchoURL, serverConfig.GetMaxAgePublicIpDuration(), vpnContainer...,
	)

	//API stuff

	//Implementation
	services := api.New(vpnAggregator)

	//Controllers
	vpnConnectionApi := vpnapi.NewVpnConnectionAPIController(services)
	vpnGatewayApi := vpnapi.NewVpnGatewayAPIController(services)

	//API Router
	router := vpnapi.NewRouter(vpnConnectionApi, vpnGatewayApi)

	// Serve the Angular SPA with proper routing support
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("ui/dist/vpn-gateway-ui/browser")))

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", serverConfig.ServerPort),
		Handler: router,
	}

	fmt.Println("Starting server with config:")
	serverConfig.LogConfig()

	log.Fatal(server.ListenAndServe())

}
