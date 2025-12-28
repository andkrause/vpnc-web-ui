## Build
FROM node:lts-alpine3.22 AS ui-builder

WORKDIR /ui
COPY ui/ .

RUN npm install
RUN npm run build

FROM golang:1.25-alpine3.22 AS go-builder

WORKDIR /app

COPY go.mod ./
COPY pkg/ pkg/
COPY gen/ gen/
RUN go mod download

COPY main.go ./

RUN go mod tidy && go build -o /vpnc-web-ui main.go

## Deploy
FROM alpine:3.22

ENV VPNC_INTERFACE=tun0
ENV WIREGUARD_INTERFACE=wg0
ENV LAN_INTERFACE=eth0

# important to remove sysctl reference. this needs to be set on docker run command
RUN apk add --no-cache vpnc iptables wireguard-tools \
    && sed -i 's|\[\[ $proto == -4 \]\] && cmd sysctl -q net\.ipv4\.conf\.all\.src_valid_mark=1|[[ $proto == -4 ]] \&\& [[ $(sysctl -n net.ipv4.conf.all.src_valid_mark) != 1 ]] \&\& cmd sysctl -q net.ipv4.conf.all.src_valid_mark=1|' /usr/bin/wg-quick

WORKDIR /app
COPY --from=ui-builder /ui/dist/vpn-gateway-ui ui/dist/vpn-gateway-ui
COPY --from=go-builder /vpnc-web-ui vpnc-web-ui
COPY conf/dockerconfig.json conf/config.json 
COPY docker/docker-entrypoint.sh ./

EXPOSE 80

ENTRYPOINT ["/app/docker-entrypoint.sh"]
