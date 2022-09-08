NAME=agile-proxy
VERSION=1.1.0
GOBUILD=go build

all: linux-amd64 windows-amd64

linux-amd64:
	set CGO_ENABLED=0
	set GOOS=linux
	set GOARCH=amd64
	$(GOBUILD) -ldflags "-X agile-proxy/config.version=$(VERSION)" -o ./bin/$(NAME)

windows-amd64:
	set GOARCH=amd64
	set GOOS=windows
	$(GOBUILD) -ldflags "-X agile-proxy/config.version=$(VERSION)" -o ./bin/$(NAME).exe
