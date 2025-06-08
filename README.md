# VPNC Web UI

## Table Of Contents 
- [VPNC Web UI](#vpnc-web-ui)
  - [Table Of Contents](#table-of-contents)
  - [Architecture](#architecture)
  - [Configuration](#configuration)
  - [REST API](#rest-api)
  - [Install](#install)
  - [Build](#build)
  - [Development](#development)
  - [Docker](#docker)

This is a small web ui on top of [vpnc](https://davidepucci.it/doc/vpnc/) and [wireguard](https://www.wireguard.com/). I use this to remote-control my IPSec/Wireguard vpn gateway running on top of a Raspberry Pi 4b. It basically replaces the need to run shell commands.

## Architecture

The application consists of two main components:
- **Backend**: Go server providing REST API endpoints for VPN management
- **Frontend**: Modern Angular web application with responsive design and real-time status updates

## Configuration

The application is configured using a json config file. This consists of 3 parts: 
* The VPNC configuration
* The Wireguard configuration
* The HTTP Server configuration

The default location is ./conf/config.json, howver this can be changed using the `-configFilePath` parameter. Below is the default config file:

```
{
    "waitTimeAfterConnect": 1,
    "serverPort": 80,
    "ipEchoURL": "https://ipecho.net/plain",
    "maxAgePublicIp": "2h",
    "vpnc": {
        "connectCommand": "/usr/sbin/vpnc",
        "disconnectCommand": "/usr/sbin/vpnc-disconnect",
        "configFolder": "/etc/vpnc/",
        "vpncNetworkInterfaceName": "tun0"
    },
    "wireguard": {
        "wgQuickCommand": "/usr/bin/wg-quick",
        "configFolder": "/var/wireguard/config/",
        "wireguardNetworkInterfaceName": "wg0",
        "wgQuickConfigSearchDir": "/etc/wireguard/"
    }
}
```

|Option|Definition|
|------|----------|
waitTimeAfterConnect|VPNC runs in background (concurrently). This wait time is used to "synchronize" the UI and the backround job. Nothing bad happens if synchronization is not perfect. However UI might display a wrong IP and or connection state|
|serverPort|Port on which the UI is exposed. By default the server binds to all IPs/Hosts|
|ipEchoURL|URL that is invoked to determine the own (server side) IP|
|maxAgePublicIp| To avoid asking for the public IP address too often it is cached. This parameter specifies the time after which the cache expires. Naturally connection chnages also expire the cache |
|vpnc/connectCommand|Location of the vpnc command, can be found using `which vpnc`. This is the command executed when connect is selected|
|vpnc/disconnectCommand|Location of the vpnc-disconnect command, can be found using `which vpnc-disconnect`. This is the command executed when disconnect is selected|
|vpnc/configFolder|Folder where vpnc configs are searched|
|vpnc/vpncNetworkInterfaceName|Is the name of the network interface that VPNC assigns to an active VPN connection. Default is `tun0`|
|wireguard/wgQuickCommand|Location of the wg-quick command, can be found using `which wg-quick`. This is the command executed when connect is selected|
|wireguard/configFolder|Folder where wg-quick configs are searched|
|wireguard/wireguardNetworkInterfaceName|Name of the OS network interface that should be assigned to all wireguard connections. Default is `wg0`|
|wireguard/wgQuickConfigSearchDir|Directory that wg-quick uses to search for *.conf files. For details see the [wg-quick manpages](https://manpages.debian.org/unstable/wireguard-tools/wg-quick.8.en.html). Default is /etc/wireguard/. |

The UI is rendered based on a web template. The template is located in `template/index.html`. If you dont like it, you can change it.

Appearance is driven by the ccs locate in `static/formatting.css`. Again, if you don't like it, change it.

## REST API

To control the Gateway from Home Automation platforms like Home Assistant or openHAB there i a REST API included. This is documented in an [Open API](api/openapi.yaml) file. You can open it in the [swagger editor](https://editor.swagger.io) or any other tool that renders Open API specs.

To regenerate the API spec the following tools are required:

- [goimports](https://pkg.go.dev/golang.org/x/tools/cmd/goimports): to remove unneccessary imports from the generator
- [OpenAPI Generator](https://openapi-generator.tech/): Generates a golang server from an Open API spec

`make generate` executes the generation.

## Install

The releases section contains a couple of versions that you can install without building. Just download and extract. Pick the file for your environment (currently linux/aarch64, macos and linux/arm64 are available) and execute.

To start the server use and init-script. Samples are found in this [init scripts directory](init-scripts/).

## Build 

The application requires both Go and Node.js environments to build:

**Prerequisites:**
- Go 1.24.4+ for the backend server
- Node.js 18+ and npm for the Angular frontend
- OpenAPI Generator and goimports (for API generation)

**Build Commands:**

```bash
# Build everything (API generation, UI build, and Go server)
make all

# Build only the backend server
make build

# Build only the Angular UI
make ui-build

# Install UI dependencies
make ui-install

# Generate API code from OpenAPI spec
make generate

# Clean build artifacts
make clean
```

The build process:
1. `make generate` - Generates Go server code from the OpenAPI specification
2. `make ui-build` - Builds the Angular application (installs npm dependencies if needed) 
3. `make build` - Compiles the Go server with the embedded Angular assets
4. Final binary: `vpnc-web-ui`

## Development

**Frontend Development:**
The Angular application is located in the `ui/` directory. For development:

```bash
cd ui
npm install
npm run serve
```

This starts the Angular development server with hot reload at http://localhost:4200. The development server proxies API requests to the Go backend.

**Backend Development:**
The Go server serves both the API and the built Angular application. During development, you can run the Go server separately and the Angular dev server will proxy requests to it.

**Available Make Targets:**
- `make all` - Complete build pipeline
- `make build` - Build Go server only
- `make ui-build` - Build Angular UI only  
- `make ui-install` - Install npm dependencies
- `make generate` - Generate API code
- `make go-deps` - Download and tidy Go dependencies
- `make clean` - Remove build artifacts

## Docker

The gateway can also be deployed as a docker container. I use this to simplify maintenance and ensure portability. 

To run the gateway on a dedicated IP with **full** network access (which is required), I use  a [macvlan network](https://docs.docker.com/network/macvlan/). This network named `docker_public_services` is pre-created using the following command:

```
docker network create -d macvlan -o parent=<network interface name> \
  --subnet <cidr of the subnet> \
  --gateway <"real" gateway in the subnet> \
  --ip-range <if addresses you want to assign> \
  docker_public_services

```

If you want to run just one instance a host network might work as well.

To start the container as a service and attach it to the pre-created macvlan `docker_public_services` the following compose file can serve as a baseline :

```
version: "2.4"
services:
    vpnc:
        image: andy008/vpnc-web-ui:latest
        init: true
        restart: "always"
        sysctls:
            - net.ipv4.conf.all.src_valid_mark=1
        cap_add:
            - NET_ADMIN
            - NET_RAW
        mem_limit: 256m
        cpus: 0.5
        networks:
            docker_public_services: 
                ipv4_address: "<your gateways fixed local IP>"
        volumes:
            -  type: volume
               source: vpnc_config
               target: /var/vpnc/mountedconfig/
               read_only: false
            -  type: volume
               source: wireguard_config
               target: /var/wireguard/mountedconfig/
               read_only: false  

networks:
  docker_public_services:
    external: true


volumes:
    vpnc_config:
       driver: local
       driver_opts:
          o: bind
          type: none
          device: <path to vpnc config>
    wireguard_config:
       driver: local
       driver_opts:
          o: bind
          type: none
          device: <path to wireguard config>

```

To enable this on a fresh ubuntu 22.04 LTS on a raspberry PI I had to install `linux-modules-extra-raspi` and add the below iptables firewall rule: `iptables -I FORWARD -i eth0 -o eth0 -j ACCEPT`

Since iptables are not persisted by default, you need to apply it after every restart (or persist iptables ;-))