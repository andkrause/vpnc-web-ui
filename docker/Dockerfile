## Build
FROM node:lts-alpine3.23 AS ui-builder
RUN apk add --no-cache make

WORKDIR /app
COPY . .
RUN make ui-build


FROM golang:1.25-alpine3.23 AS go-builder
RUN apk add --no-cache make

WORKDIR /app
COPY . .
RUN make go-build

## Put it all together to create somethiung that can run
FROM alpine:3.23

ENV VPNC_INTERFACE=tun0
ENV WIREGUARD_INTERFACE=wg0
ENV LAN_INTERFACE=eth0

# Install VPN clients and iptables (for routing)
# remove openresolv which is not needed anyway and be ready for no-of replacement
# important to remove sysctl reference. this needs to be set on docker run command
RUN apk add --no-cache vpnc iptables wireguard-tools \
    && apk del openresolv 2>/dev/null || true \
    && sed -i 's|\[\[ $proto == -4 \]\] && cmd sysctl -q net\.ipv4\.conf\.all\.src_valid_mark=1|[[ $proto == -4 ]] \&\& [[ $(sysctl -n net.ipv4.conf.all.src_valid_mark) != 1 ]] \&\& cmd sysctl -q net.ipv4.conf.all.src_valid_mark=1|' /usr/bin/wg-quick

# --- DNS handling note (important) -------------------------------------------
# This container runs WireGuard in "gateway" mode and MUST NOT let wg-quick
# modify /etc/resolv.conf.
#
# Alpine 3.23 pulls in openresolv >= 3.17, which enforces strict ownership
# of /etc/resolv.conf and expects an init system + resolver providers.
# In Docker, those do not exist. As a result:
#   - wg-quick invokes resolvconf when "DNS =" is present in wg config
#   - openresolv fails (signature mismatch / no providers)
#   - /etc/resolv.conf may be truncated or emptied, breaking DNS entirely
#
# We intentionally replace resolvconf with a no-op stub so wg-quick's DNS
# phase succeeds without touching DNS. 
#
# DNS is not handles by the container.
#
# This is a deliberate container-specific design choice, not a workaround.
# ---------------------------------------------------------------------------
COPY --chmod=+x ./docker/resolvconf /usr/bin/resolvconf

WORKDIR /app
COPY --from=ui-builder /app/ui/dist/vpn-gateway-ui ui/dist/vpn-gateway-ui
COPY --from=go-builder /app/vpnc-web-ui vpnc-web-ui
COPY conf/dockerconfig.json conf/config.json 
COPY docker/docker-entrypoint.sh ./

EXPOSE 80

ENTRYPOINT ["/app/docker-entrypoint.sh"]
