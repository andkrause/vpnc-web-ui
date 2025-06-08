all: generate ui-build build docker

.PHONY: generate docker clean ui-build ui-install go-deps

generate: 
	rm -rf gen
	openapi-generator generate -i api/openapi.yaml  -g go-server --additional-properties=sourceFolder=gen/vpnapi,featureCORS=true,outputAsLibrary=true,packageName=vpnapi,onlyInterfaces=true
	goimports -srcdir gen/vpncapi -w ./

go-deps:
	go mod download && go mod tidy

ui-install:
	cd ui && npm install

ui-build: ui-install
	cd ui && npm run build

build: main.go go-deps ui-build
	go build -o vpnc-web-ui main.go
	
clean:
	rm -f vpnc-web-ui
	rm -rf ui/dist

