## Build
FROM golang:1.22.2-alpine3.19 AS go-builder

WORKDIR /app

COPY go.mod ./
COPY pkg/ pkg/
COPY gen/ gen/
RUN go mod download

COPY main.go ./

RUN go mod tidy && go build -o /vpnc-web-ui main.go

## Deploy
FROM alpine:3.19

ENV VPNC_INTERFACE=tun0
ENV WIREGUARD_INTERFACE=wg0
ENV LAN_INTERFACE=eth0

# important to remove sysctl reference. this needs to be set on docker run command
RUN apk add --no-cache vpnc iptables wireguard-tools \
    && sed -i 's|\[\[ $proto == -4 \]\] && cmd sysctl -q net\.ipv4\.conf\.all\.src_valid_mark=1|[[ $proto == -4 ]] \&\& [[ $(sysctl -n net.ipv4.conf.all.src_valid_mark) != 1 ]] \&\& cmd sysctl -q net.ipv4.conf.all.src_valid_mark=1|' /usr/bin/wg-quick

WORKDIR /app

COPY --from=go-builder /vpnc-web-ui vpnc-web-ui
COPY static/ static/
COPY templates/ templates/ 
COPY conf/dockerconfig.json conf/config.json 
COPY docker/docker-entrypoint.sh ./

EXPOSE 80

ENTRYPOINT ["/app/docker-entrypoint.sh"]