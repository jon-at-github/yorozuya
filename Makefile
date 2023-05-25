GO_VERSION := 1.20
ARCH := armv6l
GO_BINARY_LATEST := go$(GO_VERSION).linux-$(ARCH).tar.gz
TAG := $$(git describe --abbrev=0 --tags --always)
HASH := $$(git rev-parse HEAD)
DATE := $$(date +%Y-%m-%d.%H:%M:%S)
LDFLAGS := -w -X github.com/jon-at-github/hello-api/handlers.hash=$(HASH) -X github.com/jon-at-github/hello-api/handlers.tag=$(TAG) -X github.com/jon-at-github/hello-api/handlers.date=$(DATE)

setup: install-go init-go build

install-go:
	sudo rm -rf /usr/local/go
	sudo tar -C /usr/local -xvf ${GO_BINARY_LATEST}
	rm ${GO_BINARY_LATEST}

init-go:
	echo 'export PATH=$$PATH:/usr/local/go/bin' >> $${HOME}.bashrc
	echo 'export PATH=$$PATH:$${HOME}/go/bin' >> $${HOME}.bashrc

upgrade-go:
	sudo rm -rf /usr/local/go
	wget "https://go.dev/dl/${GO_BINARY_LATEST}"
	sudo  tar -C /usr/local -xzf ${GO_BINARY_LATEST}
	rm ${GO_BINARY_LATEST}

build:
	go build -o api cmd/main.go