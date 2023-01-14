#!/bin/sh

# enable forwarding
echo 1 > /proc/sys/net/ipv4/ip_forward

# enable nating for vpn config
iptables -t nat -A POSTROUTING -o $VPN_INTERFACE -j MASQUERADE
iptables -A FORWARD -i tun0 -o $LAN_INTERFACE -m state --state RELATED,ESTABLISHED -j ACCEPT
iptables -A FORWARD -i $LAN_INTERFACE -o $VPN_INTERFACE -j ACCEPT

#start webapp
./vpnc-web-ui
