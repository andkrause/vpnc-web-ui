#!/bin/sh

# enable forwarding
echo 1 > /proc/sys/net/ipv4/ip_forward

# enable nating for vpn config
iptables -A FORWARD -i $VPNC_INTERFACE -o $LAN_INTERFACE -m state --state RELATED,ESTABLISHED -j ACCEPT
iptables -A FORWARD -i $WIREGUARD_INTERFACE -o $LAN_INTERFACE -m state --state RELATED,ESTABLISHED -j ACCEPT

iptables -t nat -A POSTROUTING -o $VPNC_INTERFACE -j MASQUERADE
iptables -t nat -A POSTROUTING -o $WIREGUARD_INTERFACE -j MASQUERADE

iptables -A FORWARD -i $LAN_INTERFACE -o $VPNC_INTERFACE -j ACCEPT
iptables -A FORWARD -i $LAN_INTERFACE -o $WIREGUARD_INTERFACE -j ACCEPT

#start webapp
./vpnc-web-ui
