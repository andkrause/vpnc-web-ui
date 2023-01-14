## Build
FROM golang:1.19.5-alpine3.17 AS go-builder

WORKDIR /app

COPY go.mod ./
COPY pkg/ pkg/
RUN go mod download

COPY main.go ./

RUN go build -o /vpnc-web-ui main.go

## Deploy
FROM alpine:3.17

ENV VPN_INTERFACE=tun0
ENV LAN_INTERFACE=eth0

RUN apk add --no-cache vpnc iptables

WORKDIR /app

COPY --from=go-builder /vpnc-web-ui vpnc-web-ui
COPY static/ static/
COPY templates/ templates/ 
COPY conf/dockerconfig.json conf/config.json 
COPY docker/docker-entrypoint.sh ./

EXPOSE 80

ENTRYPOINT ["/app/docker-entrypoint.sh"]