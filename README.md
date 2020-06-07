# VPNC Web UI

## Table Of Contents 
- [VPNC Web UI](#vpnc-web-ui)
  - [Table Of Contents](#table-of-contents)
  - [Configuration](#configuration)
  - [Install](#install)
  - [Build](#build)

This is a small web ui on top of [vpnc](https://davidepucci.it/doc/vpnc/). I use this to remote-control my IPSec vpn gateway running on top of a Raspberry Pi 4b. It basically replaces the need to run shell commands.

## Configuration

The application is configured using a json config file. This consists of 2 parts: 
* The VPNC configuration
* The HTTP Server configuration

The default location is ./conf/config.json, howver this can be changed using the `-configFilePath` parameter. Below is the default config file:

```
{
    "vpnc": {
        "connectCommand": "/usr/sbin/vpnc",
        "disconnectCommand": "/usr/sbin/vpnc-disconnect",
        "pidFile": "/var/run/vpnc.pid",
        "configFolder": "/etc/vpnc/",
        "waitTimeAfterConnect": 1
    },
    "webUI": {
        "serverPort": 80,
        "ipEchoURL": "https://ipecho.net/plain"
    }
}
```

|Section|Option|Definition|
|-------|------|----------|
|vpnc|connectCommand|Location of the vpnc command, can be found using `which vpnc`. This is the command executed when connect is selected|
|vpnc|disconnectCommand|Location of the vpnc-disconnect command, can be found using `which vpnc-disconnect`. This is the command executed when disconnect is selected|
|vpnc|pidFile|vpnc is started in background. to keep track a file containing the current process id is created. This file (its existence) is used to derrive the current connection state|
|vpnc|configFolder|Folder where vpnc configs are searched|
|vpnc|waitTimeAfterConnect|VPNC runs in background (concurrently). This wait time is used to "synchronize" the UI and the backround job. Nothing bad happens if synchronization is not perfect. However UI might display a wrong IP and or connection state|
|webUI|serverPort|Port on which the UI is exposed. By default the server binds to all IPs/Hosts|
|webUI|ipEchoURL|URL that is invoked to determine the own (server side) IP|

The usi is rendered based on a web template. The template is located in `template/index.html`. If you dont like it, you can change it.

Appearance is driven by the ccs locate in `static/formatting.css`. Again, if you don't like it, change it.

## Install

The releases section contains a couple of versions that you can install without building. Just download and extract. Pick the file for your environment (currently linux/aarch64, macos and linux/arm64 are available) and execute.

To start the server use and init-script. Samples are found in this [init scripts directory](init-scripts/).

## Build 

Since this is a "normal" golang application it requires a golang environment to be installed. It is then built using the following command:

```
go build -o vpnc-web-ui  main.go 
```

To build for alternative OS/Platform like the Rasperrry Pi use:

```
GOOS=linux GOARCH=arm64 go build -o vpnc-web-ui-aarch64  main.go 
```




