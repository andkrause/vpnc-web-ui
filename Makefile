all: generate build docker

.PHONY: generate docker clean

generate: 
	openapi-generator generate -i api/openapi.yaml  -g go-server --additional-properties=sourceFolder=gen/vpnapi,featureCORS=true,outputAsLibrary=true,packageName=vpnapi,onlyInterfaces=true
	goimports -srcdir gen/vpncapi -w ./

build: main.go generate
	go mod tidy && go build -o vpnc-web-ui main.go
	
clean:
	rm vpnc-web-ui
	rm -r gen

