package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/andkrause/vpnc-web-ui/gen/vpnapi"
	"github.com/andkrause/vpnc-web-ui/pkg/api"
	"github.com/andkrause/vpnc-web-ui/pkg/config"
	"github.com/andkrause/vpnc-web-ui/pkg/vpnc"
	"github.com/andkrause/vpnc-web-ui/pkg/web"
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
		log.Error(err.Error())
		os.Exit(2)
	}

	vpncClient := vpnc.New(serverConfig.ConnectCommand, serverConfig.DisconnectCommand,
		serverConfig.ConfigFolder, serverConfig.WaitTimeAfterConnect,
		serverConfig.IPEchoURL, serverConfig.GetMaxAgePublicIpDuration())

	//Serve UI
	ui, err := web.New(vpncClient)
	if err != nil {
		log.Error(err.Error())
		os.Exit(2)
	}

	//API stuff

	//Implementation
	services := api.New(vpncClient)

	//Controllers
	vpnConnectionApi := vpnapi.NewVpnConnectionApiController(services)
	vpnGatewayApi := vpnapi.NewVpnGatewayApiController(services)

	//API Router
	router := vpnapi.NewRouter(vpnConnectionApi, vpnGatewayApi)

	//UI Handler
	router.Handle("/", ui)

	// Serve static stuff
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static", http.FileServer(http.Dir("./static"))))

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", serverConfig.ServerPort),
		Handler: router,
	}

	fmt.Println("Starting server with config:")
	serverConfig.LogConfig()

	log.Fatal(server.ListenAndServe())

}
