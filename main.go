package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/andkrause/vpnc-web-ui/pkg/vpnc"

	"github.com/andkrause/vpnc-web-ui/pkg/config"
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

	vpncClient := vpnc.New(serverConfig.VPNC.ConnectCommand, serverConfig.VPNC.DisconnectCommand,
		serverConfig.VPNC.ConfigFolder, serverConfig.VPNC.PIDFile, serverConfig.VPNC.WaitTimeAfterConnect)

	//Serve UI
	ui, err := web.New(vpncClient, serverConfig.WebUI.IPEchoURL)
	if err != nil {
		log.Error(err.Error())
		os.Exit(2)
	}

	http.Handle("/", ui)

	// Serve static stuff
	staticFileServer := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix(strings.TrimRight("/static/", "/"), staticFileServer))

	server := http.Server{
		Addr: fmt.Sprintf(":%d", serverConfig.WebUI.ServerPort),
	}

	fmt.Println("Starting server with config:")
	serverConfig.LogConfig()

	log.Fatal(server.ListenAndServe())

}
